// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/services/logging"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var ProjectLoggingExclusionSchema = logging.ProjectLoggingExclusionSchema

func NewProjectLoggingExclusionUpdater(d *schema.ResourceData, config *transport_tpg.Config) (logging.ResourceLoggingExclusionUpdater, error) {
	return logging.NewProjectLoggingExclusionUpdater(d, config)
}

func ProjectLoggingExclusionIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	return logging.ProjectLoggingExclusionIdParseFunc(d, config)
}
