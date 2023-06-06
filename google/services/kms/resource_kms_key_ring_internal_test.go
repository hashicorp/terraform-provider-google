// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestKeyRingIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedKeyRingId   string
		Config              *transport_tpg.Config
	}{
		"id is in project/location/keyRingName format": {
			ImportId:            "test-project/us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/test-project/locations/us-central1/keyRings/test-key-ring",
		},
		"id is in domain:project/location/keyRingName format": {
			ImportId:            "example.com:test-project/us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/example.com:test-project/locations/us-central1/keyRings/test-key-ring",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/can-you-believe-that-this-key-ring-name-is-exactly-64-characters",
			ExpectedError: true,
		},
		"id is in location/keyRingName format": {
			ImportId:            "us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/test-project/locations/us-central1/keyRings/test-key-ring",
			Config:              &transport_tpg.Config{Project: "test-project"},
		},
		"id is in location/keyRingName format without project in config": {
			ImportId:      "us-central1/test-key-ring",
			ExpectedError: true,
			Config:        &transport_tpg.Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		keyRingId, err := parseKmsKeyRingId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if keyRingId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, keyRingId.TerraformId())
		}

		if keyRingId.KeyRingId() != tc.ExpectedKeyRingId {
			t.Fatalf("bad: %s, expected KeyRing ID to be `%s` but is `%s`", tn, tc.ExpectedKeyRingId, keyRingId.KeyRingId())
		}
	}
}
