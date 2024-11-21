// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-provider-google/google/fwprovider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFrameworkProvider_impl(t *testing.T) {
	primary := &schema.Provider{}
	var _ provider.ProviderWithMetaSchema = fwprovider.New(primary)
}
