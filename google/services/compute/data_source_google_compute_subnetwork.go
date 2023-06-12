// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeSubnetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSubnetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_cidr_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_google_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"secondary_ip_range": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_cidr_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeSubnetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, region, name, err := tpgresource.GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	subnetwork, err := config.NewComputeClient(userAgent).Subnetworks.Get(project, region, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Subnetwork Not Found : %s", name))
	}

	if err := d.Set("ip_cidr_range", subnetwork.IpCidrRange); err != nil {
		return fmt.Errorf("Error setting ip_cidr_range: %s", err)
	}
	if err := d.Set("private_ip_google_access", subnetwork.PrivateIpGoogleAccess); err != nil {
		return fmt.Errorf("Error setting private_ip_google_access: %s", err)
	}
	if err := d.Set("self_link", subnetwork.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("description", subnetwork.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("gateway_address", subnetwork.GatewayAddress); err != nil {
		return fmt.Errorf("Error setting gateway_address: %s", err)
	}
	if err := d.Set("network", subnetwork.Network); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("secondary_ip_range", flattenSecondaryRanges(subnetwork.SecondaryIpRanges)); err != nil {
		return fmt.Errorf("Error setting secondary_ip_range: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", project, region, name))
	return nil
}

func flattenSecondaryRanges(secondaryRanges []*compute.SubnetworkSecondaryRange) []map[string]interface{} {
	secondaryRangesSchema := make([]map[string]interface{}, 0, len(secondaryRanges))
	for _, secondaryRange := range secondaryRanges {
		data := map[string]interface{}{
			"range_name":    secondaryRange.RangeName,
			"ip_cidr_range": secondaryRange.IpCidrRange,
		}

		secondaryRangesSchema = append(secondaryRangesSchema, data)
	}
	return secondaryRangesSchema
}
