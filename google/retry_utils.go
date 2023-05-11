package google

import (
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility retry is still working,
// but all new code should use Retry in the transport_tpg package instead.
func retry(retryFunc func() error) error {
	return transport_tpg.Retry(retryFunc)
}

// Deprecated: For backward compatibility retryTime is still working,
// but all new code should use RetryTime in the transport_tpg package instead.
func retryTime(retryFunc func() error, minutes int) error {
	return transport_tpg.RetryTime(retryFunc, minutes)
}

// Deprecated: For backward compatibility RetryTimeDuration is still working,
// but all new code should use RetryTimeDuration in the transport_tpg package instead.
func RetryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) error {
	return transport_tpg.RetryTimeDuration(retryFunc, duration, errorRetryPredicates...)
}

// Deprecated: For backward compatibility isRetryableError is still working,
// but all new code should use IsRetryableError in the transport_tpg package instead.
func isRetryableError(topErr error, customPredicates ...transport_tpg.RetryErrorPredicateFunc) bool {
	return transport_tpg.IsRetryableError(topErr, customPredicates...)
}

// The polling overrides the default backoff logic with max backoff of 10s. The poll interval can be greater than 10s.
//
// Deprecated: For backward compatibility retryWithPolling is still working,
// but all new code should use RetryWithPolling in the transport_tpg package instead.
func retryWithPolling(retryFunc func() (interface{}, error), timeout time.Duration, pollInterval time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (interface{}, error) {
	return transport_tpg.RetryWithPolling(retryFunc, timeout, pollInterval, errorRetryPredicates...)
}
