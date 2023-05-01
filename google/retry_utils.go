package google

import (
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func retry(retryFunc func() error) error {
	return transport_tpg.Retry(retryFunc)
}

func retryTime(retryFunc func() error, minutes int) error {
	return transport_tpg.RetryTime(retryFunc, minutes)
}

func RetryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) error {
	return transport_tpg.RetryTimeDuration(retryFunc, duration, errorRetryPredicates...)
}

func isRetryableError(topErr error, customPredicates ...transport_tpg.RetryErrorPredicateFunc) bool {
	return transport_tpg.IsRetryableError(topErr, customPredicates...)
}

// The polling overrides the default backoff logic with max backoff of 10s. The poll interval can be greater than 10s.
func retryWithPolling(retryFunc func() (interface{}, error), timeout time.Duration, pollInterval time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (interface{}, error) {
	return transport_tpg.RetryWithPolling(retryFunc, timeout, pollInterval, errorRetryPredicates...)
}
