// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/cloudrunv2"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/run/v2"
)

// Deprecated: For backward compatibility runAdminV2OperationWaitTimeWithResponse is still working,
// but all new code should use RunAdminV2OperationWaitTimeWithResponse in the cloudrunv2 package instead.
func runAdminV2OperationWaitTimeWithResponse(config *transport_tpg.Config, op *run.GoogleLongrunningOperation, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return cloudrunv2.RunAdminV2OperationWaitTimeWithResponse(config, op, response, project, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility runAdminV2OperationWaitTime is still working,
// but all new code should use RunAdminV2OperationWaitTime in the cloudrunv2 package instead.
func runAdminV2OperationWaitTime(config *transport_tpg.Config, op *run.GoogleLongrunningOperation, project, activity, userAgent string, timeout time.Duration) error {
	return cloudrunv2.RunAdminV2OperationWaitTime(config, op, project, activity, userAgent, timeout)
}
