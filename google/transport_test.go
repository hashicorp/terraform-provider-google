package google

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This function isn't a test of transport.go; instead, it is used as an alternative
// to replaceVars inside tests.
func replaceVarsForTest(config *Config, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = rs.Primary.Attributes["region"]
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = rs.Primary.Attributes["zone"]
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider config
		if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}

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
			Config: &Config{
				Project:          "default-project",
				Region:           "default-region",
				CloudRunBasePath: "https://{{region}}-run.googleapis.com/",
			},
			Expected: "https://default-region-run.googleapis.com/namespaces/default-project/services",
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
