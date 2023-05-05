package google

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestReplaceVars(t *testing.T) {
	cases := map[string]struct {
		Template      string
		SchemaValues  map[string]interface{}
		Config        *transport_tpg.Config
		Expected      string
		ExpectedError bool
	}{
		"unspecified project fails": {
			Template:      "projects/{{project}}/global/images",
			ExpectedError: true,
		},
		"unspecified region fails": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &transport_tpg.Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"unspecified zone fails": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &transport_tpg.Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"regional with default values": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Region:  "default-region",
			},
			Expected: "projects/default-project/regions/default-region/subnetworks",
		},
		"zonal with default values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Zone:    "default-zone",
			},
			Expected: "projects/default-project/zones/default-zone/instances",
		},
		"regional schema values": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"region":  "region1",
				"name":    "subnetwork1",
			},
			Expected: "projects/project1/regions/region1/subnetworks/subnetwork1",
		},
		"regional schema self-link region": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"region":  "https://www.googleapis.com/compute/v1/projects/project1/regions/region1",
				"name":    "subnetwork1",
			},
			Expected: "projects/project1/regions/region1/subnetworks/subnetwork1",
		},
		"zonal schema values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"zone":    "zone1",
				"name":    "instance1",
			},
			Expected: "projects/project1/zones/zone1/instances/instance1",
		},
		"zonal schema self-link zone": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"zone":    "https://www.googleapis.com/compute/v1/projects/project1/zones/zone1",
				"name":    "instance1",
			},
			Expected: "projects/project1/zones/zone1/instances/instance1",
		},
		"zonal schema recursive replacement": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project":   "project1",
				"zone":      "wrapper{{innerzone}}wrapper",
				"name":      "instance1",
				"innerzone": "inner",
			},
			Expected: "projects/project1/zones/wrapperinnerwrapper/instances/instance1",
		},
		"base path recursive replacement": {
			Template: "{{CloudRunBasePath}}namespaces/{{project}}/services",
			Config: &transport_tpg.Config{
				Project:          "default-project",
				Region:           "default-region",
				CloudRunBasePath: "https://{{region}}-run.googleapis.com/",
			},
			Expected: "https://default-region-run.googleapis.com/namespaces/default-project/services",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			d := &acctest.ResourceDataMock{
				FieldsInSchema: tc.SchemaValues,
			}

			config := tc.Config
			if config == nil {
				config = &transport_tpg.Config{}
			}

			v, err := ReplaceVars(d, config, tc.Template)

			if err != nil {
				if !tc.ExpectedError {
					t.Errorf("bad: %s; unexpected error %s", tn, err)
				}
				return
			}

			if tc.ExpectedError {
				t.Errorf("bad: %s; expected error", tn)
			}

			if v != tc.Expected {
				t.Errorf("bad: %s; expected %q, got %q", tn, tc.Expected, v)
			}
		})
	}
}
