package google

import (
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"google.golang.org/api/googleapi"
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

func getAllTypes(err error, args ...interface{}) []error {
	var result []error
	for _, v := range args {
		subResult := errwrap.GetAllType(err, v)
		if subResult != nil {
			result = append(result, subResult...)
		}
	}
	return result
}

func isRetryableError(topErr error, customPredicates ...RetryErrorPredicateFunc) bool {
	retryPredicates := append(
		// Global error retry predicates are registered in this default list.
		defaultErrorRetryPredicates,
		customPredicates...)

	// Check all wrapped errors for a retryable error status.
	for _, err := range getAllTypes(topErr, &googleapi.Error{}, &url.Error{}) {
		for _, pred := range retryPredicates {
			if retry, reason := pred(err); retry {
				log.Printf("[DEBUG] Dismissed an error as retryable. %s - %s", reason, err)
				return true
			}
		}
	}
	return false
}
