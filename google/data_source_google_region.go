package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleRegionRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGoogleRegionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	d.SetId(config.Region)
	return nil
}