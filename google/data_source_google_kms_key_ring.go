package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleKmsKeyRing() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceKmsKeyRing().Schema)

	addRequiredFieldsToSchema(dsSchema, "name", "location")

	return &schema.Resource{
		Read:   datasourceGoogleKmsKeyRingRead,
		Schema: dsSchema,
	}
}

func datasourceGoogleKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	keyRing := &kmsKeyRingId{
		Project:  project,
		Location: d.Get("location").(string),
		Name:     d.Get("name").(string),
	}

	d.SetId(keyRing.keyRingId())

	return resourceKmsKeyRingRead(d, meta)
}
