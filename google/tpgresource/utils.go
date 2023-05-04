package tpgresource

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

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

func PaginatedListRequest(project, baseUrl, userAgent string, config *transport_tpg.Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	res, err := transport_tpg.SendRequest(config, "GET", project, baseUrl, userAgent, nil)
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
		res, err = transport_tpg.SendRequest(config, "GET", project, url, userAgent, nil)
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
