// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

func TestRetry(t *testing.T) {
	i := 0
	f := func() error {
		i++
		return &googleapi.Error{
			Code: 500,
		}
	}
	if err := Retry(RetryOptions{
		RetryFunc: f,
		Timeout:   time.Duration(1000) * time.Millisecond,
	}); err == nil || err.(*googleapi.Error).Code != 500 {
		t.Errorf("unexpected error retrying: %v", err)
	}
	if i < 2 {
		t.Errorf("expected error function to be called at least twice, but was called %d times", i)
	}
}

func TestRetry_wrapped(t *testing.T) {
	i := 0
	f := func() error {
		i++
		err := &googleapi.Error{
			Code: 500,
		}
		return errwrap.Wrapf("nested error: {{err}}", err)
	}
	if err := Retry(RetryOptions{
		RetryFunc: f,
		Timeout:   time.Duration(1000) * time.Millisecond,
	}); err == nil {
		t.Errorf("unexpected nil error, expected an error")
	} else {
		innerErr := errwrap.GetType(err, &googleapi.Error{})
		if innerErr == nil {
			t.Errorf("unexpected error %v does not have a google api error", err)
		}
		gerr := innerErr.(*googleapi.Error)
		if gerr.Code != 500 {
			t.Errorf("unexpected googleapi error expected code 500, error: %v", gerr)
		}
	}
	if i < 2 {
		t.Errorf("expected error function to be called at least twice, but was called %d times", i)
	}
}

func TestRetry_noretry(t *testing.T) {
	i := 0
	f := func() error {
		i++
		return &googleapi.Error{
			Code: 400,
		}
	}
	if err := Retry(RetryOptions{
		RetryFunc: f,
		Timeout:   time.Duration(1000) * time.Millisecond,
	}); err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error retrying: %v", err)
	}
	if i != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", i)
	}
}

func TestRetry_URLTimeoutsShouldRetry(t *testing.T) {
	runCount := 0
	retryFunc := func() error {
		runCount++
		if runCount == 1 {
			return &url.Error{
				Err: TimeoutErr,
			}
		}
		return nil
	}
	err := Retry(RetryOptions{
		RetryFunc: retryFunc,
		Timeout:   1 * time.Minute,
	})
	if err != nil {
		t.Errorf("unexpected error: got '%v' want 'nil'", err)
	}
	expectedRunCount := 2
	if runCount != expectedRunCount {
		t.Errorf("expected the retryFunc to be called %v time(s), instead was called %v time(s)", expectedRunCount, runCount)
	}
}

func TestRetryWithPolling_noRetry(t *testing.T) {
	retryCount := 0
	retryFunc := func() error {
		retryCount++
		return &googleapi.Error{
			Code: 400,
		}
	}
	err := Retry(RetryOptions{
		RetryFunc:    retryFunc,
		Timeout:      time.Duration(1000) * time.Millisecond,
		PollInterval: time.Duration(100) * time.Millisecond,
	})
	if err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error %v", err)
	}
	if retryCount != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_notRetryable(t *testing.T) {
	retryCount := 0
	retryFunc := func() error {
		retryCount++
		return &googleapi.Error{
			Code: 400,
		}
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	err := Retry(RetryOptions{
		RetryFunc:            retryFunc,
		Timeout:              time.Duration(1000) * time.Millisecond,
		PollInterval:         time.Duration(100) * time.Millisecond,
		ErrorRetryPredicates: []RetryErrorPredicateFunc{isRetryableFunc},
	})
	if err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error %v", err)
	}
	if retryCount != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_retriedAndSucceeded(t *testing.T) {
	retryCount := 0
	// Retry once and succeeds.
	retryFunc := func() error {
		retryCount++
		// Error code of 200 is retryable.
		if retryCount < 2 {
			return &googleapi.Error{
				Code: 200,
			}
		}
		return nil
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	err := Retry(RetryOptions{
		RetryFunc:            retryFunc,
		Timeout:              time.Duration(1000) * time.Millisecond,
		PollInterval:         time.Duration(100) * time.Millisecond,
		ErrorRetryPredicates: []RetryErrorPredicateFunc{isRetryableFunc},
	})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if retryCount != 2 {
		t.Errorf("expected error function to be called exactly twice, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_retriedAndFailed(t *testing.T) {
	retryCount := 0
	// Retry once and fails.
	retryFunc := func() error {
		retryCount++
		// Error code of 200 is retryable.
		if retryCount < 2 {
			return &googleapi.Error{
				Code: 200,
			}
		}
		return &googleapi.Error{
			Code: 400,
		}
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	err := Retry(RetryOptions{
		RetryFunc:            retryFunc,
		Timeout:              time.Duration(1000) * time.Millisecond,
		PollInterval:         time.Duration(100) * time.Millisecond,
		ErrorRetryPredicates: []RetryErrorPredicateFunc{isRetryableFunc},
	})
	if err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error %v", err)
	}
	if retryCount != 2 {
		t.Errorf("expected error function to be called exactly twice, but was called %d times", retryCount)
	}
}
