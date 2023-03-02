package google

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTpuTensorflowVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTpuTensorFlowVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTpuTensorFlowVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{TPUBasePath}}projects/{{project}}/locations/{{zone}}/tensorflowVersions")
	if err != nil {
		return err
	}

	versionsRaw, err := paginatedListRequest(project, url, userAgent, config, flattenTpuTensorflowVersions)
	if err != nil {
		return fmt.Errorf("Error listing TPU Tensorflow versions: %s", err)
	}

	versions := make([]string, len(versionsRaw))
	for i, ver := range versionsRaw {
		versions[i] = ver.(string)
	}
	sort.Strings(versions)

	log.Printf("[DEBUG] Received Google TPU Tensorflow Versions: %q", versions)

	if err := d.Set("versions", versions); err != nil {
		return fmt.Errorf("Error setting versions: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s", project, zone))

	return nil
}

func flattenTpuTensorflowVersions(resp map[string]interface{}) []interface{} {
	verObjList := resp["tensorflowVersions"].([]interface{})
	versions := make([]interface{}, len(verObjList))
	for i, v := range verObjList {
		verObj := v.(map[string]interface{})
		versions[i] = verObj["version"]
	}
	return versions
}
