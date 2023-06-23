// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
)

// General test utils

// testExtractResourceAttr navigates a test's state to find the specified resource (or data source) attribute and makes the value
// accessible via the attributeValue string pointer.
func testExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) resource.TestCheckFunc {
	return acctest.TestExtractResourceAttr(resourceName, attributeName, attributeValue)
}

// testCheckAttributeValuesEqual compares two string pointers, which have been used to retrieve attribute values from the test's state.
func testCheckAttributeValuesEqual(i *string, j *string) resource.TestCheckFunc {
	return acctest.TestCheckAttributeValuesEqual(i, j)
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func replaceVarsForFrameworkTest(prov *fwtransport.FrameworkProviderConfig, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	return acctest.ReplaceVarsForFrameworkTest(prov, rs, linkTmpl)
}
