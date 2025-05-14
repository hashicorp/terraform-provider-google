// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBeyondcorpSecurityGateway() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBeyondcorpSecurityGateway().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "security_gateway_id")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleBeyondcorpSecurityGatewayRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleBeyondcorpSecurityGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("security_gateway_id").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/locations/global/securityGateways/%s", project, name)
	d.SetId(id)
	// Remove after deprecated location variable is removed.
	d.Set("location", "global")

	err = resourceBeyondcorpSecurityGatewayRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
