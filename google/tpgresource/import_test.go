package tpgresource

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestParseImportId(t *testing.T) {
	regionalIdRegexes := []string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}
	zonalIdRegexes := []string{
		"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}
	multipleNondefaultIdRegexes := []string{
		"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/clusters/(?P<cluster>[^/]+)/nodePools/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<cluster>[^/]+)/(?P<name>[^/]+)",
		"(?P<cluster>[^/]+)/(?P<name>[^/]+)",
	}

	cases := map[string]struct {
		ImportId             string
		IdRegexes            []string
		Config               *transport_tpg.Config
		ExpectedSchemaValues map[string]interface{}
		ExpectError          bool
	}{
		"full self_link": {
			ImportId:  "https://www.googleapis.com/compute/v1/projects/my-project/regions/my-region/subnetworks/my-subnetwork",
			IdRegexes: regionalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"region":  "my-region",
				"name":    "my-subnetwork",
			},
		},
		"relative self_link": {
			ImportId:  "projects/my-project/regions/my-region/subnetworks/my-subnetwork",
			IdRegexes: regionalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"region":  "my-region",
				"name":    "my-subnetwork",
			},
		},
		"short id": {
			ImportId:  "my-project/my-region/my-subnetwork",
			IdRegexes: regionalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"region":  "my-region",
				"name":    "my-subnetwork",
			},
		},
		"short id with default project and region": {
			ImportId: "my-subnetwork",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Region:  "default-region",
			},
			IdRegexes: regionalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"region":  "default-region",
				"name":    "my-subnetwork",
			},
		},
		"short id with default project and zone": {
			ImportId: "my-instance",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Zone:    "default-zone",
			},
			IdRegexes: zonalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"zone":    "default-zone",
				"name":    "my-instance",
			},
		},
		"short id with two nondefault fields with default project and zone": {
			ImportId: "my-cluster/my-node-pool",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Zone:    "default-zone",
			},
			IdRegexes: multipleNondefaultIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"zone":    "default-zone",
				"cluster": "my-cluster",
				"name":    "my-node-pool",
			},
		},
		"short id with default project and region inferred from default zone": {
			ImportId: "my-subnetwork",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Zone:    "us-east1-a",
			},
			IdRegexes: regionalIdRegexes,
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"region":  "us-east1",
				"name":    "my-subnetwork",
			},
		},
		"invalid import id": {
			ImportId:    "i/n/v/a/l/i/d",
			IdRegexes:   regionalIdRegexes,
			ExpectError: true,
		},
		"provider-level defaults not set": {
			ImportId:    "my-subnetwork",
			IdRegexes:   regionalIdRegexes,
			ExpectError: true,
		},
	}

	for tn, tc := range cases {
		d := &ResourceDataMock{
			FieldsInSchema: make(map[string]interface{}),
		}
		d.SetId(tc.ImportId)
		config := tc.Config
		if config == nil {
			config = &transport_tpg.Config{}
		}

		if err := ParseImportId(tc.IdRegexes, d, config); err == nil {
			for k, expectedValue := range tc.ExpectedSchemaValues {
				if v, ok := d.GetOk(k); ok {
					if v != expectedValue {
						t.Errorf("%s failed; Expected value %q for field %q, got %q", tn, expectedValue, k, v)
					}
				} else {
					t.Errorf("%s failed; Expected a value for field %q", tn, k)
				}
			}
		} else if !tc.ExpectError {
			t.Errorf("%s failed; unexpected error: %s", tn, err)
		}
	}
}
