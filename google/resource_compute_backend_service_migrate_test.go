package google

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestComputeBackendServiceMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		Attributes         map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"v0 to v1": {
			StateVersion: 0,
			Attributes: map[string]string{
				"backend.#":                         "1",
				"backend.242332812.group":           "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instanceGroups/igName",
				"backend.242332812.balancing_mode":  "UTILIZATION",
				"backend.242332812.max_utilization": "0.8",
			},
			ExpectedAttributes: map[string]string{
				"backend.#":                          "1",
				"backend.2573491210.group":           "https://www.googleapis.com/compute/v1/projects/project_name/zones/zone_name/instances/instanceGroups/igName",
				"backend.2573491210.balancing_mode":  "UTILIZATION",
				"backend.2573491210.max_utilization": "0.8",
			},
			Meta: &Config{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceComputeBackendServiceMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}

		for k, v := range is.Attributes {
			if tc.ExpectedAttributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, tc.ExpectedAttributes[k], k, v, is.Attributes)
			}
		}
	}
}

func TestComputeBackendServiceMigrateState_empty(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
	}{
		"v0": {
			StateVersion: 0,
		},
	}

	for tn, tc := range cases {
		var is *terraform.InstanceState
		var meta *Config

		// should handle nil
		is, err := resourceComputeBackendServiceMigrateState(tc.StateVersion, is, meta)

		if err != nil {
			t.Fatalf("bad %s, err: %#v", tn, err)
		}
		if is != nil {
			t.Fatalf("bad %s, expected nil instancestate, got: %#v", tn, is)
		}

		// should handle non-nil but empty
		is = &terraform.InstanceState{}
		is, err = resourceComputeBackendServiceMigrateState(tc.StateVersion, is, meta)

		if err != nil {
			t.Fatalf("bad %s, err: %#v", tn, err)
		}
	}
}
