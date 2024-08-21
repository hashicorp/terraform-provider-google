// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleContainerCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceContainerCluster().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   datasourceContainerClusterRead,
		Schema: dsSchema,
	}
}

func datasourceContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	clusterName := d.Get("name").(string)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	id := containerClusterFullName(project, location, clusterName)

	d.SetId(id)

	if err := resourceContainerClusterRead(d, meta); err != nil {
		return err
	}

	// Sets the "resource_labels" field and "terraform_labels" with the value of the field "effective_labels".
	effectiveLabels := d.Get("effective_labels")
	if err := d.Set("resource_labels", effectiveLabels); err != nil {
		return fmt.Errorf("Error setting labels in data source: %s", err)
	}
	if err := d.Set("terraform_labels", effectiveLabels); err != nil {
		return fmt.Errorf("Error setting terraform_labels in data source: %s", err)
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
