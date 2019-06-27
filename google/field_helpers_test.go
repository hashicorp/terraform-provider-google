package google

import (
	"testing"
)

func TestParseGlobalFieldValue(t *testing.T) {
	const resourceType = "networks"
	cases := map[string]struct {
		FieldValue           string
		ExpectedRelativeLink string
		ExpectedError        bool
		IsEmptyValid         bool
		ProjectSchemaField   string
		ProjectSchemaValue   string
		Config               *Config
	}{
		"network is a full self link": {
			FieldValue:           "https://www.googleapis.com/compute/v1/projects/myproject/global/networks/my-network",
			ExpectedRelativeLink: "projects/myproject/global/networks/my-network",
		},
		"network is a relative self link": {
			FieldValue:           "projects/myproject/global/networks/my-network",
			ExpectedRelativeLink: "projects/myproject/global/networks/my-network",
		},
		"network is a partial relative self link": {
			FieldValue:           "global/networks/my-network",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/global/networks/my-network",
		},
		"network is the name only": {
			FieldValue:           "my-network",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/global/networks/my-network",
		},
		"network is the name only and has a project set in schema": {
			FieldValue:           "my-network",
			ProjectSchemaField:   "project",
			ProjectSchemaValue:   "schema-project",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/schema-project/global/networks/my-network",
		},
		"network is the name only and has a project set in schema but the field is not specified.": {
			FieldValue:           "my-network",
			ProjectSchemaValue:   "schema-project",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/global/networks/my-network",
		},
		"network is empty and it is valid": {
			FieldValue:           "",
			IsEmptyValid:         true,
			ExpectedRelativeLink: "",
		},
		"network is empty and it is not valid": {
			FieldValue:    "",
			IsEmptyValid:  false,
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		fieldsInSchema := make(map[string]interface{})

		if len(tc.ProjectSchemaValue) > 0 && len(tc.ProjectSchemaField) > 0 {
			fieldsInSchema[tc.ProjectSchemaField] = tc.ProjectSchemaValue
		}

		d := &ResourceDataMock{
			FieldsInSchema: fieldsInSchema,
		}

		v, err := parseGlobalFieldValue(resourceType, tc.FieldValue, tc.ProjectSchemaField, d, tc.Config, tc.IsEmptyValid)

		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("bad: %s, did not expect an error. Error: %s", tn, err)
			}
		} else {
			if v.RelativeLink() != tc.ExpectedRelativeLink {
				t.Errorf("bad: %s, expected relative link to be '%s' but got '%s'", tn, tc.ExpectedRelativeLink, v.RelativeLink())
			}
		}
	}
}

func TestParseZonalFieldValue(t *testing.T) {
	const resourceType = "instances"
	cases := map[string]struct {
		FieldValue           string
		ExpectedRelativeLink string
		ExpectedError        bool
		IsEmptyValid         bool
		ProjectSchemaField   string
		ProjectSchemaValue   string
		ZoneSchemaField      string
		ZoneSchemaValue      string
		Config               *Config
	}{
		"instance is a full self link": {
			FieldValue:           "https://www.googleapis.com/compute/v1/projects/myproject/zones/us-central1-b/instances/my-instance",
			ExpectedRelativeLink: "projects/myproject/zones/us-central1-b/instances/my-instance",
		},
		"instance is a relative self link": {
			FieldValue:           "projects/myproject/zones/us-central1-b/instances/my-instance",
			ExpectedRelativeLink: "projects/myproject/zones/us-central1-b/instances/my-instance",
		},
		"instance is a partial relative self link": {
			FieldValue:           "zones/us-central1-b/instances/my-instance",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/zones/us-central1-b/instances/my-instance",
		},
		"instance is the name only": {
			FieldValue:           "my-instance",
			ZoneSchemaField:      "zone",
			ZoneSchemaValue:      "us-east1-a",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/zones/us-east1-a/instances/my-instance",
		},
		"instance is the name only and has a project set in schema": {
			FieldValue:           "my-instance",
			ProjectSchemaField:   "project",
			ProjectSchemaValue:   "schema-project",
			ZoneSchemaField:      "zone",
			ZoneSchemaValue:      "us-east1-a",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/schema-project/zones/us-east1-a/instances/my-instance",
		},
		"instance is the name only and has a project set in schema but the field is not specified.": {
			FieldValue:           "my-instance",
			ProjectSchemaValue:   "schema-project",
			ZoneSchemaField:      "zone",
			ZoneSchemaValue:      "us-east1-a",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/zones/us-east1-a/instances/my-instance",
		},
		"instance is the name only and no zone field is specified": {
			FieldValue:    "my-instance",
			Config:        &Config{Project: "default-project"},
			ExpectedError: true,
		},
		"instance is the name only and no value for zone field is specified": {
			FieldValue:      "my-instance",
			ZoneSchemaField: "zone",
			Config:          &Config{Project: "default-project"},
			ExpectedError:   true,
		},
		"instance is empty and it is valid": {
			FieldValue:           "",
			IsEmptyValid:         true,
			ExpectedRelativeLink: "",
		},
		"instance is empty and it is not valid": {
			FieldValue:    "",
			IsEmptyValid:  false,
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		fieldsInSchema := make(map[string]interface{})

		if len(tc.ProjectSchemaValue) > 0 && len(tc.ProjectSchemaField) > 0 {
			fieldsInSchema[tc.ProjectSchemaField] = tc.ProjectSchemaValue
		}

		if len(tc.ZoneSchemaValue) > 0 && len(tc.ZoneSchemaField) > 0 {
			fieldsInSchema[tc.ZoneSchemaField] = tc.ZoneSchemaValue
		}

		d := &ResourceDataMock{
			FieldsInSchema: fieldsInSchema,
		}

		v, err := parseZonalFieldValue(resourceType, tc.FieldValue, tc.ProjectSchemaField, tc.ZoneSchemaField, d, tc.Config, tc.IsEmptyValid)

		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("bad: %s, did not expect an error. Error: %s", tn, err)
			}
		} else {
			if v.RelativeLink() != tc.ExpectedRelativeLink {
				t.Errorf("bad: %s, expected relative link to be '%s' but got '%s'", tn, tc.ExpectedRelativeLink, v.RelativeLink())
			}
		}
	}
}

func TestParseOrganizationFieldValue(t *testing.T) {
	const resourceType = "roles"
	cases := map[string]struct {
		FieldValue           string
		ExpectedRelativeLink string
		ExpectedName         string
		ExpectedOrgId        string
		ExpectedError        bool
		IsEmptyValid         bool
	}{
		"role is valid": {
			FieldValue:           "organizations/123/roles/custom",
			ExpectedRelativeLink: "organizations/123/roles/custom",
			ExpectedName:         "custom",
			ExpectedOrgId:        "123",
		},
		"role is empty and it is valid": {
			FieldValue:           "",
			IsEmptyValid:         true,
			ExpectedRelativeLink: "",
		},
		"role is empty and it is not valid": {
			FieldValue:    "",
			IsEmptyValid:  false,
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		v, err := parseOrganizationFieldValue(resourceType, tc.FieldValue, tc.IsEmptyValid)

		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("bad: %s, did not expect an error. Error: %s", tn, err)
			}
		} else {
			if v.RelativeLink() != tc.ExpectedRelativeLink {
				t.Errorf("bad: %s, expected relative link to be '%s' but got '%s'", tn, tc.ExpectedRelativeLink, v.RelativeLink())
			}
		}
	}
}

func TestParseRegionalFieldValue(t *testing.T) {
	const resourceType = "subnetworks"
	cases := map[string]struct {
		FieldValue           string
		ExpectedRelativeLink string
		ExpectedError        bool
		IsEmptyValid         bool
		ProjectSchemaField   string
		ProjectSchemaValue   string
		RegionSchemaField    string
		RegionSchemaValue    string
		ZoneSchemaField      string
		ZoneSchemaValue      string
		Config               *Config
	}{
		"subnetwork is a full self link": {
			FieldValue:           "https://www.googleapis.com/compute/v1/projects/myproject/regions/us-central1/subnetworks/my-subnetwork",
			ExpectedRelativeLink: "projects/myproject/regions/us-central1/subnetworks/my-subnetwork",
		},
		"subnetwork is a relative self link": {
			FieldValue:           "projects/myproject/regions/us-central1/subnetworks/my-subnetwork",
			ExpectedRelativeLink: "projects/myproject/regions/us-central1/subnetworks/my-subnetwork",
		},
		"subnetwork is a partial relative self link": {
			FieldValue:           "regions/us-central1/subnetworks/my-subnetwork",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/default-project/regions/us-central1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only": {
			FieldValue:           "my-subnetwork",
			RegionSchemaField:    "region",
			RegionSchemaValue:    "us-east1",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/regions/us-east1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and has a project set in schema": {
			FieldValue:           "my-subnetwork",
			ProjectSchemaField:   "project",
			ProjectSchemaValue:   "schema-project",
			RegionSchemaField:    "region",
			RegionSchemaValue:    "us-east1",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/schema-project/regions/us-east1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and has a project set in schema but the field is not specified.": {
			FieldValue:           "my-subnetwork",
			ProjectSchemaValue:   "schema-project",
			RegionSchemaField:    "region",
			RegionSchemaValue:    "us-east1",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/default-project/regions/us-east1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and region is extracted from the one field.": {
			FieldValue:           "my-subnetwork",
			ProjectSchemaValue:   "schema-project",
			RegionSchemaField:    "region",
			ZoneSchemaField:      "zone",
			ZoneSchemaValue:      "us-central1-a",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/default-project/regions/us-central1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and region is extracted from the provider-level zone.": {
			FieldValue:           "my-subnetwork",
			ProjectSchemaValue:   "schema-project",
			RegionSchemaField:    "region",
			ZoneSchemaField:      "zone",
			Config:               &Config{Project: "default-project", Zone: "us-central1-c"},
			ExpectedRelativeLink: "projects/default-project/regions/us-central1/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and no region field is specified": {
			FieldValue:           "my-subnetwork",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/default-project/regions/default-region/subnetworks/my-subnetwork",
		},
		"subnetwork is the name only and no value for region field is specified": {
			FieldValue:           "my-subnetwork",
			RegionSchemaField:    "region",
			Config:               &Config{Project: "default-project", Region: "default-region"},
			ExpectedRelativeLink: "projects/default-project/regions/default-region/subnetworks/my-subnetwork",
		},
		"subnetwork is empty and it is valid": {
			FieldValue:           "",
			IsEmptyValid:         true,
			ExpectedRelativeLink: "",
		},
		"subnetwork is empty and it is not valid": {
			FieldValue:    "",
			IsEmptyValid:  false,
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			fieldsInSchema := make(map[string]interface{})

			if tc.ProjectSchemaValue != "" && tc.ProjectSchemaField != "" {
				fieldsInSchema[tc.ProjectSchemaField] = tc.ProjectSchemaValue
			}

			if tc.RegionSchemaValue != "" && tc.RegionSchemaField != "" {
				fieldsInSchema[tc.RegionSchemaField] = tc.RegionSchemaValue
			}
			if tc.ZoneSchemaValue != "" && tc.ZoneSchemaField != "" {
				fieldsInSchema[tc.ZoneSchemaField] = tc.ZoneSchemaValue
			}

			d := &ResourceDataMock{
				FieldsInSchema: fieldsInSchema,
			}

			v, err := parseRegionalFieldValue(resourceType, tc.FieldValue, tc.ProjectSchemaField, tc.RegionSchemaField, tc.ZoneSchemaField, d, tc.Config, tc.IsEmptyValid)

			if err != nil {
				if !tc.ExpectedError {
					t.Errorf("bad: did not expect an error. Error: %s", err)
				}
			} else {
				if v.RelativeLink() != tc.ExpectedRelativeLink {
					t.Errorf("bad: expected relative link to be '%s' but got '%s'", tc.ExpectedRelativeLink, v.RelativeLink())
				}
			}
		})
	}
}

func TestParseProjectFieldValue(t *testing.T) {
	const resourceType = "instances"
	cases := map[string]struct {
		FieldValue           string
		ExpectedRelativeLink string
		ExpectedError        bool
		IsEmptyValid         bool
		ProjectSchemaField   string
		ProjectSchemaValue   string
		Config               *Config
	}{
		"instance is a full self link": {
			FieldValue:           "https://www.googleapis.com/compute/v1/projects/myproject/instances/my-instance",
			ExpectedRelativeLink: "projects/myproject/instances/my-instance",
		},
		"instance is a relative self link": {
			FieldValue:           "projects/myproject/instances/my-instance",
			ExpectedRelativeLink: "projects/myproject/instances/my-instance",
		},
		"instance is a partial relative self link": {
			FieldValue:           "projects/instances/my-instance",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/instances/my-instance",
		},
		"instance is the name only": {
			FieldValue:           "my-instance",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/instances/my-instance",
		},
		"instance is the name only and has a project set in schema": {
			FieldValue:           "my-instance",
			ProjectSchemaField:   "project",
			ProjectSchemaValue:   "schema-project",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/schema-project/instances/my-instance",
		},
		"instance is the name only and has a project set in schema but the field is not specified.": {
			FieldValue:           "my-instance",
			ProjectSchemaValue:   "schema-project",
			Config:               &Config{Project: "default-project"},
			ExpectedRelativeLink: "projects/default-project/instances/my-instance",
		},
		"instance is empty and it is valid": {
			FieldValue:           "",
			IsEmptyValid:         true,
			ExpectedRelativeLink: "",
		},
		"instance is empty and it is not valid": {
			FieldValue:    "",
			IsEmptyValid:  false,
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		fieldsInSchema := make(map[string]interface{})

		if len(tc.ProjectSchemaValue) > 0 && len(tc.ProjectSchemaField) > 0 {
			fieldsInSchema[tc.ProjectSchemaField] = tc.ProjectSchemaValue
		}

		d := &ResourceDataMock{
			FieldsInSchema: fieldsInSchema,
		}

		v, err := parseProjectFieldValue(resourceType, tc.FieldValue, tc.ProjectSchemaField, d, tc.Config, tc.IsEmptyValid)

		if err != nil {
			if !tc.ExpectedError {
				t.Errorf("bad: %s, did not expect an error. Error: %s", tn, err)
			}
		} else {
			if v.RelativeLink() != tc.ExpectedRelativeLink {
				t.Errorf("bad: %s, expected relative link to be '%s' but got '%s'", tn, tc.ExpectedRelativeLink, v.RelativeLink())
			}
		}
	}
}
