package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleClientConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClientConfigRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceClientConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.SetId(time.Now().UTC().String())
	d.Set("project", config.Project)
	d.Set("region", config.Region)
	d.Set("zone", config.Zone)

	token, err := config.tokenSource.Token()
	if err != nil {
		return err
	}
	d.Set("access_token", token.AccessToken)

	return nil
}
