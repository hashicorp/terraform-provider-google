package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type (
	// Function handling for polling for a resource
	PollReadFunc func() (resp map[string]interface{}, respErr error)

	// Function to check the response from polling once
	PollCheckResponseFunc func(resp map[string]interface{}, respErr error) PollResult

	PollResult *resource.RetryError
)

// Helper functions to construct result of single pollRead as return result for a PollCheckResponseFunc
func ErrorPollResult(err error) PollResult {
	return resource.NonRetryableError(err)
}

func PendingStatusPollResult(status string) PollResult {
	return resource.RetryableError(fmt.Errorf("got pending status %q", status))
}

func SuccessPollResult() PollResult {
	return nil
}

func PollingWaitTime(pollF PollReadFunc, checkResponse PollCheckResponseFunc, activity string, timeout time.Duration) error {
	log.Printf("[DEBUG] %s: Polling until expected state is read", activity)
	return resource.Retry(timeout, func() *resource.RetryError {
		readResp, readErr := pollF()
		return checkResponse(readResp, readErr)
	})
}

/**
 * Common PollCheckResponseFunc implementations
 */

// PollCheckForExistence waits for a successful response, continues polling on 404, and returns any other error.
func PollCheckForExistence(_ map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		if isGoogleApiErrorWithCode(respErr, 404) {
			return PendingStatusPollResult("not found")
		}
		return ErrorPollResult(respErr)
	}
	return SuccessPollResult()
}
