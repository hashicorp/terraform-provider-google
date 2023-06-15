// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/deploymentmanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility DeploymentManagerOperationWaitTime is still working,
// but all new code should use DeploymentManagerOperationWaitTime in the deploymentmanager package instead.
func DeploymentManagerOperationWaitTime(config *transport_tpg.Config, resp interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return deploymentmanager.DeploymentManagerOperationWaitTime(config, resp, project, activity, userAgent, timeout)
}
