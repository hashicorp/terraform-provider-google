// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/provider"
)

// Provider returns a *schema.Provider.

// This function stays in the google package temporarily to not break the terraform-google-conversion.
// TODO: remove it in a later PR
func Provider() *schema.Provider {
	return provider.Provider()
}

// Generated resources: <%= resource_count %>
// Generated IAM resources: <%= iam_resource_count %>
// Total generated resources: <%= resource_count + iam_resource_count %>

// This function stays in the google package temporarily to not break the tools missting tests detector and breaking changes detector.
// TODO: remove it in a later PR
func ResourceMap() map[string]*schema.Resource {
	return provider.ResourceMap()
}
