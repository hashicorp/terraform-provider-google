package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleCloudFunctionsFunction() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudFunctionsFunction().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctionsFunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctionsFunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	d.SetId(cloudFuncId.terraformId())

	// terrible hack, remove when these fields are removed
	// We're temporarily reading these fields only when they are set
	// so we need them to be set with bad values entering read
	// and then unset if those bad values are still there
	d.Set("trigger_topic", "invalid")
	d.Set("trigger_bucket", "invalid")

	err = resourceCloudFunctionsRead(d, meta)
	if err != nil {
		return err
	}

	// terrible hack, remove when these fields are removed. see above
	if v := d.Get("trigger_topic").(string); v == "invalid" {
		d.Set("trigger_topic", "")
	}
	if v := d.Get("trigger_bucket").(string); v == "invalid" {
		d.Set("trigger_bucket", "")
	}

	return nil
}
