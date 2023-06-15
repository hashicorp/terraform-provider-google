// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	tpgcomposer "github.com/hashicorp/terraform-provider-google/google/services/composer"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/composer/v1"
)

// Deprecated: For backward compatibility ComposerOperationWaitTime is still working,
// but all new code should use ComposerOperationWaitTime in the tpgcomposer package instead.
func ComposerOperationWaitTime(config *transport_tpg.Config, op *composer.Operation, project, activity, userAgent string, timeout time.Duration) error {
	return tpgcomposer.ComposerOperationWaitTime(config, op, project, activity, userAgent, timeout)
}
