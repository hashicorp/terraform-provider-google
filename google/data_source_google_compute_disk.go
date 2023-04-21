package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeDisk() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeDisk().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "zone")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeDiskRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := ReplaceVars(d, config, "projects/{{project}}/zones/{{zone}}/disks/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceComputeDiskRead(d, meta)
}
