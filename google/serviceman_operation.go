// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	tpgservicemanagement "github.com/hashicorp/terraform-provider-google/google/services/servicemanagement"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicemanagement/v1"
)

// Deprecated: For backward compatibility ServiceManagementOperationWaitTime is still working,
// but all new code should use ServiceManagementOperationWaitTime in the tpgservicemanagement package instead.
func ServiceManagementOperationWaitTime(config *transport_tpg.Config, op *servicemanagement.Operation, activity, userAgent string, timeout time.Duration) (googleapi.RawMessage, error) {
	return tpgservicemanagement.ServiceManagementOperationWaitTime(config, op, activity, userAgent, timeout)
}
