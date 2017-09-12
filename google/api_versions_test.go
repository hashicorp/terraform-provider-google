package google

import "testing"

func TestResourceWithOnlyBaseVersionFields(t *testing.T) {
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"normal_field": "foo",
		},
	}

	resourceVersion := v1
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{})
	if computeApiVersion != resourceVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", resourceVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{}, []Feature{})
	if computeApiVersion != resourceVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", resourceVersion, computeApiVersion)
	}
}

func TestResourceWithBetaFields(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"normal_field": "foo",
			"beta_field":   "bar",
		},
	}

	expectedVersion := v0beta
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "beta_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "beta_field"}}, []Feature{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithBetaFieldsNotInSchema(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"normal_field": "foo",
		},
	}

	expectedVersion := v1
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "beta_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "beta_field"}}, []Feature{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithBetaUpdateFields(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"normal_field": "foo",
			"beta_field":   "bar",
		},
		FieldsWithHasChange: []string{"beta_update_field"},
	}

	expectedVersion := v1
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	expectedVersion = v0beta
	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{}, []Feature{{Version: expectedVersion, Item: "beta_update_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithOnlyBaseNestedFields(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"list_field.#":              2,
			"list_field.0.normal_field": "foo",
			"list_field.1.normal_field": "bar",
		},
	}

	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{})
	if computeApiVersion != resourceVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", resourceVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{}, []Feature{{Version: resourceVersion, Item: "list_field.*.beta_nested_field"}})
	if computeApiVersion != resourceVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", resourceVersion, computeApiVersion)
	}
}

func TestResourceWithBetaNestedFields(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"list_field.#":                   2,
			"list_field.0.normal_field":      "foo",
			"list_field.1.normal_field":      "bar",
			"list_field.1.beta_nested_field": "baz",
		},
	}

	expectedVersion := v0beta
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "list_field.*.beta_nested_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "list_field.*.beta_nested_field"}}, []Feature{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithBetaDoubleNestedFields(t *testing.T) {
	resourceVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"list_field.#":                                       1,
			"list_field.0.nested_list_field.#":                   1,
			"list_field.0.nested_list_field.0.beta_nested_field": "foo",
		},
	}

	expectedVersion := v0beta
	computeApiVersion := getComputeApiVersion(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "list_field.*.nested_list_field.*.beta_nested_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, resourceVersion, []Feature{{Version: expectedVersion, Item: "list_field.*.nested_list_field.*.beta_nested_field"}}, []Feature{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
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
