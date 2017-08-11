package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeInstanceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instances": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"named_port": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)

	d.SetId(fmt.Sprintf("%s/%s", zone, name))

	return resourceComputeInstanceGroupRead(d, meta)
}
