package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

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

// If a permission necessary to provision a resource is created in the same config
// as the resource itself, the permission may not have propagated by the time terraform
// attempts to create the resource. This allows those errors to be retried until the timeout expires
func IamMemberMissing(err error) (bool, string) {
	return transport_tpg.IamMemberMissing(err)
}

// Cloud PubSub returns a 400 error if a topic's parent project was recently created and an
// organization policy has not propagated.
// See https://github.com/hashicorp/terraform-provider-google/issues/4349
func PubsubTopicProjectNotReady(err error) (bool, string) {
	return transport_tpg.PubsubTopicProjectNotReady(err)
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
func IsSqlInternalError(err error) (bool, string) {
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
func IsNotFoundRetryableError(opType string) transport_tpg.RetryErrorPredicateFunc {
	return func(err error) (bool, string) {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return true, fmt.Sprintf("Retry 404s for %s", opType)
		}
		return false, ""
	}
}

func IsDataflowJobUpdateRetryableError(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 404 && strings.Contains(gerr.Body, "in RUNNING OR DRAINING state") {
			return true, "Waiting for job to be in a valid state"
		}
	}
	return false, ""
}

func IsPeeringOperationInProgress(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "There is a peering operation in progress") {
			return true, "Waiting peering operation to complete"
		}
	}
	return false, ""
}

func IsCloudFunctionsSourceCodeError(err error) (bool, string) {
	if operr, ok := err.(*CommonOpError); ok {
		if operr.Code == 3 && operr.Message == "Failed to retrieve function source code" {
			return true, fmt.Sprintf("Retry on Function failing to pull code from GCS")
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
