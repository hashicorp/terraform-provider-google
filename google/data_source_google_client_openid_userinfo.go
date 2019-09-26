package google

import (
	"fmt"
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

	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, err := sendRequest(config, "GET", "", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return fmt.Errorf("error retrieving userinfo for your provider credentials; have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("email", res["email"])

	return nil
}
