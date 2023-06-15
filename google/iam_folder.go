// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var IamFolderSchema = resourcemanager.IamFolderSchema

func NewFolderIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return resourcemanager.NewFolderIamUpdater(d, config)
}

func FolderIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	return resourcemanager.FolderIdParseFunc(d, nil)
}

func canonicalFolderId(folder string) string {
	return resourcemanager.CanonicalFolderId(folder)
}

// Retrieve the existing IAM Policy for a folder
func getFolderIamPolicyByFolderName(folderName, userAgent string, config *transport_tpg.Config) (*cloudresourcemanager.Policy, error) {
	return resourcemanager.GetFolderIamPolicyByFolderName(folderName, userAgent, config)
}
