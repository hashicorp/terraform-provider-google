// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceSqlDatabaseInstance() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceSqlDatabaseInstance().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {

	return resourceSqlDatabaseInstanceRead(d, meta)

}
