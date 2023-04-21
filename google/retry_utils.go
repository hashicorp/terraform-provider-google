package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return RetryTimeDuration(retryFunc, time.Duration(minutes)*time.Minute)
}

func RetryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) error {
	return transport_tpg.RetryTimeDuration(retryFunc, duration, errorRetryPredicates...)
}

func isRetryableError(topErr error, customPredicates ...transport_tpg.RetryErrorPredicateFunc) bool {
	return transport_tpg.IsRetryableError(topErr, customPredicates...)
}

// The polling overrides the default backoff logic with max backoff of 10s. The poll interval can be greater than 10s.
func retryWithPolling(retryFunc func() (interface{}, error), timeout time.Duration, pollInterval time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		result, err := retryFunc()
		if err == nil {
			return result, "done", nil
		}

		// Check if it is a retryable error.
		if isRetryableError(err, errorRetryPredicates...) {
			return result, "retrying", nil
		}

		// The error is not retryable.
		return result, "done", err
	}
	stateChange := &resource.StateChangeConf{
		Pending: []string{
			"retrying",
		},
		Target: []string{
			"done",
		},
		Refresh:      refreshFunc,
		Timeout:      timeout,
		PollInterval: pollInterval,
	}

	return stateChange.WaitForState()
}
