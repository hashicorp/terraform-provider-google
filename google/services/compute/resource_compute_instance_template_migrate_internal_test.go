// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"testing"

	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestComputeInstanceTemplateMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion        int
		BeforeAttributes    map[string]string
		AfterAttributes     map[string]string
		ErrorStringExpected string
		Meta                interface{}
	}{
		"invalid state version": {
			StateVersion:        -1,
			ErrorStringExpected: "Unexpected schema version: -1",
		},
		"simple automatic_restart case removal": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"automatic_restart":              "true",
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
			AfterAttributes: map[string]string{
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
		},
		"simple on_host_maintenance removal": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"on_host_maintenance":            "MIGRATE",
				"automatic_restart":              "true",
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
			AfterAttributes: map[string]string{
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
		},
		"missing scheduling block": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"automatic_restart": "true",
			},
			AfterAttributes: map[string]string{
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
		},
		"empty scheduling block": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"automatic_restart": "true",
				"scheduling.#":      "0",
			},
			AfterAttributes: map[string]string{
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "true",
			},
		},
		"error upon multiple scheduling block": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"automatic_restart":              "true",
				"scheduling.#":                   "2",
				"scheduling.0.automatic_restart": "true",
				"scheduling.1.automatic_restart": "true",
			},
			ErrorStringExpected: "Found multiple scheduling blocks when there should only be one",
		},
		"error upon differing automatic_restart values": {
			StateVersion: 0,
			BeforeAttributes: map[string]string{
				"automatic_restart":              "true",
				"scheduling.#":                   "1",
				"scheduling.0.automatic_restart": "false",
			},
			ErrorStringExpected: "Found differing values for automatic_restart in state, unsure how to proceed.",
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.BeforeAttributes,
		}
		is, err := resourceComputeInstanceTemplateMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			if tc.ErrorStringExpected == "" {
				t.Fatalf("bad: %s, err: %#v", tn, err)
			} else if !strings.Contains(err.Error(), tc.ErrorStringExpected) {
				t.Fatalf("Expected error containing string %s, instead found %#v", tc.ErrorStringExpected, err)
			} else {
				continue
			}
		}

		// Compare both maps for identity
		if !reflect.DeepEqual(is.Attributes, tc.AfterAttributes) {
			t.Fatalf("Expected attributes %#v, got attributes %#v", tc.AfterAttributes, is.Attributes)
		}
	}
}

func TestComputeInstanceTemplateMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta interface{}

	// should handle nil
	is, err := resourceComputeInstanceTemplateMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	_, err = resourceComputeInstanceTemplateMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}
