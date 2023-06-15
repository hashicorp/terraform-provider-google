// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetBatchingConfig returns the batching config object given the
// provider configuration set for batching
func GetBatchingConfig(ctx context.Context, data types.List, diags *diag.Diagnostics) *batchingConfig {
	bc := &batchingConfig{
		SendAfter:      time.Second * DefaultBatchSendIntervalSec,
		EnableBatching: true,
	}

	if data.IsNull() {
		return bc
	}

	var pbConfigs []ProviderBatching
	d := data.ElementsAs(ctx, &pbConfigs, true)
	diags.Append(d...)
	if diags.HasError() {
		return bc
	}

	sendAfter, err := time.ParseDuration(pbConfigs[0].SendAfter.ValueString())
	if err != nil {
		diags.AddError("error parsing send after time duration", err.Error())
		return bc
	}

	bc.SendAfter = sendAfter

	if !pbConfigs[0].EnableBatching.IsNull() {
		bc.EnableBatching = pbConfigs[0].EnableBatching.ValueBool()
	}

	return bc
}
