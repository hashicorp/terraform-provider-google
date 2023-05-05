package google

import (
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// We've encountered a few common fingerprint-related strings; if this is one of
// them, we're confident this is an error due to fingerprints.
var FINGERPRINT_FAIL_ERRORS = transport_tpg.FINGERPRINT_FAIL_ERRORS

// Retry the operation if it looks like a fingerprint mismatch.
func IsFingerprintError(err error) (bool, string) {
	return transport_tpg.IsFingerprintError(err)
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
func IsSqlOperationInProgressError(err error) (bool, string) {
	return transport_tpg.IsSqlOperationInProgressError(err)
}

// Retry if service usage decides you're activating the same service multiple
// times. This can happen if a service and a dependent service aren't batched
// together- eg container.googleapis.com in one request followed by compute.g.c
// in the next (container relies on compute and implicitly activates it)
func ServiceUsageServiceBeingActivated(err error) (bool, string) {
	return transport_tpg.ServiceUsageServiceBeingActivated(err)
}

// Retry if Bigquery operation returns a 403 with a specific message for
// concurrent operations (which are implemented in terms of 'edit quota').
func IsBigqueryIAMQuotaError(err error) (bool, string) {
	return transport_tpg.IsBigqueryIAMQuotaError(err)
}

// Retry if Monitoring operation returns a 409 with a specific message for
// concurrent operations.
func IsMonitoringConcurrentEditError(err error) (bool, string) {
	return transport_tpg.IsMonitoringConcurrentEditError(err)
}

// Retry if KMS CryptoKeyVersions returns a 400 for PENDING_GENERATION
func IsCryptoKeyVersionsPendingGeneration(err error) (bool, string) {
	return transport_tpg.IsCryptoKeyVersionsPendingGeneration(err)
}

// Retry if getting a resource/operation returns a 404 for specific operations.
// opType should describe the operation for which 404 can be retryable.
func IsNotFoundRetryableError(opType string) transport_tpg.RetryErrorPredicateFunc {
	return transport_tpg.IsNotFoundRetryableError(opType)
}

func IsDataflowJobUpdateRetryableError(err error) (bool, string) {
	return transport_tpg.IsDataflowJobUpdateRetryableError(err)
}

func IsPeeringOperationInProgress(err error) (bool, string) {
	return transport_tpg.IsPeeringOperationInProgress(err)
}

func DatastoreIndex409Contention(err error) (bool, string) {
	return transport_tpg.DatastoreIndex409Contention(err)
}

func IapClient409Operation(err error) (bool, string) {
	return transport_tpg.IapClient409Operation(err)
}

func HealthcareDatasetNotInitialized(err error) (bool, string) {
	return transport_tpg.HealthcareDatasetNotInitialized(err)
}

// Cloud Run APIs may return a 409 on create to indicate that a resource has been deleted in the foreground
// (eg GET and LIST) but not the backing apiserver. When we encounter a 409, we can retry it.
// Note that due to limitations in MMv1's error_retry_predicates this is currently applied to all requests.
// We only expect to receive it on create, though.
func IsCloudRunCreationConflict(err error) (bool, string) {
	return transport_tpg.IsCloudRunCreationConflict(err)
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
	return transport_tpg.IamServiceAccountNotFound(err)
}

// Concurrent Apigee operations can fail with a 400 error
func IsApigeeRetryableError(err error) (bool, string) {
	return transport_tpg.IsApigeeRetryableError(err)
}
