// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnetworks_self_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	network, err := config.NewComputeClient(userAgent).Networks.Get(project, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Network Not Found : %s", name))
	}
	if err := d.Set("gateway_ipv4", network.GatewayIPv4); err != nil {
		return fmt.Errorf("Error setting gateway_ipv4: %s", err)
	}
	if err := d.Set("self_link", network.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("description", network.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("subnetworks_self_links", network.Subnetworks); err != nil {
		return fmt.Errorf("Error setting subnetworks_self_links: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/global/networks/%s", project, network.Name))
	return nil
}
