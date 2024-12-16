// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"strconv"

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

			"network_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Deprecated in favor of network_id
			"numeric_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "`numeric_id` is deprecated and will be removed in a future major release. Use `network_id` instead.",
			},

			"gateway_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"internal_ipv6_range": {
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

	id := fmt.Sprintf("projects/%s/global/networks/%s", project, name)

	network, err := config.NewComputeClient(userAgent).Networks.Get(project, name).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Network Not Found : %s", name), id)
	}
	if err := d.Set("gateway_ipv4", network.GatewayIPv4); err != nil {
		return fmt.Errorf("Error setting gateway_ipv4: %s", err)
	}
	if err := d.Set("internal_ipv6_range", network.InternalIpv6Range); err != nil {
		return fmt.Errorf("Error setting internal_ipv6_range: %s", err)
	}
	if err := d.Set("self_link", network.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("description", network.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("network_id", network.Id); err != nil {
		return fmt.Errorf("Error setting network_id: %s", err)
	}
	if err := d.Set("numeric_id", strconv.Itoa(int(network.Id))); err != nil {
		return fmt.Errorf("Error setting numeric_id: %s", err)
	}
	if err := d.Set("subnetworks_self_links", network.Subnetworks); err != nil {
		return fmt.Errorf("Error setting subnetworks_self_links: %s", err)
	}
	d.SetId(id)
	return nil
}
