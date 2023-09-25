// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeRegionInstanceTemplate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeRegionInstanceTemplate().Schema)

	dsSchema["filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["most_recent"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "filter", "most_recent", "region", "project")

	dsSchema["name"].ExactlyOneOf = []string{"name", "filter"}
	dsSchema["filter"].ExactlyOneOf = []string{"name", "filter"}

	return &schema.Resource{
		Read:   datasourceComputeRegionInstanceTemplateRead,
		Schema: dsSchema,
	}
}

func datasourceComputeRegionInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		return retrieveInstances(d, meta, project, region, v.(string))
	}
	if v, ok := d.GetOk("filter"); ok {
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return err
		}

		params := map[string]string{
			"filter": v.(string),
		}

		url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceTemplates")
		if err != nil {
			return err
		}

		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		templates, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("error retrieving list of region instance templates: %s", err)
		}

		instanceTemplates := templates["items"]

		instanceTemplatesList, err := json.Marshal(instanceTemplates)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var items []*compute.InstanceTemplate

		if err := json.Unmarshal(instanceTemplatesList, &items); err != nil {
			fmt.Println(err)
			return err
		}

		mostRecent := d.Get("most_recent").(bool)
		if mostRecent {
			sort.Sort(ByCreationTimestamp(items))
		}

		count := len(items)
		if count == 1 || count > 1 && mostRecent {
			return retrieveInstances(d, meta, project, region, items[0].Name)
		}

		return fmt.Errorf("your filter has returned %d region instance template(s). Please refine your filter or set most_recent to return exactly one region instance template", len(items))
	}

	return fmt.Errorf("one of name or filters must be set")
}

func retrieveInstances(d *schema.ResourceData, meta interface{}, project, region, name string) error {
	d.SetId("projects/" + project + "/regions/" + region + "/instanceTemplates/" + name)

	if err := resourceComputeRegionInstanceTemplateRead(d, meta); err != nil {
		return err
	}
	return tpgresource.SetDataSourceLabels(d)
}
