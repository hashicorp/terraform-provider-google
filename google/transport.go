package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

var DefaultRequestTimeout = 5 * time.Minute

func SendRequest(config *Config, method, project, rawurl, userAgent string, body map[string]interface{}, errorRetryPredicates ...RetryErrorPredicateFunc) (map[string]interface{}, error) {
	return SendRequestWithTimeout(config, method, project, rawurl, userAgent, body, DefaultRequestTimeout, errorRetryPredicates...)
}

func SendRequestWithTimeout(config *Config, method, project, rawurl, userAgent string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...RetryErrorPredicateFunc) (map[string]interface{}, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", userAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if config.UserProjectOverride && project != "" {
		// When project is "NO_BILLING_PROJECT_OVERRIDE" in the function GetCurrentUserEmail,
		// set the header X-Goog-User-Project to be empty string.
		if project == "NO_BILLING_PROJECT_OVERRIDE" {
			reqHeaders.Set("X-Goog-User-Project", "")
		} else {
			// Pass the project into this fn instead of parsing it from the URL because
			// both project names and URLs can have colons in them.
			reqHeaders.Set("X-Goog-User-Project", project)
		}
	}

	if timeout == 0 {
		timeout = time.Duration(1) * time.Hour
	}

	var res *http.Response
	err := RetryTimeDuration(
		func() error {
			var buf bytes.Buffer
			if body != nil {
				err := json.NewEncoder(&buf).Encode(body)
				if err != nil {
					return err
				}
			}

			u, err := AddQueryParams(rawurl, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}
			req, err := http.NewRequest(method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = config.Client.Do(req)
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
		return nil, fmt.Errorf("Unable to parse server response. This is most likely a terraform problem, please file a bug at https://github.com/hashicorp/terraform-provider-google/issues.")
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

func AddQueryParams(rawurl string, params map[string]string) (string, error) {
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

func ReplaceVars(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, false, 0)
}

// relaceVarsForId shortens variables by running them through GetResourceNameFromSelfLink
// this allows us to use long forms of variables from configs without needing
// custom id formats. For instance:
// accessPolicies/{{access_policy}}/accessLevels/{{access_level}}
// with values:
// access_policy: accessPolicies/foo
// access_level: accessPolicies/foo/accessLevels/bar
// becomes accessPolicies/foo/accessLevels/bar
func replaceVarsForId(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, true, 0)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func replaceVarsRecursive(d TerraformResourceData, config *Config, linkTmpl string, shorten bool, depth int) (string, error) {
	if depth > 10 {
		return "", errors.New("Recursive substitution detcted")
	}

	// https://github.com/google/re2/wiki/Syntax
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f, err := buildReplacementFunc(re, d, config, linkTmpl, shorten)
	if err != nil {
		return "", err
	}
	final := re.ReplaceAllStringFunc(linkTmpl, f)

	if re.Match([]byte(final)) {
		return replaceVarsRecursive(d, config, final, shorten, depth+1)
	}

	return final, nil
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
func buildReplacementFunc(re *regexp.Regexp, d TerraformResourceData, config *Config, linkTmpl string, shorten bool) (func(string) string, error) {
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
				if shorten {
					return GetResourceNameFromSelfLink(fmt.Sprintf("%v", v))
				} else {
					return fmt.Sprintf("%v", v)
				}
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

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func IsGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}

func isApiNotEnabledError(err error) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	if !ok {
		return false
	}
	if gerr == nil {
		return false
	}
	if gerr.Code != 403 {
		return false
	}
	for _, e := range gerr.Errors {
		if e.Reason == "accessNotConfigured" {
			return true
		}
	}
	return false
}
