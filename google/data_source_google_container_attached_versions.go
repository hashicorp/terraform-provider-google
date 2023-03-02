package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleContainerAttachedVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerAttachedVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"valid_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleContainerAttachedVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	url, err := replaceVars(d, config, "{{ContainerAttachedBasePath}}projects/{{project}}/locations/{{location}}/attachedServerConfig")
	if err != nil {
		return err
	}
	res, err := SendRequest(config, "GET", project, url, userAgent, nil)
	if err != nil {
		return err
	}
	var validVersions []string
	for _, v := range res["validVersions"].([]interface{}) {
		vm := v.(map[string]interface{})
		validVersions = append(validVersions, vm["version"].(string))
	}
	if err := d.Set("valid_versions", validVersions); err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
