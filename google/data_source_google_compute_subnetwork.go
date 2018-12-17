package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeSubnetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSubnetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	subnetwork, err := config.clientCompute.Subnetworks.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Subnetwork Not Found : %s", name))
	}

	d.Set("ip_cidr_range", subnetwork.IpCidrRange)
	d.Set("private_ip_google_access", subnetwork.PrivateIpGoogleAccess)
	d.Set("self_link", subnetwork.SelfLink)
	d.Set("description", subnetwork.Description)
	d.Set("gateway_address", subnetwork.GatewayAddress)
	d.Set("network", subnetwork.Network)
	d.Set("project", project)
	d.Set("region", region)
	// Flattening code defined in resource_compute_subnetwork.go
	d.Set("secondary_ip_range", flattenSecondaryRanges(subnetwork.SecondaryIpRanges))

	//Subnet id creation is defined in resource_compute_subnetwork.go
	subnetwork.Region = region
	d.SetId(createSubnetID(subnetwork))
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

func createSubnetID(s *compute.Subnetwork) string {
	return fmt.Sprintf("%s/%s", s.Region, s.Name)
}
