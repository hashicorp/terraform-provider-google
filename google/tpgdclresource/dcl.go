// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgdclresource

import (
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

var (
	// CreateDirective restricts Apply to creating resources for Create
	CreateDirective = []dcl.ApplyOption{
		dcl.WithLifecycleParam(dcl.BlockAcquire),
		dcl.WithLifecycleParam(dcl.BlockDestruction),
		dcl.WithLifecycleParam(dcl.BlockModification),
	}

	// UpdateDirective restricts Apply to modifying resources for Update
	UpdateDirective = []dcl.ApplyOption{
		dcl.WithLifecycleParam(dcl.BlockCreation),
		dcl.WithLifecycleParam(dcl.BlockDestruction),
	}
)
