package google

import (
	"fmt"

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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	email, err := GetCurrentUserEmail(config, userAgent)
	if err != nil {
		return err
	}
	d.SetId(email)
	if err := d.Set("email", email); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	return nil
}
