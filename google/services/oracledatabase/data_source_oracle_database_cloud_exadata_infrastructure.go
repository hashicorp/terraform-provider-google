// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseCloudExadataInfrastructure() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseCloudExadataInfrastructure().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location", "cloud_exadata_infrastructure_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseCloudExadataInfrastructureRead,
		Schema: dsSchema,
	}

}

func dataSourceOracleDatabaseCloudExadataInfrastructureRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures/{{cloud_exadata_infrastructure_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	err = resourceOracleDatabaseCloudExadataInfrastructureRead(d, meta)
	if err != nil {
		return err
	}
	d.SetId(id)

	return nil
}
