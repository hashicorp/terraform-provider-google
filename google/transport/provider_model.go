// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderBatching struct {
	SendAfter      types.String `tfsdk:"send_after"`
	EnableBatching types.Bool   `tfsdk:"enable_batching"`
}
