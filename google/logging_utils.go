// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/logging"
)

// loggingSinkResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingSinkResourceTypes = logging.LoggingSinkResourceTypes

type LoggingSinkId = logging.LoggingSinkId

// parseLoggingSinkId parses a canonical id into a LoggingSinkId, or returns an error on failure.
func parseLoggingSinkId(id string) (*LoggingSinkId, error) {
	return logging.ParseLoggingSinkId(id)
}
