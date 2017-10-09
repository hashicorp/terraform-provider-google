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
