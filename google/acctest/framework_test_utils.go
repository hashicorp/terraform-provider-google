// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func GetFwTestProvider(t *testing.T) *frameworkTestProvider {
	configsLock.RLock()
	fwProvider, ok := fwProviders[t.Name()]
	configsLock.RUnlock()
	if ok {
		return fwProvider
	}

	var diags diag.Diagnostics
	p := NewFrameworkTestProvider(t.Name())
	configureApiClient(context.Background(), &p.FrameworkProvider, &diags)
	if diags.HasError() {
		log.Fatalf("%d errors when configuring test provider client: first is %s", diags.ErrorsCount(), diags.Errors()[0].Detail())
	}

	return p
}

// General test utils

// TestExtractResourceAttr navigates a test's state to find the specified resource (or data source) attribute and makes the value
// accessible via the attributeValue string pointer.
func TestExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName] // To find a datasource, include `data.` at the start of the resourceName value

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

// TestCheckAttributeValuesEqual compares two string pointers, which have been used to retrieve attribute values from the test's state.
func TestCheckAttributeValuesEqual(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) != testStringValue(j) {
			return fmt.Errorf("attribute values are different, got %s and %s", testStringValue(i), testStringValue(j))
		}

		return nil
	}
}

// testStringValue returns string values from string pointers, handling nil pointers.
func testStringValue(sPtr *string) string {
	if sPtr == nil {
		return ""
	}

	return *sPtr
}
