// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NOTE(craigatgoogle): An out of band aspect of this API is that it uses a unique formatting of network
// different from the standard self_link URI. It requires a call to the resource manager to get the project
// number for the current project.
func retrieveServiceNetworkingNetworkName(d *schema.ResourceData, config *transport_tpg.Config, network, userAgent string) (string, error) {
	return servicenetworking.RetrieveServiceNetworkingNetworkName(d, config, network, userAgent)

}
