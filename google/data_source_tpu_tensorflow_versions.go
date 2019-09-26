package google

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTpuTensorflowVersions() *schema.Resource {
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

	versionsRaw, err := paginatedListRequest(project, url, config, flattenTpuTensorflowVersions)
	if err != nil {
		return fmt.Errorf("Error listing TPU Tensorflow versions: %s", err)
	}

	versions := make([]string, len(versionsRaw))
	for i, ver := range versionsRaw {
		versions[i] = ver.(string)
	}
	sort.Strings(versions)

	log.Printf("[DEBUG] Received Google TPU Tensorflow Versions: %q", versions)

	d.Set("versions", versions)
	d.Set("zone", zone)
	d.Set("project", project)
	d.SetId(time.Now().UTC().String())

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
