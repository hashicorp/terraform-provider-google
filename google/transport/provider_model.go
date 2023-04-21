package transport

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderBatching struct {
	SendAfter      types.String `tfsdk:"send_after"`
	EnableBatching types.Bool   `tfsdk:"enable_batching"`
}
