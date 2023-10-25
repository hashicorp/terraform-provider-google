// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigqueryDataset() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBigQueryDataset().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "dataset_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleBigqueryDatasetRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBigqueryDatasetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	dataset_id := d.Get("dataset_id").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

	id := fmt.Sprintf("projects/%s/datasets/%s", project, dataset_id)
	d.SetId(id)
	err = resourceBigQueryDatasetRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
