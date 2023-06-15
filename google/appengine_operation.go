// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/appengine"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility AppEngineOperationWaitTimeWithResponse is still working,
// but all new code should use AppEngineOperationWaitTimeWithResponse in the appengine package instead.
func AppEngineOperationWaitTimeWithResponse(config *transport_tpg.Config, res interface{}, response *map[string]interface{}, appId, activity, userAgent string, timeout time.Duration) error {
	return appengine.AppEngineOperationWaitTimeWithResponse(config, res, response, appId, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility AppEngineOperationWaitTime is still working,
// but all new code should use AppEngineOperationWaitTime in the appengine package instead.
func AppEngineOperationWaitTime(config *transport_tpg.Config, res interface{}, appId, activity, userAgent string, timeout time.Duration) error {
	return appengine.AppEngineOperationWaitTime(config, res, appId, activity, userAgent, timeout)
}
