// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleIapClient() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceIapClient().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "brand", "client_id")

	return &schema.Resource{
		Read:   dataSourceGoogleIapClientRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleIapClientRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "{{brand}}/identityAwareProxyClients/{{client_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceIapClientRead(d, meta)
}
