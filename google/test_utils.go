package google

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type ResourceDataMock struct {
	FieldsInSchema      map[string]interface{}
	FieldsWithHasChange []string
	id                  string
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

func (d *ResourceDataMock) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

func (d *ResourceDataMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.GetOkExists(key)
	if ok && !isEmptyValue(reflect.ValueOf(v)) {
		return v, true
	} else {
		return v, false
	}
}

func (d *ResourceDataMock) GetOkExists(key string) (interface{}, bool) {
	for k, v := range d.FieldsInSchema {
		if key == k {
			return v, true
		}
	}

	return nil, false
}

func (d *ResourceDataMock) Set(key string, value interface{}) error {
	d.FieldsInSchema[key] = value
	return nil
}

func (d *ResourceDataMock) SetId(v string) {
	d.id = v
}

func (d *ResourceDataMock) Id() string {
	return d.id
}

type ResourceDiffMock struct {
	Before  map[string]interface{}
	After   map[string]interface{}
	Cleared map[string]struct{}
}

func (d *ResourceDiffMock) GetChange(key string) (interface{}, interface{}) {
	return d.Before[key], d.After[key]
}

func (d *ResourceDiffMock) Clear(key string) error {
	if d.Cleared == nil {
		d.Cleared = map[string]struct{}{}
	}
	d.Cleared[key] = struct{}{}
	return nil
}

func toBool(attribute string) (bool, error) {
	// Handle the case where an unset value defaults to false
	if attribute == "" {
		return false, nil
	}
	return strconv.ParseBool(attribute)
}

func checkDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return checkDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, map[string]struct{}{})
}

func checkDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		errMsg := ""
		// Data sources are often derived from resources, so iterate over the resource fields to
		// make sure all fields are accounted for in the data source.
		// If a field exists in the data source but not in the resource, its expected value should
		// be checked separately.
		for k := range rsAttr {
			if _, ok := ignoreFields[k]; ok {
				continue
			}
			if dsAttr[k] != rsAttr[k] {
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], rsAttr[k])
			}
		}

		if errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		return nil
	}
}
