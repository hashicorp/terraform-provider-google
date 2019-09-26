package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	projects := make([]map[string]interface{}, 0)

	for {
		params["filter"] = d.Get("filter").(string)
		url := "https://cloudresourcemanager.googleapis.com/v1/projects"

		url, err := addQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := sendRequest(config, "GET", "", url, nil)
		if err != nil {
			return fmt.Errorf("Error retrieving projects: %s", err)
		}

		pageProjects := flattenDatasourceGoogleProjectsList(res["projects"])
		projects = append(projects, pageProjects...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("projects", projects); err != nil {
		return fmt.Errorf("Error retrieving projects: %s", err)
	}

	d.SetId(d.Get("filter").(string))

	return nil
}

func flattenDatasourceGoogleProjectsList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	projects := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		p := raw.(map[string]interface{})
		if pId, ok := p["projectId"]; ok {
			projects = append(projects, map[string]interface{}{
				"project_id": pId,
			})
		}
	}

	return projects
}
