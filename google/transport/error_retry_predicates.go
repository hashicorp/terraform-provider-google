package transport

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RetryErrorPredicateFunc func(error) (bool, string)

/** ADD GLOBAL ERROR RETRY PREDICATES HERE **/
// Retry predicates that shoud apply to all requests should be added here.
var defaultErrorRetryPredicates = []RetryErrorPredicateFunc{
	// Common network errors (usually wrapped by URL error)
	isNetworkTemporaryError,
	isNetworkTimeoutError,
	isIoEOFError,
	isConnectionResetNetworkError,

	// Common GCP error codes
	isCommonRetryableErrorCode,

	//While this might apply only to Cloud SQL, historically,
	// we had this in our global default error retries.
	// Keeping it as a default for now.
	is409OperationInProgressError,

	// GCE Error codes- we don't have a way to add these to all GCE resources
	// easily, so add them globally.

	// GCE Subnetworks are considered unready for a brief period when certain
	// operations are performed on them, and the scope is likely too broad to
	// apply a mutex. If we attempt an operation w/ an unready subnetwork, retry
	// it.
	isSubnetworkUnreadyError,

	// As of February 2022 GCE seems to have added extra quota enforcement on
	// reads, causing significant failure for our CI and for large customers.
	// GCE returns the wrong error code, as this should be a 429, which we retry
	// already.
	is403QuotaExceededPerMinuteError,
}

/** END GLOBAL ERROR RETRY PREDICATES HERE **/

func isNetworkTemporaryError(err error) (bool, string) {
	if netErr, ok := err.(*net.OpError); ok && netErr.Temporary() {
		return true, "marked as timeout"
	}
	if urlerr, ok := err.(*url.Error); ok && urlerr.Temporary() {
		return true, "marked as timeout"
	}
	return false, ""
}

func isNetworkTimeoutError(err error) (bool, string) {
	if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
		return true, "marked as timeout"
	}
	if urlerr, ok := err.(*url.Error); ok && urlerr.Timeout() {
		return true, "marked as timeout"
	}
	return false, ""
}

func isIoEOFError(err error) (bool, string) {
	if err == io.ErrUnexpectedEOF {
		return true, "Got unexpected EOF"
	}

	if urlerr, urlok := err.(*url.Error); urlok {
		wrappedErr := urlerr.Unwrap()
		if wrappedErr == io.ErrUnexpectedEOF {
			return true, "Got unexpected EOF"
		}
	}
	return false, ""
}

const connectionResetByPeerErr = ": connection reset by peer"

func isConnectionResetNetworkError(err error) (bool, string) {
	if strings.HasSuffix(err.Error(), connectionResetByPeerErr) {
		return true, fmt.Sprintf("reset connection error: %v", err)
	}
	return false, ""
}

// Retry 409s because some APIs like Cloud SQL throw a 409 if concurrent calls
// are being made.
//
// The only way right now to determine it is a retryable 409 due to
// concurrent calls is to look at the contents of the error message.
// See https://github.com/hashicorp/terraform-provider-google/issues/3279
func is409OperationInProgressError(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code == 409 && strings.Contains(gerr.Body, "operationInProgress") {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code 409 and error reason 'operationInProgress': %s", err)
		return true, "Operation still in progress"
	}
	return false, ""
}

func isSubnetworkUnreadyError(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code == 400 && strings.Contains(gerr.Body, "resourceNotReady") && strings.Contains(gerr.Body, "subnetworks") {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code 400 and error reason 'resourceNotReady' w/ `subnetwork`: %s", err)
		return true, "Subnetwork not ready"
	}
	return false, ""
}

// GCE (and possibly other APIs) incorrectly return a 403 rather than a 429 on
// rate limits.
func is403QuotaExceededPerMinuteError(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}
	var QuotaRegex = regexp.MustCompile(`Quota exceeded for quota metric '(?P<Metric>.*)' and limit '(?P<Limit>.* per minute)' of service`)
	if gerr.Code == 403 && QuotaRegex.MatchString(gerr.Body) {
		matches := QuotaRegex.FindStringSubmatch(gerr.Body)
		metric := matches[QuotaRegex.SubexpIndex("Metric")]
		limit := matches[QuotaRegex.SubexpIndex("Limit")]
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code 403 and error message 'Quota exceeded for quota metric `%s`: %s", metric, err)
		return true, fmt.Sprintf("Waiting for quota limit %s to refresh", limit)
	}
	return false, ""
}

// We've encountered a few common fingerprint-related strings; if this is one of
// them, we're confident this is an error due to fingerprints.
var FINGERPRINT_FAIL_ERRORS = []string{"Invalid fingerprint.", "Supplied fingerprint does not match current metadata fingerprint."}

// Retry the operation if it looks like a fingerprint mismatch.
func IsFingerprintError(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code != 412 {
		return false, ""
	}

	for _, msg := range FINGERPRINT_FAIL_ERRORS {
		if strings.Contains(err.Error(), msg) {
			return true, "fingerprint mismatch"
		}
	}

	return false, ""
}

const METADATA_FINGERPRINT_RETRIES = 10

// Since the google compute API uses optimistic locking, there is a chance
// we need to resubmit our updated metadata. To do this, you need to provide
// an update function that attempts to submit your metadata
func MetadataRetryWrapper(update func() error) error {
	attempt := 0
	for attempt < METADATA_FINGERPRINT_RETRIES {
		err := update()
		if err == nil {
			return nil
		}

		if ok, _ := IsFingerprintError(err); !ok {
			// Something else went wrong, don't retry
			return err
		}

		log.Printf("[DEBUG] Dismissed an error as retryable as a fingerprint mismatch: %s", err)
		attempt++
	}
	return fmt.Errorf("Failed to update metadata after %d retries", attempt)
}

// If a permission necessary to provision a resource is created in the same config
// as the resource itself, the permission may not have propagated by the time terraform
// attempts to create the resource. This allows those errors to be retried until the timeout expires
func IamMemberMissing(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "permission") {
			return true, "Waiting for IAM member permissions to propagate."
		}
	}
	return false, ""
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
func IsSqlOperationInProgressError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
		if strings.Contains(gerr.Body, "instanceAlreadyExists") {
			return false, ""
		}

		return true, "Waiting for other concurrent Cloud SQL operations to finish"
	}
	return false, ""
}

// Retry if service usage decides you're activating the same service multiple
// times. This can happen if a service and a dependent service aren't batched
// together- eg container.googleapis.com in one request followed by compute.g.c
// in the next (container relies on compute and implicitly activates it)
func ServiceUsageServiceBeingActivated(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 400 {
		if strings.Contains(gerr.Body, "Another activation or deactivation is in progress") {
			return true, "Waiting for same service activation/deactivation to finish"
		}

		return false, ""
	}
	return false, ""
}

// See https://github.com/hashicorp/terraform-provider-google/issues/14691 for
// details on the error message this handles
// This is a post-operation error so it uses tpgresource.CommonOpError instead of googleapi.Error
func ServiceUsageInternalError160009(err error) (bool, string) {
	// a cyclical dependency between transport/tpgresource blocks using tpgresource.CommonOpError
	// so just work off the error string. Ideally, we'd use that type instead.
	s := err.Error()
	if strings.Contains(s, "encountered internal error") && strings.Contains(s, "160009") && strings.Contains(s, "with failed services") {
		return true, "retrying internal error 160009."
	}

	return false, ""
}

// Retry if Bigquery operation returns a 403 with a specific message for
// concurrent operations (which are implemented in terms of 'edit quota').
func IsBigqueryIAMQuotaError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 403 && strings.Contains(strings.ToLower(gerr.Body), "exceeded rate limits") {
			return true, "Waiting for Bigquery edit quota to refresh"
		}
	}
	return false, ""
}

// Retry if Monitoring operation returns a 409 with a specific message for
// concurrent operations.
func IsMonitoringConcurrentEditError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && (strings.Contains(strings.ToLower(gerr.Body), "too many concurrent edits") ||
			strings.Contains(strings.ToLower(gerr.Body), "could not fulfill the request")) {
			return true, "Waiting for other Monitoring changes to finish"
		}
	}
	return false, ""
}

// Retry if KMS CryptoKeyVersions returns a 400 for PENDING_GENERATION
func IsCryptoKeyVersionsPendingGeneration(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 400 {
		if strings.Contains(gerr.Body, "PENDING_GENERATION") {
			return true, "Waiting for pending key generation"
		}
	}
	return false, ""
}

// Retry if getting a resource/operation returns a 404 for specific operations.
// opType should describe the operation for which 404 can be retryable.
func IsNotFoundRetryableError(opType string) RetryErrorPredicateFunc {
	return func(err error) (bool, string) {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return true, fmt.Sprintf("Retry 404s for %s", opType)
		}
		return false, ""
	}
}

func IsPeeringOperationInProgress(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "There is a peering operation in progress") {
			return true, "Waiting peering operation to complete"
		}
	}
	return false, ""
}

func DatastoreIndex409Contention(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && strings.Contains(gerr.Body, "too much contention") {
			return true, "too much contention - waiting for less activity"
		}
	}
	return false, ""
}

func IapClient409Operation(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && strings.Contains(strings.ToLower(gerr.Body), "operation was aborted") {
			return true, "operation was aborted possibly due to concurrency issue - retrying"
		}
	}
	return false, ""
}

func HealthcareDatasetNotInitialized(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 404 && strings.Contains(strings.ToLower(gerr.Body), "dataset not initialized") {
			return true, "dataset not initialized - retrying"
		}
	}
	return false, ""
}

// Cloud Run APIs may return a 409 on create to indicate that a resource has been deleted in the foreground
// (eg GET and LIST) but not the backing apiserver. When we encounter a 409, we can retry it.
// Note that due to limitations in MMv1's error_retry_predicates this is currently applied to all requests.
// We only expect to receive it on create, though.
func IsCloudRunCreationConflict(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 {
			return true, "saw a 409 - waiting until background deletion completes"
		}
	}

	return false, ""
}

// If a service account is deleted in the middle of updating an IAM policy
// it can cause the API to return an error. In fine-grained IAM resources we
// read the policy, modify it, then send it back to the API. Retrying is
// useful particularly in high-traffic projects.
// We don't want to retry _every_ time we see this error because the
// user-provided SA could trigger this too. At the callsite, we should check
// if the current etag matches the old etag and short-circuit if they do as
// that indicates the new config is the likely problem.
func IamServiceAccountNotFound(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "Service account") && strings.Contains(gerr.Body, "does not exist") {
			return true, "service account not found in IAM"
		}
	}

	return false, ""
}

// Concurrent Apigee operations can fail with a 400 error
func IsApigeeRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(strings.ToLower(gerr.Body), "the resource is locked by another operation") {
			return true, "Waiting for other concurrent operations to finish"
		}
	}

	return false, ""
}

func IsDataflowJobUpdateRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 404 && strings.Contains(gerr.Body, "in RUNNING OR DRAINING state") {
			return true, "Waiting for job to be in a valid state"
		}
	}
	return false, ""
}

// Cloud PubSub returns a 400 error if a topic's parent project was recently created and an
// organization policy has not propagated.
// See https://github.com/hashicorp/terraform-provider-google/issues/4349
func PubsubTopicProjectNotReady(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "retry this operation") {
			log.Printf("[DEBUG] Dismissed error as a retryable operation: %s", err)
			return true, "Waiting for Pubsub topic's project to properly initialize with organiation policy"
		}
	}
	return false, ""
}

// Retry on comon googleapi error codes for retryable errors.
// TODO(#5609): This may not need to be applied globally - figure out
// what retryable error codes apply to which API.
func isCommonRetryableErrorCode(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503 {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}
	return false, ""
}

// Retry if filestore operation returns a 429 with a specific message for
// concurrent operations.
func IsNotFilestoreQuotaError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 429 {
			return false, ""
		}
	}
	return isCommonRetryableErrorCode(err)
}

// Retry if App Engine operation returns a 409 with a specific message for
// concurrent operations, or a 404 indicating p4sa has not yet propagated.
func IsAppEngineRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && strings.Contains(strings.ToLower(gerr.Body), "operation is already in progress") {
			return true, "Waiting for other concurrent App Engine changes to finish"
		}
		if gerr.Code == 404 && strings.Contains(strings.ToLower(gerr.Body), "unable to retrieve p4sa") {
			return true, "Waiting for P4SA propagation to GAIA"
		}
	}
	return false, ""
}

// Bigtable uses gRPC and thus does not return errors of type *googleapi.Error.
// Instead the errors returned are *status.Error. See the types of codes returned
// here (https://pkg.go.dev/google.golang.org/grpc/codes#Code).
func IsBigTableRetryableError(err error) (bool, string) {
	// The error is retryable if the error code is not OK and has a retry delay.
	// The retry delay is currently not used.
	if errorStatus, ok := status.FromError(err); ok && errorStatus.Code() != codes.OK {
		var retryDelayDuration time.Duration
		for _, detail := range errorStatus.Details() {
			retryInfo, ok := detail.(*errdetails.RetryInfo)
			if !ok {
				continue
			}
			retryDelay := retryInfo.GetRetryDelay()
			retryDelayDuration = time.Duration(retryDelay.Seconds)*time.Second + time.Duration(retryDelay.Nanos)*time.Nanosecond
			break
		}
		if retryDelayDuration != 0 {
			// TODO: Consider sleep for `retryDelayDuration` before retrying.
			return true, "Bigtable operation failed with a retryable error, will retry"
		}
	}

	return false, ""
}
