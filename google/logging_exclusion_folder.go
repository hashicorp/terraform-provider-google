// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/services/logging"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var FolderLoggingExclusionSchema = logging.FolderLoggingExclusionSchema

func NewFolderLoggingExclusionUpdater(d *schema.ResourceData, config *transport_tpg.Config) (logging.ResourceLoggingExclusionUpdater, error) {
	return logging.NewFolderLoggingExclusionUpdater(d, config)
}

func FolderLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	return logging.FolderLoggingExclusionIdParseFunc(d, nil)
}
