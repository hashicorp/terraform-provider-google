// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/dialogflowcx"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// nolint: deadcode,unused
//
// Deprecated: For backward compatibility DialogflowCXOperationWaitTimeWithResponse is still working,
// but all new code should use DialogflowCXOperationWaitTimeWithResponse in the dialogflowcx package instead.
func DialogflowCXOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, activity, userAgent, location string, timeout time.Duration) error {
	return dialogflowcx.DialogflowCXOperationWaitTimeWithResponse(config, op, response, activity, userAgent, location, timeout)
}

// Deprecated: For backward compatibility DialogflowCXOperationWaitTime is still working,
// but all new code should use DialogflowCXOperationWaitTime in the dialogflowcx package instead.
func DialogflowCXOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, activity, userAgent, location string, timeout time.Duration) error {
	return dialogflowcx.DialogflowCXOperationWaitTime(config, op, activity, userAgent, location, timeout)
}
