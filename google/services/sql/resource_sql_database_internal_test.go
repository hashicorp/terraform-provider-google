// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"testing"
)

func TestCaseDiffDashSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"PD_HDD": {
			Old:                "PD_HDD",
			New:                "pd-hdd",
			ExpectDiffSuppress: true,
		},
		"PD_SSD": {
			Old:                "PD_SSD",
			New:                "pd-ssd",
			ExpectDiffSuppress: true,
		},
		"pd-hdd": {
			Old:                "pd-hdd",
			New:                "PD_HDD",
			ExpectDiffSuppress: false,
		},
		"pd-ssd": {
			Old:                "pd-ssd",
			New:                "PD_SSD",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if caseDiffDashSuppress(tn, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
