package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
)

var DefaultRequestTimeout = 5 * time.Minute

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

func sendRequest(config *Config, method, project, rawurl string, body map[string]interface{}, errorRetryPredicates ...func(e error) (bool, string)) (map[string]interface{}, error) {
	return sendRequestWithTimeout(config, method, project, rawurl, body, DefaultRequestTimeout, errorRetryPredicates...)
}

func sendRequestWithTimeout(config *Config, method, project, rawurl string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...func(e error) (bool, string)) (map[string]interface{}, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", config.userAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if config.UserProjectOverride && project != "" {
		// Pass the project into this fn instead of parsing it from the URL because
		// both project names and URLs can have colons in them.
		reqHeaders.Set("X-Goog-User-Project", project)
	}

	if timeout == 0 {
		timeout = time.Duration(1) * time.Hour
	}

	var res *http.Response
	err := retryTimeDuration(
		func() error {
			var buf bytes.Buffer
			if body != nil {
				err := json.NewEncoder(&buf).Encode(body)
				if err != nil {
					return err
				}
			}

			u, err := addQueryParams(rawurl, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}
			req, err := http.NewRequest(method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = config.client.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return err
			}

			return nil
		},
		timeout,
		errorRetryPredicates...,
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("Unable to parse server response. This is most likely a terraform problem, please file a bug at https://github.com/terraform-providers/terraform-provider-google/issues.")
	}

	// The defer call must be made outside of the retryFunc otherwise it's closed too soon.
	defer googleapi.CloseBody(res)

	// 204 responses will have no body, so we're going to error with "EOF" if we
	// try to parse it. Instead, we can just return nil.
	if res.StatusCode == 204 {
		return nil, nil
	}
	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func addQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func replaceVars(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, 0)
}

// replaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func replaceVarsRecursive(d TerraformResourceData, config *Config, linkTmpl string, depth int) (string, error) {
	if depth > 10 {
		return "", errors.New("Recursive substitution detcted")
	}

	// https://github.com/google/re2/wiki/Syntax
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f, err := buildReplacementFunc(re, d, config, linkTmpl)
	if err != nil {
		return "", err
	}
	final := re.ReplaceAllStringFunc(linkTmpl, f)

	if re.Match([]byte(final)) {
		return replaceVarsRecursive(d, config, final, depth+1)
	}

	return final, nil
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
func buildReplacementFunc(re *regexp.Regexp, d TerraformResourceData, config *Config, linkTmpl string) (func(string) string, error) {
	var project, projectID, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = getProject(d, config)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{project_id_or_project}}") {
		v, ok := d.GetOkExists("project_id")
		if ok {
			projectID, _ = v.(string)
		}
		if projectID == "" {
			project, err = getProject(d, config)
		}
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = getRegion(d, config)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = getZone(d, config)
		if err != nil {
			return nil, err
		}
	}

	f := func(s string) string {

		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "project_id_or_project" {
			if projectID != "" {
				return projectID
			}
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}
		if string(m[0]) == "%" {
			v, ok := d.GetOkExists(m[1:])
			if ok {
				return url.PathEscape(fmt.Sprintf("%v", v))
			}
		} else {
			v, ok := d.GetOkExists(m)
			if ok {
				return fmt.Sprintf("%v", v)
			}
		}

		// terraform-google-conversion doesn't provide a provider config in tests.
		if config != nil {
			// Attempt to draw values from the provider config if it's present.
			if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
				return f.String()
			}
		}
		return ""
	}

	return f, nil
}
