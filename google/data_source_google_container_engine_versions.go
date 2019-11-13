package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleContainerEngineVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerEngineVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use location instead",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use location instead",
			},
			"default_cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_node_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_master_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"valid_node_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleContainerEngineVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	location = fmt.Sprintf("projects/%s/locations/%s", project, location)
	resp, err := config.clientContainerBeta.Projects.Locations.GetServerConfig(location).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving available container cluster versions: %s", err.Error())
	}

	validMasterVersions := make([]string, 0)
	for _, v := range resp.ValidMasterVersions {
		if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
			validMasterVersions = append(validMasterVersions, v)
		}
	}

	validNodeVersions := make([]string, 0)
	for _, v := range resp.ValidNodeVersions {
		if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
			validNodeVersions = append(validNodeVersions, v)
		}
	}

	d.Set("valid_master_versions", validMasterVersions)
	if len(validMasterVersions) > 0 {
		d.Set("latest_master_version", validMasterVersions[0])
	}

	d.Set("valid_node_versions", validNodeVersions)
	if len(validNodeVersions) > 0 {
		d.Set("latest_node_version", validNodeVersions[0])
	}

	d.Set("default_cluster_version", resp.DefaultClusterVersion)

	d.SetId(time.Now().UTC().String())
	return nil
}
