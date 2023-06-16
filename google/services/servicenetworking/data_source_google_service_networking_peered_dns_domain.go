// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleServiceNetworkingPeeredDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: resourceGoogleServiceNetworkingPeeredDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
