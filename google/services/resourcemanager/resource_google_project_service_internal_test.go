// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"testing"
)

func TestProjectServiceServiceValidateFunc(t *testing.T) {
	cases := map[string]struct {
		val                   interface{}
		ExpectValidationError bool
	}{
		"ignoredProjectService": {
			val:                   "dataproc-control.googleapis.com",
			ExpectValidationError: true,
		},
		"bannedProjectService": {
			val:                   "bigquery-json.googleapis.com",
			ExpectValidationError: true,
		},
		"third party API": {
			val:                   "whatever.example.com",
			ExpectValidationError: false,
		},
		"not a domain": {
			val:                   "monitoring",
			ExpectValidationError: true,
		},
		"not a string": {
			val:                   5,
			ExpectValidationError: true,
		},
	}

	for tn, tc := range cases {
		_, errs := validateProjectServiceService(tc.val, "service")
		if tc.ExpectValidationError && len(errs) == 0 {
			t.Errorf("bad: %s, %q passed validation but was expected to fail", tn, tc.val)
		}
		if !tc.ExpectValidationError && len(errs) > 0 {
			t.Errorf("bad: %s, %q failed validation but was expected to pass. errs: %q", tn, tc.val, errs)
		}
	}
}
