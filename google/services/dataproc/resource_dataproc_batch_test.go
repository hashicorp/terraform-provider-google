// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc

import (
	"testing"
)

func TestCloudDataprocBatchRuntimeConfigVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"old version is empty, new version has a value": {
			Old:                "",
			New:                "2.2.100",
			ExpectDiffSuppress: false,
		},
		"old version is the prefix of the new version": {
			Old:                "2.2",
			New:                "2.2.100",
			ExpectDiffSuppress: true,
		},
		"old version is not the prefix of the new version": {
			Old:                "2.1",
			New:                "2.2.100",
			ExpectDiffSuppress: false,
		},
		"new version is empty, old version has a value": {
			Old:                "2.2.100",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"new version is the prefix of the old version": {
			Old:                "2.2.100",
			New:                "2.2",
			ExpectDiffSuppress: true,
		},
		"new version is not the prefix of the old version": {
			Old:                "2.2.100",
			New:                "2.1",
			ExpectDiffSuppress: false,
		},
		"old version is the same with the new version": {
			Old:                "2.2.100",
			New:                "2.2.100",
			ExpectDiffSuppress: true,
		},
		"both new version and old version are empty string": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		if CloudDataprocBatchRuntimeConfigVersionDiffSuppressFunc(tc.Old, tc.New) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
