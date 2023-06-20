// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/services/dataproc"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var IamDataprocJobSchema = dataproc.IamDataprocJobSchema

func NewDataprocJobUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return dataproc.NewDataprocJobUpdater(d, config)
}

func DataprocJobIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	return dataproc.DataprocJobIdParseFunc(d, config)
}
