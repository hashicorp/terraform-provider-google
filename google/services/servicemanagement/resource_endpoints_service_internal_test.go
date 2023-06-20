// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicemanagement

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestEndpointsService_grpcMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		Attributes         map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"update from protoc_output to protoc_output_base64": {
			StateVersion: 0,
			Attributes: map[string]string{
				"protoc_output": "123456789",
				"name":          "testcase",
			},
			ExpectedAttributes: map[string]string{
				"protoc_output_base64": "MTIzNDU2Nzg5",
				"protoc_output":        "",
				"name":                 "testcase",
			},
			Meta: &transport_tpg.Config{Project: "gcp-project", Region: "us-central1"},
		},
		"update from non-protoc_output": {
			StateVersion: 0,
			Attributes: map[string]string{
				"openapi_config": "foo bar baz",
				"name":           "testcase-2",
			},
			ExpectedAttributes: map[string]string{
				"openapi_config": "foo bar baz",
				"name":           "testcase-2",
			},
			Meta: &transport_tpg.Config{Project: "gcp-project", Region: "us-central1"},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.Attributes["name"],
			Attributes: tc.Attributes,
		}

		is, err := migrateEndpointsService(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if !reflect.DeepEqual(is.Attributes, tc.ExpectedAttributes) {
			t.Fatalf("Attributes should be `%s` but are `%s`", tc.ExpectedAttributes, is.Attributes)
		}
	}
}
