// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"testing"
)

func TestHasTerminationTimeChanged(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		Old, New map[string]interface{}
		Expect   bool
	}{
		"empty": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": ""},
			Expect: false,
		},
		"new": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"changed": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"same": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			Expect: false,
		},
	}
	for tn, tc := range cases {
		if hasTerminationTimeChanged(tc.Old, tc.New) != tc.Expect {
			t.Errorf("%s: expected %t for whether termination time matched for old = %q, new = %q", tn, tc.Expect, tc.Old, tc.New)
		}
	}
}
