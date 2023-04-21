package transport

import (
	"log"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func RetryTimeDuration(retryFunc func() error, duration time.Duration, errorRetryPredicates ...RetryErrorPredicateFunc) error {
	return resource.Retry(duration, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if IsRetryableError(err, errorRetryPredicates...) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
}

func IsRetryableError(topErr error, customPredicates ...RetryErrorPredicateFunc) bool {
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
