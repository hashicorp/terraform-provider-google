package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleProject() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceGoogleProject().Schema)

	addOptionalFieldsToSchema(dsSchema, "project_id")

	dsSchema["project_id"].ValidateFunc = validateDSProjectID()
	return &schema.Resource{
		Read:   datasourceGoogleProjectRead,
		Schema: dsSchema,
	}
}

func datasourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if v, ok := d.GetOk("project_id"); ok {
		project := v.(string)
		d.SetId(fmt.Sprintf("projects/%s", project))
	} else {
		project, err := getProject(d, config)
		if err != nil {
			return fmt.Errorf("no project value set. `project_id` must be set at the resource level, or a default `project` value must be specified on the provider")
		}
		d.SetId(fmt.Sprintf("projects/%s", project))
	}

	return resourceGoogleProjectRead(d, meta)
}
