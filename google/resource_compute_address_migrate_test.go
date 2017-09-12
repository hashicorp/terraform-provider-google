package google

import (
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestComputeAddressMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		ExpectedId   string
		Meta         interface{}
	}{
		"update id from name to region/name": {
			StateVersion: 0,
			Attributes: map[string]string{
				"name": "address-1",
			},
			ExpectedId: "projects/gcp-project/regions/us-central1/addresses/address-1",
			Meta:       &Config{Region: "us-central1", Project: "gcp-project"},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.Attributes["name"],
			Attributes: tc.Attributes,
		}

		is, err := resourceComputeAddressMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if is.ID != tc.ExpectedId {
			t.Fatalf("Id should be set to `%s` but is `%s`", tc.ExpectedId, is.ID)
		}
	}
}

func TestComputeAddressMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta *Config

	// should handle nil
	is, err := resourceComputeAddressMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	is, err = resourceComputeAddressMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
