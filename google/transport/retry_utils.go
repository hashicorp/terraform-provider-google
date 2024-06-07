// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"log"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type RetryOptions struct {
	RetryFunc            func() error
	Timeout              time.Duration
	PollInterval         time.Duration
	ErrorRetryPredicates []RetryErrorPredicateFunc
	ErrorAbortPredicates []RetryErrorPredicateFunc
}

func Retry(opt RetryOptions) error {
	if opt.Timeout == 0 {
		opt.Timeout = 1 * time.Minute
	}

	if opt.PollInterval != 0 {
		refreshFunc := func() (interface{}, string, error) {
			err := opt.RetryFunc()
			if err == nil {
				return "", "done", nil
			}

			// Check if it is a retryable error.
			if IsRetryableError(err, opt.ErrorRetryPredicates, opt.ErrorAbortPredicates) {
				return "", "retrying", nil
			}

			// The error is not retryable.
			return "", "done", err
		}
		stateChange := &retry.StateChangeConf{
			Pending: []string{
				"retrying",
			},
			Target: []string{
				"done",
			},
			Refresh:      refreshFunc,
			Timeout:      opt.Timeout,
			PollInterval: opt.PollInterval,
		}

		_, err := stateChange.WaitForState()
		return err
	}

	return retry.Retry(opt.Timeout, func() *retry.RetryError {
		err := opt.RetryFunc()
		if err == nil {
			return nil
		}
		if IsRetryableError(err, opt.ErrorRetryPredicates, opt.ErrorAbortPredicates) {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
}

func IsRetryableError(topErr error, retryPredicates, abortPredicates []RetryErrorPredicateFunc) bool {
	if topErr == nil {
		return false
	}

	retryPredicates = append(
		// Global error retry predicates are registered in this default list.
		defaultErrorRetryPredicates,
		retryPredicates...)

	// Check all wrapped errors for an abortable error status.
	isAbortable := false
	errwrap.Walk(topErr, func(werr error) {
		for _, pred := range abortPredicates {
			if predAbort, predReason := pred(werr); predAbort {
				log.Printf("[DEBUG] Dismissed an error as abortable. %s - %s", predReason, werr)
				isAbortable = true
				return
			}
		}
	})
	if isAbortable {
		return false
	}

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
