package google

import "testing"

func TestResourceWithOnlyBaseVersionFields(t *testing.T) {
	d := &ResourceDataMock{
		FieldsInSchema: []string{"normal_field"},
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
		FieldsInSchema: []string{"normal_field", "beta_field"},
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
		FieldsInSchema: []string{"normal_field"},
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
		FieldsInSchema:      []string{"normal_field", "beta_update_field"},
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

type ResourceDataMock struct {
	FieldsInSchema      []string
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
	exists := false
	for _, val := range d.FieldsInSchema {
		if key == val {
			exists = true
		}

	}

	return nil, exists
}
