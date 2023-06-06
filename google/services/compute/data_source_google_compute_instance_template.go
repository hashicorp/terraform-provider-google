// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeInstanceTemplate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstanceTemplate().Schema)

	dsSchema["filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["self_link_unique"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["most_recent"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "filter", "most_recent", "project", "self_link_unique")

	mutuallyExclusive := []string{"name", "filter", "self_link_unique"}
	for _, n := range mutuallyExclusive {
		dsSchema[n].ExactlyOneOf = mutuallyExclusive
	}

	return &schema.Resource{
		Read:   datasourceComputeInstanceTemplateRead,
		Schema: dsSchema,
	}
}

func datasourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		return retrieveInstance(d, meta, project, v.(string))
	}
	if v, ok := d.GetOk("filter"); ok {
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return err
		}

		templates, err := config.NewComputeClient(userAgent).InstanceTemplates.List(project).Filter(v.(string)).Do()
		if err != nil {
			return fmt.Errorf("error retrieving list of instance templates: %s", err)
		}

		mostRecent := d.Get("most_recent").(bool)
		if mostRecent {
			sort.Sort(ByCreationTimestamp(templates.Items))
		}

		count := len(templates.Items)
		if count == 1 || count > 1 && mostRecent {
			return retrieveInstance(d, meta, project, templates.Items[0].Name)
		}

		return fmt.Errorf("your filter has returned %d instance template(s). Please refine your filter or set most_recent to return exactly one instance template", len(templates.Items))
	}
	if v, ok := d.GetOk("self_link_unique"); ok {
		return retrieveInstanceFromUniqueId(d, meta, project, v.(string))
	}

	return fmt.Errorf("one of name, filters or self_link_unique must be set")
}

func retrieveInstance(d *schema.ResourceData, meta interface{}, project, name string) error {
	d.SetId("projects/" + project + "/global/instanceTemplates/" + name)

	return resourceComputeInstanceTemplateRead(d, meta)
}

func retrieveInstanceFromUniqueId(d *schema.ResourceData, meta interface{}, project, self_link_unique string) error {
	normalId, _ := parseUniqueId(self_link_unique)
	d.SetId(normalId)
	d.Set("self_link_unique", self_link_unique)

	return resourceComputeInstanceTemplateRead(d, meta)
}

// ByCreationTimestamp implements sort.Interface for []*InstanceTemplate based on
// the CreationTimestamp field.
type ByCreationTimestamp []*compute.InstanceTemplate

func (a ByCreationTimestamp) Len() int      { return len(a) }
func (a ByCreationTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreationTimestamp) Less(i, j int) bool {
	return a[i].CreationTimestamp > a[j].CreationTimestamp
}
