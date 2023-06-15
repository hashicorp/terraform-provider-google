// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitBigtableGCPolicy_customizeDiff(t *testing.T) {
	for _, tc := range testUnitBigtableGCPolicyCustomizeDiffTestcases {
		tc.check(t)
	}
}

func (testcase *testUnitBigtableGCPolicyCustomizeDiffTestcase) check(t *testing.T) {
	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{},
		After:  map[string]interface{}{},
	}

	d.Before["max_age.0.days"] = testcase.oldDays
	d.Before["max_age.0.duration"] = testcase.oldDuration

	d.After["max_age.#"] = testcase.arraySize
	d.After["max_age.0.days"] = testcase.newDays
	d.After["max_age.0.duration"] = testcase.newDuration

	err := resourceBigtableGCPolicyCustomizeDiffFunc(d)
	if err != nil {
		t.Errorf("error on testcase %s - %v", testcase.testName, err)
	}

	var cleared bool = d.Cleared != nil && d.Cleared["max_age.0.duration"] == true && d.Cleared["max_age.0.days"] == true
	if cleared != testcase.cleared {
		t.Errorf("%s: expected diff clear to be %v, but was %v", testcase.testName, testcase.cleared, cleared)
	}
}

type testUnitBigtableGCPolicyJSONRules struct {
	name          string
	gcJSONString  string
	want          string
	errorExpected bool
}

var testUnitBigtableGCPolicyRulesTestCases = []testUnitBigtableGCPolicyJSONRules{
	{
		name:          "Simple policy",
		gcJSONString:  `{"rules":[{"max_age":"10h"}]}`,
		want:          "age() > 10h",
		errorExpected: false,
	},
	{
		name:          "Simple multiple policies",
		gcJSONString:  `{"mode":"union", "rules":[{"max_age":"10h"},{"max_version":2}]}`,
		want:          "(age() > 10h || versions() > 2)",
		errorExpected: false,
	},
	{
		name:          "Nested policy",
		gcJSONString:  `{"mode":"union", "rules":[{"max_age":"10h"},{"mode": "intersection", "rules":[{"max_age":"2h"}, {"max_version":2}]}]}`,
		want:          "(age() > 10h || (age() > 2h && versions() > 2))",
		errorExpected: false,
	},
	{
		name:          "JSON with no `rules`",
		gcJSONString:  `{"mode": "union"}`,
		errorExpected: true,
	},
	{
		name:          "Empty JSON",
		gcJSONString:  "{}",
		errorExpected: true,
	},
	{
		name:          "Invalid duration string",
		errorExpected: true,
		gcJSONString:  `{"mode":"union","rules":[{"max_age":"12o"},{"max_version":2}]}`,
	},
	{
		name:          "Empty mode policy with more than 1 rules",
		gcJSONString:  `{"rules":[{"max_age":"10h"}, {"max_version":2}]}`,
		errorExpected: true,
	},
	{
		name:          "Less than 2 rules with mode specified",
		gcJSONString:  `{"mode":"union", "rules":[{"max_version":2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule object",
		gcJSONString:  `{"mode": "union", "rules": [{"mode": "intersection"}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: not max_version or max_age",
		gcJSONString:  `{"mode": "union", "rules": [{"max_versions": 2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: additional fields",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10h", "something_else": 100}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: more than 2 fields in a gc rule object",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10h", "max_version": 10, "something": 100}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: max_version or max_age is in the wrong type",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10d", "max_version": 2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule: wrong data type for child gc_rule",
		gcJSONString:  `{"rules": {"max_version": "456"}}`,
		errorExpected: true,
	},
}

func TestUnitBigtableGCPolicy_getGCPolicyFromJSON(t *testing.T) {
	for _, tc := range testUnitBigtableGCPolicyRulesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var topLevelPolicy map[string]interface{}
			err := json.Unmarshal([]byte(tc.gcJSONString), &topLevelPolicy)
			if err != nil {
				t.Fatalf("error unmarshalling JSON string: %v", err)
			}
			got, err := getGCPolicyFromJSON(topLevelPolicy /*isTopLevel=*/, true)
			if tc.errorExpected && err == nil {
				t.Fatal("expect error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else {
				if got != nil && got.String() != tc.want {
					t.Errorf("error getting policy from JSON, got: %v, want: %v", got, tc.want)
				}
			}
		})
	}
}

type testUnitBigtableGCPolicyCustomizeDiffTestcase struct {
	testName    string
	arraySize   int
	oldDays     int
	newDays     int
	oldDuration string
	newDuration string
	cleared     bool
}

var testUnitBigtableGCPolicyCustomizeDiffTestcases = []testUnitBigtableGCPolicyCustomizeDiffTestcase{
	{
		testName:  "ArraySize0",
		arraySize: 0,
		cleared:   false,
	},
	{
		testName:  "DaysChange",
		arraySize: 1,
		oldDays:   3,
		newDays:   2,
		cleared:   false,
	},
	{
		testName:    "DurationChanges",
		arraySize:   1,
		oldDuration: "3h",
		newDuration: "4h",
		cleared:     false,
	},
	{
		testName:    "DaysToDurationEq",
		arraySize:   1,
		oldDays:     3,
		newDuration: "72h",
		cleared:     true,
	},
	{
		testName:    "DaysToDurationNotEq",
		arraySize:   1,
		oldDays:     3,
		newDuration: "70h",
		cleared:     false,
	},
}
