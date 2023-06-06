// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package redis

import (
	"testing"
)

func TestSecondaryIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty strings": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"auto range": {
			Old:                "",
			New:                "auto",
			ExpectDiffSuppress: false,
		},
		"auto on already applied range": {
			Old:                "10.0.0.0/28",
			New:                "auto",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.0.0.0/28",
			ExpectDiffSuppress: true,
		},
		"different ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.1.2.3/28",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if secondaryIpDiffSuppress("whatever", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestUnitRedisInstance_redisVersionIsDecreasing(t *testing.T) {
	t.Parallel()
	type testcase struct {
		name       string
		old        interface{}
		new        interface{}
		decreasing bool
	}
	tcs := []testcase{
		{
			name:       "stays the same",
			old:        "REDIS_4_0",
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "increases",
			old:        "REDIS_4_0",
			new:        "REDIS_5_0",
			decreasing: false,
		},
		{
			name:       "nil vals",
			old:        nil,
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "corrupted",
			old:        "REDIS_4_0",
			new:        "REDIS_banana",
			decreasing: false,
		},
		{
			name:       "decreases",
			old:        "REDIS_6_0",
			new:        "REDIS_4_0",
			decreasing: true,
		},
	}

	for _, tc := range tcs {
		decreasing := isRedisVersionDecreasingFunc(tc.old, tc.new)
		if decreasing != tc.decreasing {
			t.Errorf("%s: expected decreasing to be %v, but was %v", tc.name, tc.decreasing, decreasing)
		}
	}
}
