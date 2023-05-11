package tpgresource

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
	if ok && !IsEmptyValue(reflect.ValueOf(v)) {
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

func (d *ResourceDataMock) GetProviderMeta(dst interface{}) error {
	return nil
}

func (d *ResourceDataMock) Timeout(key string) time.Duration {
	return time.Duration(1)
}

type ResourceDiffMock struct {
	Before     map[string]interface{}
	After      map[string]interface{}
	Cleared    map[string]interface{}
	IsForceNew bool
}

func (d *ResourceDiffMock) GetChange(key string) (interface{}, interface{}) {
	return d.Before[key], d.After[key]
}

func (d *ResourceDiffMock) HasChange(key string) bool {
	old, new := d.GetChange(key)
	return old != new
}

func (d *ResourceDiffMock) Get(key string) interface{} {
	return d.After[key]
}

func (d *ResourceDiffMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.After[key]
	return v, ok
}

func (d *ResourceDiffMock) Clear(key string) error {
	if d.Cleared == nil {
		d.Cleared = map[string]interface{}{}
	}
	d.Cleared[key] = true
	return nil
}

func (d *ResourceDiffMock) ForceNew(key string) error {
	d.IsForceNew = true
	return nil
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func ReplaceVarsForTest(config *transport_tpg.Config, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
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

// Used to create populated schema.ResourceData structs in tests.
// Pass in a schema and a config map containing the fields and values you wish to set
// The returned schema.ResourceData can represent a configured resource, data source or provider.
func SetupTestResourceDataFromConfigMap(t *testing.T, s map[string]*schema.Schema, configValues map[string]interface{}) *schema.ResourceData {

	// Create empty schema.ResourceData using the SDK Provider schema
	emptyConfigMap := map[string]interface{}{}
	d := schema.TestResourceDataRaw(t, s, emptyConfigMap)

	// Load Terraform config data
	if len(configValues) > 0 {
		for k, v := range configValues {
			err := d.Set(k, v)
			if err != nil {
				t.Fatalf("error during test setup: %v", err)
			}
		}
	}

	return d
}
