package google

import (
	"testing"
)

func TestReplaceVars(t *testing.T) {
	cases := map[string]struct {
		Template      string
		SchemaValues  map[string]interface{}
		Config        *Config
		Expected      string
		ExpectedError bool
	}{
		"unspecified project fails": {
			Template:      "projects/{{project}}/global/images",
			ExpectedError: true,
		},
		"unspecified region fails": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"unspecified zone fails": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"regional with default values": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &Config{
				Project: "default-project",
				Region:  "default-region",
			},
			Expected: "projects/default-project/regions/default-region/subnetworks",
		},
		"zonal with default values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &Config{
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
		"zonal schema values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"zone":    "zone1",
				"name":    "instance1",
			},
			Expected: "projects/project1/zones/zone1/instances/instance1",
		},
	}

	for tn, tc := range cases {
		d := &ResourceDataMock{
			FieldsInSchema: tc.SchemaValues,
		}

		config := tc.Config
		if config == nil {
			config = &Config{}
		}

		v, err := replaceVars(d, config, tc.Template)

		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("bad: %s; unexpected error %s", tn, err)
			}
			continue
		}

		if tc.ExpectedError {
			t.Errorf("bad: %s; expected error", tn)
		}

		if v != tc.Expected {
			t.Errorf("bad: %s; expected %q, got %q", tn, tc.Expected, v)
		}
	}
}
