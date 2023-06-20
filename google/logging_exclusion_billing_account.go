// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/services/logging"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var BillingAccountLoggingExclusionSchema = logging.BillingAccountLoggingExclusionSchema

func NewBillingAccountLoggingExclusionUpdater(d *schema.ResourceData, config *transport_tpg.Config) (logging.ResourceLoggingExclusionUpdater, error) {
	return logging.NewBillingAccountLoggingExclusionUpdater(d, config)
}

func BillingAccountLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	return logging.BillingAccountLoggingExclusionIdParseFunc(d, nil)
}
