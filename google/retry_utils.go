package google

import (
	"log"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return retryTimeDuration(retryFunc, time.Duration(minutes)*time.Minute)
}

func retryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...RetryErrorPredicateFunc) error {
	return resource.Retry(duration, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if isRetryableError(err, errorRetryPredicates...) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
}

func isRetryableError(topErr error, customPredicates ...RetryErrorPredicateFunc) bool {
	if topErr == nil {
		return false
	}

	retryPredicates := append(
		// Global error retry predicates are registered in this default list.
		defaultErrorRetryPredicates,
		customPredicates...)

	// Check all wrapped errors for a retryable error status.
	isRetryable := false
	errwrap.Walk(topErr, func(werr error) {
		for _, pred := range retryPredicates {
			if predRetry, predReason := pred(werr); predRetry {
				log.Printf("[DEBUG] Dismissed an error as retryable. %s - %s", predReason, werr)
				isRetryable = true
				return
			}
		}
	})
	return isRetryable
}

// The polling overrides the default backoff logic with max backoff of 10s. The poll interval can be greater than 10s.
func retryWithPolling(retryFunc func() (interface{}, error), timeout time.Duration, pollInterval time.Duration, errorRetryPredicates ...RetryErrorPredicateFunc) (interface{}, error) {
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
