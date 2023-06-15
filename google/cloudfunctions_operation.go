// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	tpgcloudfunctions "github.com/hashicorp/terraform-provider-google/google/services/cloudfunctions"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudfunctions/v1"
)

// Deprecated: For backward compatibility cloudFunctionsOperationWait is still working,
// but all new code should use CloudFunctionsOperationWait in the tpgcloudfunctions package instead.
func cloudFunctionsOperationWait(config *transport_tpg.Config, op *cloudfunctions.Operation, activity, userAgent string, timeout time.Duration) error {
	return tpgcloudfunctions.CloudFunctionsOperationWait(config, op, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility IsCloudFunctionsSourceCodeError is still working,
// but all new code should use IsCloudFunctionsSourceCodeError in the tpgcloudfunctions package instead.
func IsCloudFunctionsSourceCodeError(err error) (bool, string) {
	return tpgcloudfunctions.IsCloudFunctionsSourceCodeError(err)
}
