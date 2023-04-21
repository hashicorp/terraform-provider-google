package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsKeyRing() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceKMSKeyRing().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addRequiredFieldsToSchema(dsSchema, "location")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleKmsKeyRingRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	keyRingId := KmsKeyRingId{
		Name:     d.Get("name").(string),
		Location: d.Get("location").(string),
		Project:  project,
	}
	d.SetId(keyRingId.KeyRingId())

	return resourceKMSKeyRingRead(d, meta)
}
