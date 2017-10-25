package google

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestContainerClusterMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"change additional_zones from list to set": {
			StateVersion: 0,
			Attributes: map[string]string{
				"additional_zones.#": "2",
				"additional_zones.0": "us-central1-c",
				"additional_zones.1": "us-central1-b",
			},
			Expected: map[string]string{
				"additional_zones.#":          "2",
				"additional_zones.90274510":   "us-central1-c",
				"additional_zones.1919306328": "us-central1-b",
			},
			Meta: &Config{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceContainerClusterMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestContainerClusterMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta *Config

	// should handle nil
	is, err := resourceContainerClusterMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	is, err = resourceContainerClusterMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
