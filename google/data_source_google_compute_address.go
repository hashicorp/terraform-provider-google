// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func parseComputeAddressId(id string, config *transport_tpg.Config) (*compute.ComputeAddressId, error) {
	return compute.ParseComputeAddressId(id, config)
}
