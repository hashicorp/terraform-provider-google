package google

import "testing"

type ExpectedApiVersions struct {
	Create     ComputeApiVersion
	ReadDelete ComputeApiVersion
	Update     ComputeApiVersion
}

func TestComputeApiVersion(t *testing.T) {
	baseVersion := v1
	betaVersion := v0beta

	cases := map[string]struct {
		Features         []Feature
		FieldsInSchema   map[string]interface{}
		UpdatedFields    []string
		UpdateOnlyFields []Feature
		ExpectedApiVersions
	}{
		"no beta field": {
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     baseVersion,
			},
		},
		"beta field not set": {
			Features: []Feature{{Version: betaVersion, Item: "beta_field"}},
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     baseVersion,
			},
		},
		"beta field set": {
			Features: []Feature{{Version: betaVersion, Item: "beta_field"}},
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
				"beta_field":   "bar",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     betaVersion,
				ReadDelete: betaVersion,
				Update:     betaVersion,
			},
		},
		"update only beta field": {
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
			},
			UpdatedFields:    []string{"beta_update_field"},
			UpdateOnlyFields: []Feature{{Version: betaVersion, Item: "beta_update_field"}},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     betaVersion,
			},
		},
		"nested beta field not set": {
			Features: []Feature{{Version: betaVersion, Item: "list_field.*.beta_nested_field"}},
			FieldsInSchema: map[string]interface{}{
				"list_field.#":              2,
				"list_field.0.normal_field": "foo",
				"list_field.1.normal_field": "bar",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     baseVersion,
			},
		},
		"nested beta field set": {
			Features: []Feature{{Version: betaVersion, Item: "list_field.*.beta_nested_field"}},
			FieldsInSchema: map[string]interface{}{
				"list_field.#":                   2,
				"list_field.0.normal_field":      "foo",
				"list_field.1.normal_field":      "bar",
				"list_field.1.beta_nested_field": "baz",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     betaVersion,
				ReadDelete: betaVersion,
				Update:     betaVersion,
			},
		},
		"double nested fields set": {
			Features: []Feature{{Version: betaVersion, Item: "list_field.*.nested_list_field.*.beta_nested_field"}},
			FieldsInSchema: map[string]interface{}{
				"list_field.#":                                       1,
				"list_field.0.nested_list_field.#":                   1,
				"list_field.0.nested_list_field.0.beta_nested_field": "foo",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     betaVersion,
				ReadDelete: betaVersion,
				Update:     betaVersion,
			},
		},
		"beta field has default value": {
			Features: []Feature{{Version: betaVersion, Item: "beta_field", DefaultValue: "bar"}},
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
				"beta_field":   "bar",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     baseVersion,
			},
		},
		"beta field is updated to default value": {
			Features: []Feature{{Version: betaVersion, Item: "beta_field", DefaultValue: "bar"}},
			FieldsInSchema: map[string]interface{}{
				"normal_field": "foo",
				"beta_field":   "bar",
			},
			UpdatedFields: []string{"beta_field"},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     betaVersion,
			},
		},
		"nested beta field has default value": {
			Features: []Feature{{Version: betaVersion, Item: "list_field.*.beta_nested_field", DefaultValue: "baz"}},
			FieldsInSchema: map[string]interface{}{
				"list_field.#":                   2,
				"list_field.0.normal_field":      "foo",
				"list_field.1.normal_field":      "bar",
				"list_field.1.beta_nested_field": "baz",
			},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     baseVersion,
			},
		},
		"nested beta field is updated default value": {
			Features: []Feature{{Version: betaVersion, Item: "list_field.*.beta_nested_field", DefaultValue: "baz"}},
			FieldsInSchema: map[string]interface{}{
				"list_field.#":                   2,
				"list_field.0.normal_field":      "foo",
				"list_field.1.normal_field":      "bar",
				"list_field.1.beta_nested_field": "baz",
			},
			UpdatedFields: []string{"list_field.1.beta_nested_field"},
			ExpectedApiVersions: ExpectedApiVersions{
				Create:     baseVersion,
				ReadDelete: baseVersion,
				Update:     betaVersion,
			},
		},
	}

	for tn, tc := range cases {
		// Create
		// All fields with value have HasChange set to true.
		keys := make([]string, 0, len(tc.FieldsInSchema))
		for key := range tc.FieldsInSchema {
			keys = append(keys, key)
		}

		d := &ResourceDataMock{
			FieldsInSchema:      tc.FieldsInSchema,
			FieldsWithHasChange: keys,
		}

		apiVersion := getComputeApiVersion(d, v1, tc.Features)
		if apiVersion != tc.ExpectedApiVersions.Create {
			t.Errorf("bad: %s, Expected to see version %v for create, got version %v", tn, tc.ExpectedApiVersions.Create, apiVersion)
		}

		// Read/Delete
		// All fields have HasChange set to false.
		d = &ResourceDataMock{
			FieldsInSchema: tc.FieldsInSchema,
		}

		apiVersion = getComputeApiVersion(d, v1, tc.Features)
		if apiVersion != tc.ExpectedApiVersions.ReadDelete {
			t.Errorf("bad: %s, Expected to see version %v for read/delete, got version %v", tn, tc.ExpectedApiVersions.ReadDelete, apiVersion)
		}

		// Update
		// Only fields defined as updated in the test case have HasChange set to true.
		d = &ResourceDataMock{
			FieldsInSchema:      tc.FieldsInSchema,
			FieldsWithHasChange: tc.UpdatedFields,
		}

		apiVersion = getComputeApiVersionUpdate(d, v1, tc.Features, tc.UpdateOnlyFields)
		if apiVersion != tc.ExpectedApiVersions.Update {
			t.Errorf("bad: %s, Expected to see version %v for update, got version %v", tn, tc.ExpectedApiVersions.Update, apiVersion)
		}
	}
}

type ResourceDataMock struct {
	FieldsInSchema      map[string]interface{}
	FieldsWithHasChange []string
}

func (d *ResourceDataMock) HasChange(key string) bool {
	exists := false
	for _, val := range d.FieldsWithHasChange {
		if key == val {
			exists = true
		}
	}

	return exists
}

func (d *ResourceDataMock) GetOk(key string) (interface{}, bool) {
	for k, v := range d.FieldsInSchema {
		if key == k {
			return v, true
		}
	}

	return nil, false
}
