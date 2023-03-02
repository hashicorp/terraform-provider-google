package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleContainerAttachedInstallManifest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerAttachedInstallManifestRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manifest": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleContainerAttachedInstallManifestRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	clusterId := d.Get("cluster_id").(string)
	platformVersion := d.Get("platform_version").(string)

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

	url, err := replaceVars(d, config, "{{ContainerAttachedBasePath}}projects/{{project}}/locations/{{location}}:generateAttachedClusterInstallManifest")
	if err != nil {
		return err
	}
	params := map[string]string{
		"attached_cluster_id": clusterId,
		"platform_version":    platformVersion,
	}
	url, err = addQueryParams(url, params)
	if err != nil {
		return err
	}
	res, err := SendRequest(config, "GET", project, url, userAgent, nil)
	if err != nil {
		return err
	}

	if err := d.Set("manifest", res["manifest"]); err != nil {
		return fmt.Errorf("Error setting manifest: %s", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
