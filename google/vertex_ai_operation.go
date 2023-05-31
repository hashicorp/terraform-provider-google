// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// nolint: deadcode,unused
//
// Deprecated: For backward compatibility VertexAIOperationWaitTimeWithResponse is still working,
// but all new code should use VertexAIOperationWaitTimeWithResponse in the vertexai package instead.
func VertexAIOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return vertexai.VertexAIOperationWaitTimeWithResponse(config, op, response, project, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility VertexAIOperationWaitTime is still working,
// but all new code should use VertexAIOperationWaitTime in the vertexai package instead.
func VertexAIOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	return vertexai.VertexAIOperationWaitTime(config, op, project, activity, userAgent, timeout)
}
