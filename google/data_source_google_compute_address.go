package google

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleComputeAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeAddressRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	address, err := config.clientCompute.Addresses.Get(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore

			return fmt.Errorf("Address Not Found")
		}

		return fmt.Errorf("Error reading Address: %s", err)
	}

	d.Set("address", address.Address)
	d.Set("status", address.Status)
	d.Set("self_link", address.SelfLink)
	d.Set("project", project)
	d.Set("region", region)

	d.SetId(strconv.FormatUint(uint64(address.Id), 10))
	return nil
}
