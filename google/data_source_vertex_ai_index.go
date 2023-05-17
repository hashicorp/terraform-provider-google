package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func dataSourceVertexAIIndex() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceVertexAIIndex().Schema)

	addRequiredFieldsToSchema(dsSchema, "name", "region")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceVertexAIIndexRead,
		Schema: dsSchema,
	}
}

func dataSourceVertexAIIndexRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/indexes/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceVertexAIIndexRead(d, meta)
}
