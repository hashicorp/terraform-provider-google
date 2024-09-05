// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-provider-google/google/fwprovider"
)

func TestFrameworkProvider_impl(t *testing.T) {
	var _ provider.ProviderWithMetaSchema = fwprovider.New()
}
