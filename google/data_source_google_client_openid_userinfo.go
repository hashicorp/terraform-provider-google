package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleClientOpenIDUserinfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleClientOpenIDUserinfoRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleClientOpenIDUserinfoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	email, err := GetCurrentUserEmail(config)
	if err != nil {
		return err
	}
	d.SetId(time.Now().UTC().String())
	if err := d.Set("email", email); err != nil {
		return fmt.Errorf("Error reading email: %s", err)
	}
	return nil
}
