// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"testing"
)

func TestInstanceGroupManager_parseUniqueId(t *testing.T) {
	expectations := map[string][]string{
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       {"projects/imre-test/global/instanceTemplates/example-template-custom", "123"},
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": {"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom", "123"},
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                                    {"projects/imre-test/global/instanceTemplates/example-template-custom", ""},
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom":              {"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom", ""},
		"example-template-custom?uniqueId=123":                                                                                   {"example-template-custom", "123"},

		// this test demonstrates that uniqueIds can't override eachother
		"projects/imre-test/global/instanceTemplates/example?uniqueId=123?uniqueId=456": {"projects/imre-test/global/instanceTemplates/example", "123?uniqueId=456"},
	}

	for k, v := range expectations {
		aName, aUniqueId := parseUniqueId(k)
		if v[0] != aName {
			t.Errorf("parseUniqueId failed; name of %v should be %v, not %v", k, v[0], aName)
		}
		if v[1] != aUniqueId {
			t.Errorf("parseUniqueId failed; uniqueId of %v should be %v, not %v", k, v[1], aUniqueId)
		}
	}
}

func TestInstanceGroupManager_compareInstanceTemplate(t *testing.T) {
	shouldAllMatch := []string{
		// uniqueId not present
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom",
		"projects/imre-test/global/instanceTemplates/example-template-custom",
		// uniqueId present
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123",
	}
	shouldNotMatch := map[string]string{
		// mismatching name
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom": "projects/imre-test/global/instanceTemplates/example-template-custom2",
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                       "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom2",
		// matching name, but mismatching uniqueId
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": "projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=1234",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=1234",
	}
	for _, v1 := range shouldAllMatch {
		for _, v2 := range shouldAllMatch {
			if !compareSelfLinkRelativePathsIgnoreParams("", v1, v2, nil) {
				t.Fatalf("compareSelfLinkRelativePathsIgnoreParams did not match (and should have) %v and %v", v1, v2)
			}
		}
	}

	for v1, v2 := range shouldNotMatch {
		if compareSelfLinkRelativePathsIgnoreParams("", v1, v2, nil) {
			t.Fatalf("compareSelfLinkRelativePathsIgnoreParams did match (and shouldn't) %v and %v", v1, v2)
		}
	}
}

func TestInstanceGroupManager_convertUniqueId(t *testing.T) {
	matches := map[string]string{
		// uniqueId not present (should return the same)
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom": "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom",
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                       "projects/imre-test/global/instanceTemplates/example-template-custom",
		// uniqueId present (should return the last component replaced)
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/123",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       "projects/imre-test/global/instanceTemplates/123",
		"tf-test-igm-8amncgtq22?uniqueId=8361222501423044003":                                                                    "8361222501423044003",
	}
	for input, expected := range matches {
		actual := ConvertToUniqueIdWhenPresent(input)
		if actual != expected {
			t.Fatalf("invalid return value by ConvertToUniqueIdWhenPresent for input %v; expected: %v, actual: %v", input, expected, actual)
		}
	}
}
