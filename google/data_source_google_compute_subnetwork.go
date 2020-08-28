package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeSubnetwork() *schema.Resource {
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
	config := meta.(*Config)

	project, region, name, err := GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	subnetwork, err := config.clientCompute.Subnetworks.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Subnetwork Not Found : %s", name))
	}

	if err := d.Set("ip_cidr_range", subnetwork.IpCidrRange); err != nil {
		return fmt.Errorf("Error reading ip_cidr_range: %s", err)
	}
	if err := d.Set("private_ip_google_access", subnetwork.PrivateIpGoogleAccess); err != nil {
		return fmt.Errorf("Error reading private_ip_google_access: %s", err)
	}
	if err := d.Set("self_link", subnetwork.SelfLink); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("description", subnetwork.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("gateway_address", subnetwork.GatewayAddress); err != nil {
		return fmt.Errorf("Error reading gateway_address: %s", err)
	}
	if err := d.Set("network", subnetwork.Network); err != nil {
		return fmt.Errorf("Error reading network: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error reading region: %s", err)
	}
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("secondary_ip_range", flattenSecondaryRanges(subnetwork.SecondaryIpRanges)); err != nil {
		return fmt.Errorf("Error reading secondary_ip_range: %s", err)
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
