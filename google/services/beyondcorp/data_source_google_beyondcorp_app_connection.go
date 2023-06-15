// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBeyondcorpAppConnection() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBeyondcorpAppConnection().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpAppConnectionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpAppConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("name").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/appConnections/%s", project, region, name))

	return resourceBeyondcorpAppConnectionRead(d, meta)
}
