package google

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeInstanceTemplate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeInstanceTemplate().Schema)

	dsSchema["filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["most_recent"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name", "filter", "most_recent", "project")

	dsSchema["name"].ExactlyOneOf = []string{"name", "filter"}
	dsSchema["filter"].ExactlyOneOf = []string{"name", "filter"}

	return &schema.Resource{
		Read:   datasourceComputeInstanceTemplateRead,
		Schema: dsSchema,
	}
}

func datasourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		return retrieveInstance(d, meta, project, v.(string))
	}
	if v, ok := d.GetOk("filter"); ok {
		userAgent, err := generateUserAgentString(d, config.userAgent)
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

	return fmt.Errorf("one of name or filters must be set")
}

func retrieveInstance(d *schema.ResourceData, meta interface{}, project, name string) error {
	d.SetId("projects/" + project + "/global/instanceTemplates/" + name)

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
