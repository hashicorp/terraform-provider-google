// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/containerattached"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// nolint: deadcode,unused
//
// Deprecated: For backward compatibility ContainerAttachedOperationWaitTimeWithResponse is still working,
// but all new code should use ContainerAttachedOperationWaitTimeWithResponse in the containerattached package instead.
func ContainerAttachedOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return containerattached.ContainerAttachedOperationWaitTimeWithResponse(config, op, response, project, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility ContainerAttachedOperationWaitTime is still working,
// but all new code should use ContainerAttachedOperationWaitTime in the containerattached package instead.
func ContainerAttachedOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return containerattached.ContainerAttachedOperationWaitTime(config, op, project, activity, userAgent, timeout)
}
