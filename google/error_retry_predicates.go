package google

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
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
//The only way right now to determine it is a retryable 409 due to
// concurrent calls is to look at the contents of the error message.
// See https://github.com/terraform-providers/terraform-provider-google/issues/3279
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

// We've encountered a few common fingerprint-related strings; if this is one of
// them, we're confident this is an error due to fingerprints.
var FINGERPRINT_FAIL_ERRORS = []string{"Invalid fingerprint.", "Supplied fingerprint does not match current metadata fingerprint."}

// Retry the operation if it looks like a fingerprint mismatch.
func isFingerprintError(err error) (bool, string) {
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

// If a permission necessary to provision a resource is created in the same config
// as the resource itself, the permission may not have propagated by the time terraform
// attempts to create the resource. This allows those errors to be retried until the timeout expires
func iamMemberMissing(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "permission") {
			return true, "Waiting for IAM member permissions to propagate."
		}
	}
	return false, ""
}

// Cloud PubSub returns a 400 error if a topic's parent project was recently created and an
// organization policy has not propagated.
// See https://github.com/terraform-providers/terraform-provider-google/issues/4349
func pubsubTopicProjectNotReady(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "retry this operation") {
			log.Printf("[DEBUG] Dismissed error as a retryable operation: %s", err)
			return true, "Waiting for Pubsub topic's project to properly initialize with organiation policy"
		}
	}
	return false, ""
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
func isSqlInternalError(err error) (bool, string) {
	if gerr, ok := err.(*SqlAdminOperationError); ok {
		// SqlAdminOperationError is a non-interface type so we need to cast it through
		// a layer of interface{}.  :)
		var ierr interface{}
		ierr = gerr
		if serr, ok := ierr.(*sqladmin.OperationErrors); ok && serr.Errors[0].Code == "INTERNAL_ERROR" {
			return true, "Received an internal error, which is sometimes retryable for some SQL resources.  Optimistically retrying."
		}

	}
	return false, ""
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
func isSqlOperationInProgressError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
		if strings.Contains(gerr.Body, "you cannot reuse the name of the deleted instance until one week from the deletion date.") {
			return false, ""
		}

		return true, "Waiting for other concurrent Cloud SQL operations to finish"
	}
	return false, ""
}

// Retry if Monitoring operation returns a 429 with a specific message for
// concurrent operations.
func isMonitoringRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && strings.Contains(strings.ToLower(gerr.Body), "too many concurrent edits") {
			return true, "Waiting for other Monitoring changes to finish"
		}
	}
	return false, ""
}

// Retry if App Engine operation returns a 429 with a specific message for
// concurrent operations.
func isAppEngineRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 409 && strings.Contains(strings.ToLower(gerr.Body), "operation is already in progress") {
			return true, "Waiting for other concurrent App Engine changes to finish"
		}
	}
	return false, ""
}

// Retry if KMS CryptoKeyVersions returns a 400 for PENDING_GENERATION
func isCryptoKeyVersionsPendingGeneration(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 400 {
		if strings.Contains(gerr.Body, "PENDING_GENERATION") {
			return true, "Waiting for pending key generation"
		}
	}
	return false, ""
}

// Retry if getting a resource/operation returns a 404 for specific operations.
// opType should describe the operation for which 404 can be retryable.
func isNotFoundRetryableError(opType string) RetryErrorPredicateFunc {
	return func(err error) (bool, string) {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return true, fmt.Sprintf("Retry 404s for %s", opType)
		}
		return false, ""
	}
}

func isStoragePreconditionError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 412 {
		return true, fmt.Sprintf("Retry on storage precondition not met")
	}
	return false, ""
}

func isDataflowJobUpdateRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 404 && strings.Contains(gerr.Body, "in RUNNING OR DRAINING state") {
			return true, "Waiting for job to be in a valid state"
		}
	}
	return false, ""
}
