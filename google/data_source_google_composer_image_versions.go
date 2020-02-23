package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComposerImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComposerImageVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_python_versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComposerImageVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComposerBasePath}}projects/{{project}}/locations/{{region}}/imageVersions")
	if err != nil {
		return err
	}

	versions, err := paginatedListRequest(project, url, config, flattenGoogleComposerImageVersions)
	if err != nil {
		return fmt.Errorf("Error listing Composer image versions: %s", err)
	}

	log.Printf("[DEBUG] Received Composer Image Versions: %q", versions)

	d.Set("image_versions", versions)
	d.Set("region", region)
	d.Set("project", project)
	d.SetId(time.Now().UTC().String())

	return nil
}

func flattenGoogleComposerImageVersions(resp map[string]interface{}) []interface{} {
	verObjList := resp["imageVersions"].([]interface{})
	versions := make([]interface{}, len(verObjList))
	for i, v := range verObjList {
		verObj := v.(map[string]interface{})
		versions[i] = map[string]interface{}{
			"image_version_id":          verObj["imageVersionId"],
			"supported_python_versions": verObj["supportedPythonVersions"],
		}
	}
	return versions
}
