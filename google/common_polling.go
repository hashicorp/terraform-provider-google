// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Helper functions to construct result of single pollRead as return result for a PollCheckResponseFunc
//
// Deprecated: For backward compatibility ErrorPollResult is still working,
// but all new code should use ErrorPollResult in the transport_tpg package instead.
func ErrorPollResult(err error) transport_tpg.PollResult {
	return transport_tpg.ErrorPollResult(err)
}

// Deprecated: For backward compatibility PendingStatusPollResult is still working,
// but all new code should use PendingStatusPollResult in the transport_tpg package instead.
func PendingStatusPollResult(status string) transport_tpg.PollResult {
	return transport_tpg.PendingStatusPollResult(status)
}

// Deprecated: For backward compatibility SuccessPollResult is still working,
// but all new code should use SuccessPollResult in the transport_tpg package instead.
func SuccessPollResult() transport_tpg.PollResult {
	return transport_tpg.SuccessPollResult()
}

// Deprecated: For backward compatibility PollingWaitTime is still working,
// but all new code should use PollingWaitTime in the transport_tpg package instead.
func PollingWaitTime(pollF transport_tpg.PollReadFunc, checkResponse transport_tpg.PollCheckResponseFunc, activity string,
	timeout time.Duration, targetOccurrences int) error {
	return transport_tpg.PollingWaitTime(pollF, checkResponse, activity, timeout, targetOccurrences)
}

// RetryWithTargetOccurrences is a basic wrapper around StateChangeConf that will retry
// a function until it returns the specified amount of target occurrences continuously.
// Adapted from the Retry function in the go SDK.
//
// Deprecated: For backward compatibility RetryWithTargetOccurrences is still working,
// but all new code should use RetryWithTargetOccurrences in the transport_tpg package instead.
func RetryWithTargetOccurrences(timeout time.Duration, targetOccurrences int,
	f resource.RetryFunc) error {
	return transport_tpg.RetryWithTargetOccurrences(timeout, targetOccurrences, f)
}

/**
 * Common PollCheckResponseFunc implementations
 */

// PollCheckForExistence waits for a successful response, continues polling on 404,
// and returns any other error.
//
// Deprecated: For backward compatibility PollCheckForExistence is still working,
// but all new code should use PollCheckForExistence in the transport_tpg package instead.
func PollCheckForExistence(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForExistence(nil, respErr)
}

// PollCheckForExistenceWith403 waits for a successful response, continues polling on 404 or 403,
// and returns any other error.
//
// Deprecated: For backward compatibility PollCheckForExistenceWith403 is still working,
// but all new code should use PollCheckForExistenceWith403 in the transport_tpg package instead.
func PollCheckForExistenceWith403(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForExistenceWith403(nil, respErr)
}

// PollCheckForAbsence waits for a 404/403 response, continues polling on a successful
// response, and returns any other error.
//
// Deprecated: For backward compatibility PollCheckForAbsenceWith403 is still working,
// but all new code should use PollCheckForAbsenceWith403 in the transport_tpg package instead.
func PollCheckForAbsenceWith403(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForAbsenceWith403(nil, respErr)
}

// PollCheckForAbsence waits for a 404 response, continues polling on a successful
// response, and returns any other error.
//
// Deprecated: For backward compatibility PollCheckForAbsence is still working,
// but all new code should use PollCheckForAbsence in the transport_tpg package instead.
func PollCheckForAbsence(_ map[string]interface{}, respErr error) transport_tpg.PollResult {
	return transport_tpg.PollCheckForAbsence(nil, respErr)
}
