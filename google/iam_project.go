// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var IamProjectSchema = resourcemanager.IamProjectSchema

func NewProjectIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return resourcemanager.NewProjectIamUpdater(d, config)
}

func ProjectIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	return resourcemanager.ProjectIdParseFunc(d, nil)
}

func compareProjectName(_, old, new string, _ *schema.ResourceData) bool {
	// We can either get "projects/project-id" or "project-id", so strip any prefixes
	return resourcemanager.CompareProjectName("", old, new, nil)
}
