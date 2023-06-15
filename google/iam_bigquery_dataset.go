// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/bigquery"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var IamBigqueryDatasetSchema = bigquery.IamBigqueryDatasetSchema

func NewBigqueryDatasetIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return bigquery.NewBigqueryDatasetIamUpdater(d, config)
}

func BigqueryDatasetIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	return bigquery.BigqueryDatasetIdParseFunc(d, config)
}
