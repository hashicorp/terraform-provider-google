package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleProjects() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleProjectsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleProjectsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	params := make(map[string]string)

	params["filter"] = d.Get("filter").(string)
	url := "https://cloudresourcemanager.googleapis.com/v1/projects"

	url, err := addQueryParams(url, params)
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving projects: %s", err)
	}

	if err := d.Set("projects", flattenDatasourceGoogleProjectsProjects(res["projects"], d)); err != nil {
		return fmt.Errorf("Error retrieving projects: %s", err)
	}

	d.SetId(d.Get("filter").(string))

	return nil
}

func flattenDatasourceGoogleProjectsProjects(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}

	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"project_id": original["projectId"],
		})
	}

	return transformed
}
