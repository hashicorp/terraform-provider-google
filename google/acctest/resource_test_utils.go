// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// General test utils

var _ plancheck.PlanCheck = expectNoDelete{}

type expectNoDelete struct{}

func (e expectNoDelete) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result error
	for _, rc := range req.Plan.ResourceChanges {
		if slices.Contains(rc.Change.Actions, tfjson.ActionDelete) {
			result = errors.Join(result, fmt.Errorf("expected no deletion of resources, but %s has planned deletion", rc.Address))
		}
	}
	resp.Error = result
}

func ExpectNoDelete() plancheck.PlanCheck {
	return expectNoDelete{}
}

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

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Assume we never want to sleep when we're in replaying mode.
		if IsVcrEnabled() && os.Getenv("VCR_MODE") == "REPLAYING" {
			return nil
		}
		time.Sleep(time.Duration(t) * time.Second)
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
