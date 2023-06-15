// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/tags"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility TagsLocationOperationWaitTimeWithResponse is still working,
// but all new code should use TagsLocationOperationWaitTimeWithResponse in the tags package instead.
func TagsLocationOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	return tags.TagsLocationOperationWaitTimeWithResponse(config, op, response, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility TagsLocationOperationWaitTime is still working,
// but all new code should use TagsLocationOperationWaitTime in the tags package instead.
func TagsLocationOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	return tags.TagsLocationOperationWaitTime(config, op, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility GetLocationFromOpName is still working,
// but all new code should use GetLocationFromOpName in the tags package instead.
func GetLocationFromOpName(opName string) string {
	return tags.GetLocationFromOpName(opName)
}
