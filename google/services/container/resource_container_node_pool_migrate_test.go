// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestContainerNodePoolMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		ExpectedId   string
		Meta         interface{}
	}{
		"update id from name to zone/cluster/name": {
			StateVersion: 0,
			Attributes: map[string]string{
				"name":    "node-pool-1",
				"zone":    "us-central1-c",
				"cluster": "cluster-1",
			},
			ExpectedId: "us-central1-c/cluster-1/node-pool-1",
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.Attributes["name"],
			Attributes: tc.Attributes,
		}

		is, err := resourceContainerNodePoolMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if is.ID != tc.ExpectedId {
			t.Fatalf("Id should be set to `%s` but is `%s`", tc.ExpectedId, is.ID)
		}
	}
}

func TestContainerNodePoolMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta *transport_tpg.Config

	// should handle nil
	is, err := resourceContainerNodePoolMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	_, err = resourceContainerNodePoolMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
