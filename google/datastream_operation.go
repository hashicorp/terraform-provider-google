// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/datastream"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// nolint: deadcode,unused
//
// Deprecated: For backward compatibility DatastreamOperationWaitTimeWithResponse is still working,
// but all new code should use DatastreamOperationWaitTimeWithResponse in the datastream package instead.
func DatastreamOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return datastream.DatastreamOperationWaitTimeWithResponse(config, op, response, project, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility DatastreamOperationWaitTime is still working,
// but all new code should use DatastreamOperationWaitTime in the datastream package instead.
func DatastreamOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return datastream.DatastreamOperationWaitTime(config, op, project, activity, userAgent, timeout)
}
