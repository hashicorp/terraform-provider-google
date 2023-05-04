package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudFunctionsFunction() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCloudFunctionsFunction().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctionsFunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctionsFunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	cloudFuncId := &cloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	d.SetId(cloudFuncId.cloudFunctionId())

	err = resourceCloudFunctionsRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
