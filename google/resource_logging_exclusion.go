// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import "github.com/hashicorp/terraform-provider-google/google/services/logging"

// parseLoggingExclusionId parses a canonical id into a LoggingExclusionId, or returns an error on failure.
func parseLoggingExclusionId(id string) (*logging.LoggingExclusionId, error) {
	return logging.ParseLoggingExclusionId(id)
}
