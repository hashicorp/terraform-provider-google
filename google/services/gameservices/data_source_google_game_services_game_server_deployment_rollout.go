// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gameservices

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGameServicesGameServerDeploymentRollout() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGameServicesGameServerDeploymentRollout().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "deployment_id")

	return &schema.Resource{
		Read:   dataSourceGameServicesGameServerDeploymentRolloutRead,
		Schema: dsSchema,
	}
}

func dataSourceGameServicesGameServerDeploymentRolloutRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/global/gameServerDeployments/{{deployment_id}}/rollout")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	return resourceGameServicesGameServerDeploymentRolloutRead(d, meta)
}
