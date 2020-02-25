package google

import (
	"log"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
