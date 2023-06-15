// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/dataproc"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility dataprocJobOperationWait is still working,
// but all new code should use DataprocJobOperationWait in the dataproc package instead.
func dataprocJobOperationWait(config *transport_tpg.Config, region, projectId, jobId, activity, userAgent string, timeout time.Duration) error {
	return dataproc.DataprocJobOperationWait(config, region, projectId, jobId, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility dataprocDeleteOperationWait is still working,
// but all new code should use DataprocDeleteOperationWait in the dataproc package instead.
func dataprocDeleteOperationWait(config *transport_tpg.Config, region, projectId, jobId, activity, userAgent string, timeout time.Duration) error {
	return dataproc.DataprocDeleteOperationWait(config, region, projectId, jobId, activity, userAgent, timeout)
}
