package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/zones/%s", config.Project, config.Region, config.Zone))
	if err := d.Set("project", config.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", config.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("zone", config.Zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	token, err := config.tokenSource.Token()
	if err != nil {
		return err
	}
	if err := d.Set("access_token", token.AccessToken); err != nil {
		return fmt.Errorf("Error setting access_token: %s", err)
	}

	return nil
}
