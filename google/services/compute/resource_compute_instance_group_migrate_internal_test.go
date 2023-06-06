// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestComputeInstanceGroupMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		Attributes         map[string]string
		ExpectedAttributes map[string]string
		ExpectedId         string
		Meta               interface{}
	}{
		"v1 to v2": {
			StateVersion: 1,
			Attributes: map[string]string{
				"zone": "us-central1-c",
				"name": "instancegroup-test",
			},
			ExpectedAttributes: map[string]string{
				"zone": "us-central1-c",
				"name": "instancegroup-test",
			},
			ExpectedId: "us-central1-c/instancegroup-test",
			Meta:       &transport_tpg.Config{},
		},
		"v0 to v2": {
			StateVersion: 0,
			Attributes: map[string]string{
				"zone":        "us-central1-c",
				"name":        "instancegroup-test",
				"instances.#": "1",
				"instances.0": "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instancegroup-test-1",
				"instances.1": "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instancegroup-test-0",
			},
			ExpectedAttributes: map[string]string{
				"zone":                 "us-central1-c",
				"name":                 "instancegroup-test",
				"instances.#":          "1",
				"instances.764135222":  "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instancegroup-test-1",
				"instances.1519187872": "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instancegroup-test-0",
			},
			ExpectedId: "us-central1-c/instancegroup-test",
			Meta:       &transport_tpg.Config{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceComputeInstanceGroupMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if is.ID != tc.ExpectedId {
			t.Fatalf("bad: %s\n\n expected: %s\n got: %s", tn, tc.ExpectedId, is.ID)
		}

		for k, v := range tc.ExpectedAttributes {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestComputeInstanceGroupMigrateState_empty(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
	}{
		"v0": {
			StateVersion: 0,
		},
		"v1": {
			StateVersion: 1,
		},
	}

	for tn, tc := range cases {
		var is *terraform.InstanceState
		var meta *transport_tpg.Config

		// should handle nil
		is, err := resourceComputeInstanceGroupMigrateState(tc.StateVersion, is, meta)

		if err != nil {
			t.Fatalf("bad %s, err: %#v", tn, err)
		}
		if is != nil {
			t.Fatalf("bad %s, expected nil instancestate, got: %#v", tn, is)
		}

		// should handle non-nil but empty
		is = &terraform.InstanceState{}
		_, err = resourceComputeInstanceGroupMigrateState(tc.StateVersion, is, meta)

		if err != nil {
			t.Fatalf("bad %s, err: %#v", tn, err)
		}
	}
}
