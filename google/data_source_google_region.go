package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func dataSourceGoogleRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleRegionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleRegionRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(time.Now().UTC().String())
	return nil
}