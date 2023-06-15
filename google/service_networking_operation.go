// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	tpgservicenetworking "github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/servicenetworking/v1"
)

// Deprecated: For backward compatibility ServiceNetworkingOperationWaitTime is still working,
// but all new code should use ServiceNetworkingOperationWaitTime in the tpgservicenetworking package instead.
func ServiceNetworkingOperationWaitTime(config *transport_tpg.Config, op *servicenetworking.Operation, activity, userAgent, project string, timeout time.Duration) error {
	return tpgservicenetworking.ServiceNetworkingOperationWaitTime(config, op, activity, userAgent, project, timeout)
}
