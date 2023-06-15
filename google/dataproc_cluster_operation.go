// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	tpgdataproc "github.com/hashicorp/terraform-provider-google/google/services/dataproc"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/dataproc/v1"
)

// Deprecated: For backward compatibility dataprocClusterOperationWait is still working,
// but all new code should use DataprocClusterOperationWait in the tpgdataproc package instead.
func dataprocClusterOperationWait(config *transport_tpg.Config, op *dataproc.Operation, activity, userAgent string, timeout time.Duration) error {
	return tpgdataproc.DataprocClusterOperationWait(config, op, activity, userAgent, timeout)
}
