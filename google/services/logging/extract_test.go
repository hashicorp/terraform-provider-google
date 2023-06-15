// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import "testing"

func TestExtractFieldByPattern(t *testing.T) {
	tests := []struct {
		name             string
		fieldValue       string
		parentFieldValue string
		pattern          string
		expected         string
		hasError         bool
	}{
		{"value-is-set", "my-region", "", "", "my-region", false},
		{"use-regex", "", "projects/my-project/regions/my-region/instances/my-instance", "projects/.*/regions/([a-z09A-Z_-]*)/", "my-region", false},
		{"no-mismatch", "my-region", "projects/my-project/regions/my-region/instances/my-instance", "projects/.*/regions/([a-z09A-Z_-]*)/", "my-region", false},
		{"mismatch", "my-region", "projects/my-project/regions/not-my-region/instances/my-instance", "projects/.*/regions/([a-z09A-Z_-]*)/", "ignored", true},
		{"no-values", "", "", "projects/.*/regions/([a-z09A-Z_-]*)/", "ignored", true},
		{"all-short-form", "my-region", "my-instance", "projects/.*/regions/([a-z09A-Z_-]*)/", "my-region", false},
		{"no-submatch", "my-region", "projects/my-project/regions/my-region/instances/my-instance", "projects/.*/regions/[a-z09A-Z_-]*/", "ignored", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := ExtractFieldByPattern(tc.name, tc.fieldValue, tc.parentFieldValue, tc.pattern)
			if err != nil && !tc.hasError {
				t.Errorf("ValueOnRegexFromField(%v, %v, %v) got error %v", tc.fieldValue, tc.parentFieldValue, tc.pattern, err)
			}

			if !tc.hasError {
				if val != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, val)
				}
			}
		})
	}
}
