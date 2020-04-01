package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	d.Set("email", email)
	return nil
}
