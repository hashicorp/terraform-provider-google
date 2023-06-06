// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudbuild

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBuildTrigger() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCloudBuildTrigger().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "trigger_id", "location")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudBuildTriggerRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleCloudBuildTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/triggers/{{trigger_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	id = strings.ReplaceAll(id, "/locations/global/", "/")

	d.SetId(id)
	return resourceCloudBuildTriggerRead(d, meta)
}
