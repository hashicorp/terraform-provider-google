package google

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func PollingWaitTime(pollF PollReadFunc, checkResponse PollCheckResponseFunc, activity string,
	timeout time.Duration, targetOccurrences int) error {
	log.Printf("[DEBUG] %s: Polling until expected state is read", activity)
	log.Printf("[DEBUG] Target occurrences: %d", targetOccurrences)
	if targetOccurrences == 1 {
		return resource.Retry(timeout, func() *resource.RetryError {
			readResp, readErr := pollF()
			return checkResponse(readResp, readErr)
		})
	}
	return RetryWithTargetOccurrences(timeout, targetOccurrences, func() *resource.RetryError {
		readResp, readErr := pollF()
		return checkResponse(readResp, readErr)
	})
}

// RetryWithTargetOccurrences is a basic wrapper around StateChangeConf that will retry
// a function until it returns the specified amount of target occurrences continuously.
// Adapted from the Retry function in the go SDK.
func RetryWithTargetOccurrences(timeout time.Duration, targetOccurrences int,
	f resource.RetryFunc) error {
	// These are used to pull the error out of the function; need a mutex to
	// avoid a data race.
	var resultErr error
	var resultErrMu sync.Mutex

	c := &resource.StateChangeConf{
		Pending:                   []string{"retryableerror"},
		Target:                    []string{"success"},
		Timeout:                   timeout,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: targetOccurrences,
		Refresh: func() (interface{}, string, error) {
			rerr := f()

			resultErrMu.Lock()
			defer resultErrMu.Unlock()

			if rerr == nil {
				resultErr = nil
				return 42, "success", nil
			}

			resultErr = rerr.Err

			if rerr.Retryable {
				return 42, "retryableerror", nil
			}
			return nil, "quit", rerr.Err
		},
	}

	_, waitErr := c.WaitForState()

	// Need to acquire the lock here to be able to avoid race using resultErr as
	// the return value
	resultErrMu.Lock()
	defer resultErrMu.Unlock()

	// resultErr may be nil because the wait timed out and resultErr was never
	// set; this is still an error
	if resultErr == nil {
		return waitErr
	}
	// resultErr takes precedence over waitErr if both are set because it is
	// more likely to be useful
	return resultErr
}

/**
 * Common PollCheckResponseFunc implementations
 */

// PollCheckForExistence waits for a successful response, continues polling on 404,
// and returns any other error.
func PollCheckForExistence(_ map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		if isGoogleApiErrorWithCode(respErr, 404) {
			return PendingStatusPollResult("not found")
		}
		return ErrorPollResult(respErr)
	}
	return SuccessPollResult()
}

// PollCheckForExistenceWith403 waits for a successful response, continues polling on 404 or 403,
// and returns any other error.
func PollCheckForExistenceWith403(_ map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		if isGoogleApiErrorWithCode(respErr, 404) || isGoogleApiErrorWithCode(respErr, 403) {
			return PendingStatusPollResult("not found")
		}
		return ErrorPollResult(respErr)
	}
	return SuccessPollResult()
}

// PollCheckForAbsence waits for a 404/403 response, continues polling on a successful
// response, and returns any other error.
func PollCheckForAbsenceWith403(_ map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		if isGoogleApiErrorWithCode(respErr, 404) || isGoogleApiErrorWithCode(respErr, 403) {
			return SuccessPollResult()
		}
		return ErrorPollResult(respErr)
	}
	return PendingStatusPollResult("found")
}

// PollCheckForAbsence waits for a 404 response, continues polling on a successful
// response, and returns any other error.
func PollCheckForAbsence(_ map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		if isGoogleApiErrorWithCode(respErr, 404) {
			return SuccessPollResult()
		}
		return ErrorPollResult(respErr)
	}
	return PendingStatusPollResult("found")
}
