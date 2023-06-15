// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// nolint: deadcode,unused
func ResourceManagerOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	return resourcemanager.ResourceManagerOperationWaitTimeWithResponse(config, op, response, activity, userAgent, timeout)
}

func ResourceManagerOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	return resourcemanager.ResourceManagerOperationWaitTime(config, op, activity, userAgent, timeout)
}
