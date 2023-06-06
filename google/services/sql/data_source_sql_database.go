// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSqlDatabase() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceSQLDatabase().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, d.Get("instance").(string), d.Get("name").(string)))
	err = resourceSQLDatabaseRead(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("deletion_policy", nil); err != nil {
		return fmt.Errorf("Error setting deletion_policy: %s", err)
	}
	return nil
}
