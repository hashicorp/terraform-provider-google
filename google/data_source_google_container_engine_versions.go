package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleContainerEngineVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerEngineVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"zone"},
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
		return fmt.Errorf("Cannot determine location: set zone or region in this data source or at provider-level")
	}

	location = fmt.Sprintf("projects/%s/locations/%s", project, location)
	resp, err := config.clientContainerBeta.Projects.Locations.GetServerConfig(location).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving available container cluster versions: %s", err.Error())
	}

	d.Set("valid_master_versions", resp.ValidMasterVersions)
	d.Set("default_cluster_version", resp.DefaultClusterVersion)
	d.Set("valid_node_versions", resp.ValidNodeVersions)
	if len(resp.ValidMasterVersions) > 0 {
		d.Set("latest_master_version", resp.ValidMasterVersions[0])
	}
	if len(resp.ValidNodeVersions) > 0 {
		d.Set("latest_node_version", resp.ValidNodeVersions[0])
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
