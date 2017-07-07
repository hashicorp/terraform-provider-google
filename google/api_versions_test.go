package google

import "testing"

func TestResourceWithOnlyBaseVersionFields(t *testing.T) {
	d := &ResourceDataMock{
		FieldsInSchema: []string{"normal_field"},
	}

	baseVersion := v1
	computeApiVersion := getComputeApiVersion(d, baseVersion, []Key{})
	if computeApiVersion != baseVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", baseVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, baseVersion, []Key{}, []Key{})
	if computeApiVersion != baseVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", baseVersion, computeApiVersion)
	}
}

func TestResourceWithBetaFields(t *testing.T) {
	baseVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: []string{"normal_field", "beta_field"},
	}

	expectedVersion := v0beta
	computeApiVersion := getComputeApiVersion(d, baseVersion, []Key{{Version: expectedVersion, Item: "beta_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, baseVersion, []Key{{Version: expectedVersion, Item: "beta_field"}}, []Key{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithBetaFieldsNotInSchema(t *testing.T) {
	baseVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema: []string{"normal_field"},
	}

	expectedVersion := v1
	computeApiVersion := getComputeApiVersion(d, baseVersion, []Key{{Version: expectedVersion, Item: "beta_field"}})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	computeApiVersion = getComputeApiVersionUpdate(d, baseVersion, []Key{{Version: expectedVersion, Item: "beta_field"}}, []Key{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}
}

func TestResourceWithBetaUpdateFields(t *testing.T) {
	baseVersion := v1
	d := &ResourceDataMock{
		FieldsInSchema:      []string{"normal_field", "beta_update_field"},
		FieldsWithHasChange: []string{"beta_update_field"},
	}

	expectedVersion := v1
	computeApiVersion := getComputeApiVersion(d, baseVersion, []Key{})
	if computeApiVersion != expectedVersion {
		t.Errorf("Expected to see version: %v. Saw version: %v.", expectedVersion, computeApiVersion)
	}

	expectedVersion = v0beta
	computeApiVersion = getComputeApiVersionUpdate(d, baseVersion, []Key{}, []Key{{Version: expectedVersion, Item: "beta_update_field"}})
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
	for _, val := range d.FieldsInSchema {
		if key == val {
			exists = true
		}
	}

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
