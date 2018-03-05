package google

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/api/googleapi"
	"reflect"
)

type serializableBody struct {
	body map[string]interface{}

	// ForceSendFields is a list of field names (e.g. "UtilizationTarget")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string

	// NullFields is a list of field names (e.g. "UtilizationTarget") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string
}

// MarshalJSON returns a JSON encoding of schema containing only selected fields.
// A field is selected if any of the following is true:
//   * it has a non-empty value
//   * its field name is present in forceSendFields and it is not a nil pointer or nil interface
//   * its field name is present in nullFields.
func (b *serializableBody) MarshalJSON() ([]byte, error) {
	// By default, all fields in a map are added to the json output
	// This changes that to remove the entry with an empty value.
	// This mimics the "omitempty" behavior.

	// The "omitempty" option specifies that the field should be omitted
	// from the encoding if the field has an empty value, defined as
	// false, 0, a nil pointer, a nil interface value, and any empty array,
	// slice, map, or string.

	// TODO: Add support for ForceSendFields and NullFields.
	for k, v := range b.body {
		if isEmptyValue(reflect.ValueOf(v)) {
			delete(b.body, k)
		}
	}

	return json.Marshal(b.body)
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func Post(config *Config, url string, body map[string]interface{}) (map[string]interface{}, error) {
	return sendRequest(config, "POST", url, body)
}

func Get(config *Config, url string) (map[string]interface{}, error) {
	return sendRequest(config, "GET", url, nil)
}

func Put(config *Config, url string, body map[string]interface{}) (map[string]interface{}, error) {
	return sendRequest(config, "PUT", url, body)
}

func Delete(config *Config, url string) (map[string]interface{}, error) {
	return sendRequest(config, "DELETE", url, nil)
}

func sendRequest(config *Config, method, url string, body map[string]interface{}) (map[string]interface{}, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", config.userAgent)
	reqHeaders.Set("Content-Type", "application/json")

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(&serializableBody{
			body: body})
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url+"?alt=json", &buf)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	res, err := config.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func replaceVars(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = getProject(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = getRegion(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = getZone(d, config)
		if err != nil {
			return "", err
		}
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
		v, ok := d.GetOk(m)
		if ok {
			return v.(string)
		}
		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}
