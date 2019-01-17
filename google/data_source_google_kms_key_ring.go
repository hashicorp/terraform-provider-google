package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleKmsKeyRing() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceKmsKeyRing().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addRequiredFieldsToSchema(dsSchema, "location")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleKmsKeyRingRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	keyRingId := kmsKeyRingId{
		Name:     d.Get("name").(string),
		Location: d.Get("location").(string),
		Project:  project,
	}
	d.SetId(keyRingId.terraformId())

	return resourceKmsKeyRingRead(d, meta)
}
