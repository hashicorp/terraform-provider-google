// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

const regexGCEName = "^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$"

func DataSourceComputeNetworkPeering() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeNetworkPeering().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name", "network")

	dsSchema["name"].ValidateFunc = verify.ValidateRegexp(regexGCEName)
	dsSchema["network"].ValidateFunc = verify.ValidateRegexp(peerNetworkLinkRegex)
	return &schema.Resource{
		Read:   dataSourceComputeNetworkPeeringRead,
		Schema: dsSchema,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(4 * time.Minute),
		},
	}
}

func dataSourceComputeNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	networkFieldValue, err := tpgresource.ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", networkFieldValue.Name, d.Get("name").(string)))

	return resourceComputeNetworkPeeringRead(d, meta)
}
