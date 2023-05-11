package google

import (
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// We've encountered a few common fingerprint-related strings; if this is one of
// them, we're confident this is an error due to fingerprints.
//
// Deprecated: For backward compatibility FINGERPRINT_FAIL_ERRORS is still working,
// but all new code should use FINGERPRINT_FAIL_ERRORS in the transport_tpg package instead.
var FINGERPRINT_FAIL_ERRORS = transport_tpg.FINGERPRINT_FAIL_ERRORS

// Retry the operation if it looks like a fingerprint mismatch.
//
// Deprecated: For backward compatibility IsFingerprintError is still working,
// but all new code should use IsFingerprintError in the transport_tpg package instead.
func IsFingerprintError(err error) (bool, string) {
	return transport_tpg.IsFingerprintError(err)
}

// If a permission necessary to provision a resource is created in the same config
// as the resource itself, the permission may not have propagated by the time terraform
// attempts to create the resource. This allows those errors to be retried until the timeout expires
//
// Deprecated: For backward compatibility IamMemberMissing is still working,
// but all new code should use IamMemberMissing in the transport_tpg package instead.
func IamMemberMissing(err error) (bool, string) {
	return transport_tpg.IamMemberMissing(err)
}

// Cloud PubSub returns a 400 error if a topic's parent project was recently created and an
// organization policy has not propagated.
// See https://github.com/hashicorp/terraform-provider-google/issues/4349
//
// Deprecated: For backward compatibility PubsubTopicProjectNotReady is still working,
// but all new code should use PubsubTopicProjectNotReady in the transport_tpg package instead.
func PubsubTopicProjectNotReady(err error) (bool, string) {
	return transport_tpg.PubsubTopicProjectNotReady(err)
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
//
// Deprecated: For backward compatibility IsSqlOperationInProgressError is still working,
// but all new code should use IsSqlOperationInProgressError in the transport_tpg package instead.
func IsSqlOperationInProgressError(err error) (bool, string) {
	return transport_tpg.IsSqlOperationInProgressError(err)
}

// Retry if service usage decides you're activating the same service multiple
// times. This can happen if a service and a dependent service aren't batched
// together- eg container.googleapis.com in one request followed by compute.g.c
// in the next (container relies on compute and implicitly activates it)
//
// Deprecated: For backward compatibility ServiceUsageServiceBeingActivated is still working,
// but all new code should use ServiceUsageServiceBeingActivated in the transport_tpg package instead.
func ServiceUsageServiceBeingActivated(err error) (bool, string) {
	return transport_tpg.ServiceUsageServiceBeingActivated(err)
}

// Retry if Bigquery operation returns a 403 with a specific message for
// concurrent operations (which are implemented in terms of 'edit quota').
//
// Deprecated: For backward compatibility IsBigqueryIAMQuotaError is still working,
// but all new code should use IsBigqueryIAMQuotaError in the transport_tpg package instead.
func IsBigqueryIAMQuotaError(err error) (bool, string) {
	return transport_tpg.IsBigqueryIAMQuotaError(err)
}

// Retry if Monitoring operation returns a 409 with a specific message for
// concurrent operations.
//
// Deprecated: For backward compatibility IsMonitoringConcurrentEditError is still working,
// but all new code should use IsMonitoringConcurrentEditError in the transport_tpg package instead.
func IsMonitoringConcurrentEditError(err error) (bool, string) {
	return transport_tpg.IsMonitoringConcurrentEditError(err)
}

// Retry if KMS CryptoKeyVersions returns a 400 for PENDING_GENERATION
//
// Deprecated: For backward compatibility IsCryptoKeyVersionsPendingGeneration is still working,
// but all new code should use IsCryptoKeyVersionsPendingGeneration in the transport_tpg package instead.
func IsCryptoKeyVersionsPendingGeneration(err error) (bool, string) {
	return transport_tpg.IsCryptoKeyVersionsPendingGeneration(err)
}

// Retry if getting a resource/operation returns a 404 for specific operations.
// opType should describe the operation for which 404 can be retryable.
//
// Deprecated: For backward compatibility IsNotFoundRetryableError is still working,
// but all new code should use IsNotFoundRetryableError in the transport_tpg package instead.
func IsNotFoundRetryableError(opType string) transport_tpg.RetryErrorPredicateFunc {
	return transport_tpg.IsNotFoundRetryableError(opType)
}

// Deprecated: For backward compatibility IsDataflowJobUpdateRetryableError is still working,
// but all new code should use IsDataflowJobUpdateRetryableError in the transport_tpg package instead.
func IsDataflowJobUpdateRetryableError(err error) (bool, string) {
	return transport_tpg.IsDataflowJobUpdateRetryableError(err)
}

// Deprecated: For backward compatibility IsPeeringOperationInProgress is still working,
// but all new code should use IsPeeringOperationInProgress in the transport_tpg package instead.
func IsPeeringOperationInProgress(err error) (bool, string) {
	return transport_tpg.IsPeeringOperationInProgress(err)
}

// Deprecated: For backward compatibility DatastoreIndex409Contention is still working,
// but all new code should use DatastoreIndex409Contention in the transport_tpg package instead.
func DatastoreIndex409Contention(err error) (bool, string) {
	return transport_tpg.DatastoreIndex409Contention(err)
}

// Deprecated: For backward compatibility IapClient409Operation is still working,
// but all new code should use IapClient409Operation in the transport_tpg package instead.
func IapClient409Operation(err error) (bool, string) {
	return transport_tpg.IapClient409Operation(err)
}

// Deprecated: For backward compatibility HealthcareDatasetNotInitialized is still working,
// but all new code should use HealthcareDatasetNotInitialized in the transport_tpg package instead.
func HealthcareDatasetNotInitialized(err error) (bool, string) {
	return transport_tpg.HealthcareDatasetNotInitialized(err)
}

// Cloud Run APIs may return a 409 on create to indicate that a resource has been deleted in the foreground
// (eg GET and LIST) but not the backing apiserver. When we encounter a 409, we can retry it.
// Note that due to limitations in MMv1's error_retry_predicates this is currently applied to all requests.
// We only expect to receive it on create, though.
//
// Deprecated: For backward compatibility IsCloudRunCreationConflict is still working,
// but all new code should use IsCloudRunCreationConflict in the transport_tpg package instead.
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
//
// Deprecated: For backward compatibility IamServiceAccountNotFound is still working,
// but all new code should use IamServiceAccountNotFound in the transport_tpg package instead.
func IamServiceAccountNotFound(err error) (bool, string) {
	return transport_tpg.IamServiceAccountNotFound(err)
}

// Concurrent Apigee operations can fail with a 400 error
//
// Deprecated: For backward compatibility IsApigeeRetryableError is still working,
// but all new code should use IsApigeeRetryableError in the transport_tpg package instead.
func IsApigeeRetryableError(err error) (bool, string) {
	return transport_tpg.IsApigeeRetryableError(err)
}
