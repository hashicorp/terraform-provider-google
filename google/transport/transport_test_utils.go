// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

type TimeoutError struct {
	timeout bool
}

func (e *TimeoutError) Timeout() bool {
	return e.timeout
}

func (e *TimeoutError) Error() string {
	return "timeout error"
}

var TimeoutErr = &TimeoutError{timeout: true}
