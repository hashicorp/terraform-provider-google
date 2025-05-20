// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleProjects() *schema.Resource {
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
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `A set of key/value label pairs assigned on a project.`,
						},
						"parent": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `An optional reference to a parent Resource.`,
						},
						"number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The numeric identifier of the project.`,
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The numeric identifier of the project.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The optional user-assigned display name of the Project.`,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleProjectsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	projects := make([]map[string]interface{}, 0)

	for {
		params["filter"] = d.Get("filter").(string)
		url := "https://cloudresourcemanager.googleapis.com/v1/projects"

		url, err := transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
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

		var mId, mNumber, mLabels, mLifecycleState, mName, mCreateTime, mParent interface{}
		if pId, ok := p["projectId"]; ok {
			mId = pId
		}
		if pNumber, ok := p["projectNumber"]; ok {
			mNumber = pNumber
		}
		if pName, ok := p["name"]; ok {
			mName = pName
		}
		if pLabels, ok := p["labels"]; ok {
			mLabels = pLabels
		}
		if pLifecycleState, ok := p["lifecycleState"]; ok {
			mLifecycleState = pLifecycleState
		}
		if pCreateTime, ok := p["createTime"]; ok {
			mCreateTime = pCreateTime
		}
		if pParent, ok := p["parent"]; ok {
			mParent = pParent
		}
		projects = append(projects, map[string]interface{}{
			"project_id":      mId,
			"number":          mNumber,
			"name":            mName,
			"labels":          mLabels,
			"lifecycle_state": mLifecycleState,
			"create_time":     mCreateTime,
			"parent":          mParent,
		})
	}

	return projects
}
