// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions

import (
	"testing"
)

func TestCloudFunctionsFunction_nameValidator(t *testing.T) {
	validNames := []string{
		"a",
		"aA",
		"a0",
		"has-hyphen",
		"has_underscore",
		"hasUpperCase",
		"allChars_-A0",
		"StartsUpperCase",
		"endsUpperCasE",
	}
	for _, tc := range validNames {
		wrns, errs := validateResourceCloudFunctionsFunctionName(tc, "function.name")
		if len(wrns) > 0 {
			t.Errorf("Expected no validation warnings for test case %q, got: %+v", tc, wrns)
		}
		if len(errs) > 0 {
			t.Errorf("Expected no validation errors for test name %q, got: %+v", tc, errs)
		}
	}

	invalidNames := []string{
		"0startsWithNumber",
		"endsWith_",
		"endsWith-",
		"bad*Character",
		"aCloudFunctionsFunctionNameThatIsSeventyFiveCharactersLongWhichIsMoreThan63",
	}
	for _, tc := range invalidNames {
		_, errs := validateResourceCloudFunctionsFunctionName(tc, "function.name")
		if len(errs) == 0 {
			t.Errorf("Expected errors for invalid test name %q, got none", tc)
		}
	}
}

func TestValidLabelKeys(t *testing.T) {
	testCases := []struct {
		labelKey string
		valid    bool
	}{
		{
			"test-label", true,
		},
		{
			"test_label", true,
		},
		{
			"MixedCase", false,
		},
		{
			"number-09-dash", true,
		},
		{
			"", false,
		},
		{
			"test-label", true,
		},
		{
			"mixed*symbol", false,
		},
		{
			"intérnätional", true,
		},
	}

	for _, tc := range testCases {
		labels := make(map[string]interface{})
		labels[tc.labelKey] = "test value"

		_, errs := labelKeyValidator(labels, "")
		if tc.valid && len(errs) > 0 {
			t.Errorf("Validation failure, key: '%s' should be valid but actual errors were %q", tc.labelKey, errs)
		}
		if !tc.valid && len(errs) < 1 {
			t.Errorf("Validation failure, key: '%s' should fail but actual errors were %q", tc.labelKey, errs)
		}
	}
}

func TestCompareSelfLinkOrResourceNameWithMultipleParts(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"projects to no projects doc": {
			Old:                "projects/myproject/databases/default/documents/resource",
			New:                "resource",
			ExpectDiffSuppress: true,
		},
		"no projects to projects doc": {
			Old:                "resource",
			New:                "projects/myproject/databases/default/documents/resource",
			ExpectDiffSuppress: true,
		},
		"projects to projects doc": {
			Old:                "projects/myproject/databases/default/documents/resource",
			New:                "projects/myproject/databases/default/documents/resource",
			ExpectDiffSuppress: true,
		},
		"multi messages doc": {
			Old:                "messages/{messageId}",
			New:                "projects/myproject/databases/(default)/documents/messages/{messageId}",
			ExpectDiffSuppress: true,
		},
		"multi messages 2 doc": {
			Old:                "projects/myproject/databases/(default)/documents/messages/{messageId}",
			New:                "messages/{messageId}",
			ExpectDiffSuppress: true,
		},
		"projects to no projects topics": {
			Old:                "projects/myproject/topics/resource",
			New:                "resource",
			ExpectDiffSuppress: true,
		},
		"no projects to projects topics": {
			Old:                "resource",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: true,
		},
		"projects to projects topics": {
			Old:                "projects/myproject/topics/resource",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: true,
		},

		"unmatched projects to no projects doc": {
			Old:                "projects/myproject/databases/default/documents/resource",
			New:                "resourcex",
			ExpectDiffSuppress: false,
		},
		"unmatched no projects to projects doc": {
			Old:                "resourcex",
			New:                "projects/myproject/databases/default/documents/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to projects doc": {
			Old:                "projects/myproject/databases/default/documents/resource",
			New:                "projects/myproject/databases/default/documents/resourcex",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to projects 2 doc": {
			Old:                "projects/myprojectx/databases/default/documents/resource",
			New:                "projects/myproject/databases/default/documents/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to empty doc": {
			Old:                "",
			New:                "projects/myproject/databases/default/documents/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched empty to projects 2 doc": {
			Old:                "projects/myprojectx/databases/default/documents/resource",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"unmatched default to default2 doc": {
			Old:                "projects/myproject/databases/default/documents/resource",
			New:                "projects/myproject/databases/default2/documents/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to no projects topics": {
			Old:                "projects/myproject/topics/resource",
			New:                "resourcex",
			ExpectDiffSuppress: false,
		},
		"unmatched no projects to projects topics": {
			Old:                "resourcex",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to projects topics": {
			Old:                "projects/myproject/topics/resource",
			New:                "projects/myproject/topics/resourcex",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to projects 2 topics": {
			Old:                "projects/myprojectx/topics/resource",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched projects to empty topics": {
			Old:                "projects/myproject/topics/resource",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"unmatched empty to projects topics": {
			Old:                "",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: false,
		},
		"unmatched resource to resource-partial": {
			Old:                "resource",
			New:                "resource-partial",
			ExpectDiffSuppress: false,
		},
		"unmatched resource-partial to projects": {
			Old:                "resource-partial",
			New:                "projects/myproject/topics/resource",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if compareSelfLinkOrResourceNameWithMultipleParts("resource", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
