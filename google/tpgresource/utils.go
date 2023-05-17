package tpgresource

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/errwrap"
	fwDiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TerraformResourceDataChange interface {
	GetChange(string) (interface{}, interface{})
}

type TerraformResourceData interface {
	HasChange(string) bool
	GetOkExists(string) (interface{}, bool)
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
	Set(string, interface{}) error
	SetId(string)
	Id() string
	GetProviderMeta(interface{}) error
	Timeout(key string) time.Duration
}

type TerraformResourceDiff interface {
	HasChange(string) bool
	GetChange(string) (interface{}, interface{})
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Clear(string) error
	ForceNew(string) error
}

// Contains functions that don't really belong anywhere else.

// GetRegionFromZone returns the region from a zone for Google cloud.
// This is by removing the last two chars from the zone name to leave the region
// If there aren't enough characters in the input string, an empty string is returned
// e.g. southamerica-west1-a => southamerica-west1
func GetRegionFromZone(zone string) string {
	if zone != "" && len(zone) > 2 {
		region := zone[:len(zone)-2]
		return region
	}
	return ""
}

// Infers the region based on the following (in order of priority):
// - `region` field in resource schema
// - region extracted from the `zone` field in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func GetRegion(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return GetRegionFromSchema("region", "zone", d, config)
}

// GetProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func GetProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return GetProjectFromSchema("project", d, config)
}

// GetBillingProject reads the "billing_project" field from the given resource data and falls
// back to the provider's value if not given. If no value is found, an error is returned.
func GetBillingProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return GetBillingProjectFromSchema("billing_project", d, config)
}

// GetProjectFromDiff reads the "project" field from the given diff and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func GetProjectFromDiff(d *schema.ResourceDiff, config *transport_tpg.Config) (string, error) {
	res, ok := d.GetOk("project")
	if ok {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%s: required field is not set", "project")
}

func GetRouterLockName(region string, router string) string {
	return fmt.Sprintf("router/%s/%s", region, router)
}

func IsFailedPreconditionError(err error) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	if !ok {
		return false
	}
	if gerr == nil {
		return false
	}
	if gerr.Code != 400 {
		return false
	}
	for _, e := range gerr.Errors {
		if e.Reason == "failedPrecondition" {
			return true
		}
	}
	return false
}

func IsConflictError(err error) bool {
	if e, ok := err.(*googleapi.Error); ok && (e.Code == 409 || e.Code == 412) {
		return true
	} else if !ok && errwrap.ContainsType(err, &googleapi.Error{}) {
		e := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
		if e.Code == 409 || e.Code == 412 {
			return true
		}
	}
	return false
}

// gRPC does not return errors of type *googleapi.Error. Instead the errors returned are *status.Error.
// See the types of codes returned here (https://pkg.go.dev/google.golang.org/grpc/codes#Code).
func IsNotFoundGrpcError(err error) bool {
	if errorStatus, ok := status.FromError(err); ok && errorStatus.Code() == codes.NotFound {
		return true
	}
	return false
}

// ExpandLabels pulls the value of "labels" out of a TerraformResourceData as a map[string]string.
func ExpandLabels(d TerraformResourceData) map[string]string {
	return ExpandStringMap(d, "labels")
}

// ExpandEnvironmentVariables pulls the value of "environment_variables" out of a schema.ResourceData as a map[string]string.
func ExpandEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return ExpandStringMap(d, "environment_variables")
}

// ExpandBuildEnvironmentVariables pulls the value of "build_environment_variables" out of a schema.ResourceData as a map[string]string.
func ExpandBuildEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return ExpandStringMap(d, "build_environment_variables")
}

// ExpandStringMap pulls the value of key out of a TerraformResourceData as a map[string]string.
func ExpandStringMap(d TerraformResourceData, key string) map[string]string {
	v, ok := d.GetOk(key)

	if !ok {
		return map[string]string{}
	}

	return ConvertStringMap(v.(map[string]interface{}))
}

func ConvertStringMap(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, val := range v {
		m[k] = val.(string)
	}
	return m
}

func ConvertStringArr(ifaceArr []interface{}) []string {
	return ConvertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func ConvertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func MapStringArr(original []string, f func(string) string) []string {
	var arr []string
	for _, v := range original {
		arr = append(arr, f(v))
	}
	return arr
}

func ConvertStringArrToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

func ConvertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
}

func GolangSetFromStringSlice(strings []string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, v := range strings {
		set[v] = struct{}{}
	}

	return set
}

func StringSliceFromGolangSet(sset map[string]struct{}) []string {
	ls := make([]string, 0, len(sset))
	for s := range sset {
		ls = append(ls, s)
	}
	sort.Strings(ls)

	return ls
}

func ReverseStringMap(m map[string]string) map[string]string {
	o := map[string]string{}
	for k, v := range m {
		o[v] = k
	}
	return o
}

func MergeStringMaps(a, b map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func MergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func StringToFixed64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func ExtractFirstMapConfig(m []interface{}) map[string]interface{} {
	if len(m) == 0 || m[0] == nil {
		return map[string]interface{}{}
	}

	return m[0].(map[string]interface{})
}

// This is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
func Nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}", fmt.Sprintf("%v", val), -1)
	}
	return format
}

//	ServiceAccountFQN will attempt to generate the fully qualified name in the format of:
//
// "projects/(-|<project>)/serviceAccounts/<service_account_id>@<project>.iam.gserviceaccount.com"
// A project is required if we are trying to build the FQN from a service account id and
// and error will be returned in this case if no project is set in the resource or the
// provider-level config
func ServiceAccountFQN(serviceAccount string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	// If the service account id is already the fully qualified name
	if strings.HasPrefix(serviceAccount, "projects/") {
		return serviceAccount, nil
	}

	// If the service account id is an email
	if strings.Contains(serviceAccount, "@") {
		return "projects/-/serviceAccounts/" + serviceAccount, nil
	}

	// Get the project from the resource or fallback to the project
	// in the provider configuration
	project, err := GetProject(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", serviceAccount, project), nil
}

func PaginatedListRequest(project, baseUrl, userAgent string, config *transport_tpg.Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    baseUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	ls := flattener(res)
	pageToken, ok := res["pageToken"]
	for ok {
		if pageToken.(string) == "" {
			break
		}
		url := fmt.Sprintf("%s?pageToken=%s", baseUrl, pageToken.(string))
		res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return nil, err
		}
		ls = append(ls, flattener(res))
		pageToken, ok = res["pageToken"]
	}

	return ls, nil
}

func GetInterconnectAttachmentLink(config *transport_tpg.Config, project, region, ic, userAgent string) (string, error) {
	if !strings.Contains(ic, "/") {
		icData, err := config.NewComputeClient(userAgent).InterconnectAttachments.Get(
			project, region, ic).Do()
		if err != nil {
			return "", fmt.Errorf("Error reading interconnect attachment: %s", err)
		}
		ic = icData.SelfLink
	}

	return ic, nil
}

// Given two sets of references (with "from" values in self link form),
// determine which need to be added or removed // during an update using
// addX/removeX APIs.
func CalcAddRemove(from []string, to []string) (add, remove []string) {
	add = make([]string, 0)
	remove = make([]string, 0)
	for _, u := range to {
		found := false
		for _, v := range from {
			if CompareSelfLinkOrResourceName("", v, u, nil) {
				found = true
				break
			}
		}
		if !found {
			add = append(add, u)
		}
	}
	for _, u := range from {
		found := false
		for _, v := range to {
			if CompareSelfLinkOrResourceName("", u, v, nil) {
				found = true
				break
			}
		}
		if !found {
			remove = append(remove, u)
		}
	}
	return add, remove
}

func StringInSlice(arr []string, str string) bool {
	for _, i := range arr {
		if i == str {
			return true
		}
	}

	return false
}

func MigrateStateNoop(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return is, nil
}

func ExpandString(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return v.(string), nil
}

func ChangeFieldSchemaToForceNew(sch *schema.Schema) {
	sch.ForceNew = true
	switch sch.Type {
	case schema.TypeList:
	case schema.TypeSet:
		if nestedR, ok := sch.Elem.(*schema.Resource); ok {
			for _, nestedSch := range nestedR.Schema {
				ChangeFieldSchemaToForceNew(nestedSch)
			}
		}
	}
}

func GenerateUserAgentString(d TerraformResourceData, currentUserAgent string) (string, error) {
	var m transport_tpg.ProviderMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return currentUserAgent, err
	}

	if m.ModuleName != "" {
		return strings.Join([]string{currentUserAgent, m.ModuleName}, " "), nil
	}

	return currentUserAgent, nil
}

func SnakeToPascalCase(s string) string {
	split := strings.Split(s, "_")
	for i := range split {
		split[i] = strings.Title(split[i])
	}
	return strings.Join(split, "")
}

func CheckStringMap(v interface{}) map[string]string {
	m, ok := v.(map[string]string)
	if ok {
		return m
	}
	return ConvertStringMap(v.(map[string]interface{}))
}

// return a fake 404 so requests get retried or nested objects are considered deleted
func Fake404(reasonResourceType, resourceName string) *googleapi.Error {
	return &googleapi.Error{
		Code:    404,
		Message: fmt.Sprintf("%v object %v not found", reasonResourceType, resourceName),
	}
}

// validate name of the gcs bucket. Guidelines are located at https://cloud.google.com/storage/docs/naming-buckets
// this does not attempt to check for IP addresses or close misspellings of "google"
func CheckGCSName(name string) error {
	if strings.HasPrefix(name, "goog") {
		return fmt.Errorf("error: bucket name %s cannot start with %q", name, "goog")
	}

	if strings.Contains(name, "google") {
		return fmt.Errorf("error: bucket name %s cannot contain %q", name, "google")
	}

	valid, _ := regexp.MatchString("^[a-z0-9][a-z0-9_.-]{1,220}[a-z0-9]$", name)
	if !valid {
		return fmt.Errorf("error: bucket name validation failed %v. See https://cloud.google.com/storage/docs/naming-buckets", name)
	}

	for _, str := range strings.Split(name, ".") {
		valid, _ := regexp.MatchString("^[a-z0-9_-]{1,63}$", str)
		if !valid {
			return fmt.Errorf("error: bucket name validation failed %v", str)
		}
	}
	return nil
}

// CheckGoogleIamPolicy makes assertions about the contents of a google_iam_policy data source's policy_data attribute
func CheckGoogleIamPolicy(value string) error {
	if strings.Contains(value, "\"description\":\"\"") {
		return fmt.Errorf("found an empty description field (should be omitted) in google_iam_policy data source: %s", value)
	}
	return nil
}

func FrameworkDiagsToSdkDiags(fwD fwDiags.Diagnostics) *diag.Diagnostics {
	var diags diag.Diagnostics
	for _, e := range fwD.Errors() {
		diags = append(diags, diag.Diagnostic{
			Detail:   e.Detail(),
			Severity: diag.Error,
			Summary:  e.Summary(),
		})
	}
	for _, w := range fwD.Warnings() {
		diags = append(diags, diag.Diagnostic{
			Detail:   w.Detail(),
			Severity: diag.Warning,
			Summary:  w.Summary(),
		})
	}

	return &diags
}

func IsEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

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

func ReplaceVars(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return ReplaceVarsRecursive(d, config, linkTmpl, false, 0)
}

// relaceVarsForId shortens variables by running them through GetResourceNameFromSelfLink
// this allows us to use long forms of variables from configs without needing
// custom id formats. For instance:
// accessPolicies/{{access_policy}}/accessLevels/{{access_level}}
// with values:
// access_policy: accessPolicies/foo
// access_level: accessPolicies/foo/accessLevels/bar
// becomes accessPolicies/foo/accessLevels/bar
func ReplaceVarsForId(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return ReplaceVarsRecursive(d, config, linkTmpl, true, 0)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func ReplaceVarsRecursive(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool, depth int) (string, error) {
	if depth > 10 {
		return "", errors.New("Recursive substitution detcted")
	}

	// https://github.com/google/re2/wiki/Syntax
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f, err := BuildReplacementFunc(re, d, config, linkTmpl, shorten)
	if err != nil {
		return "", err
	}
	final := re.ReplaceAllStringFunc(linkTmpl, f)

	if re.Match([]byte(final)) {
		return ReplaceVarsRecursive(d, config, final, shorten, depth+1)
	}

	return final, nil
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
func BuildReplacementFunc(re *regexp.Regexp, d TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool) (func(string) string, error) {
	var project, projectID, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = GetProject(d, config)
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
			project, err = GetProject(d, config)
		}
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = GetRegion(d, config)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = GetZone(d, config)
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
