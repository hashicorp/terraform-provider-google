// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeSubnetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSubnetworksRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnetworks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_cidr_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_google_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeSubnetworksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for subnetwork: %s", err)
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching region for subnetwork: %s", err)
	}

	filter := d.Get("filter").(string)

	subnetworks := make([]map[string]interface{}, 0)

	subnetworkList, err := config.NewComputeClient(userAgent).Subnetworks.List(project, region).Filter(filter).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Subnetworks : %s %s", project, region))
	}

	for _, subnet := range subnetworkList.Items {
		subnetworks = append(subnetworks, map[string]interface{}{
			"description":              subnet.Description,
			"ip_cidr_range":            subnet.IpCidrRange,
			"name":                     subnet.Name,
			"network_self_link":        filepath.Base(subnet.Network),
			"network":                  subnet.Network,
			"private_ip_google_access": subnet.PrivateIpGoogleAccess,
			"self_link":                subnet.SelfLink,
		})
	}

	if err := d.Set("subnetworks", subnetworks); err != nil {
		return fmt.Errorf("Error retrieving subnetworks: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s/regions/%s/subnetworks",
		project,
		region,
	))

	return nil
}
