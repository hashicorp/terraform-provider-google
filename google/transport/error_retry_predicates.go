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
