// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

// getRegionFromZone returns the region from a zone for Google cloud.
func getRegionFromZone(zone string) string {
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
func getRegion(d TerraformResourceData, config *Config) (string, error) {
	return getRegionFromSchema("region", "zone", d, config)
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d TerraformResourceData, config *Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

// getBillingProject reads the "billing_project" field from the given resource data and falls
// back to the provider's value if not given. If no value is found, an error is returned.
func getBillingProject(d TerraformResourceData, config *Config) (string, error) {
	return getBillingProjectFromSchema("billing_project", d, config)
}

// getProjectFromDiff reads the "project" field from the given diff and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFromDiff(d *schema.ResourceDiff, config *Config) (string, error) {
	res, ok := d.GetOk("project")
	if ok {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%s: required field is not set", "project")
}

func getRouterLockName(region string, router string) string {
	return fmt.Sprintf("router/%s/%s", region, router)
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if isGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func isGoogleApiErrorWithCode(err error, errCode int) bool {
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

func isFailedPreconditionError(err error) bool {
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

func isConflictError(err error) bool {
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
func isNotFoundGrpcError(err error) bool {
	if errorStatus, ok := status.FromError(err); ok && errorStatus.Code() == codes.NotFound {
		return true
	}
	return false
}

// expandLabels pulls the value of "labels" out of a TerraformResourceData as a map[string]string.
func expandLabels(d TerraformResourceData) map[string]string {
	return expandStringMap(d, "labels")
}

// expandEnvironmentVariables pulls the value of "environment_variables" out of a schema.ResourceData as a map[string]string.
func expandEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "environment_variables")
}

// expandBuildEnvironmentVariables pulls the value of "build_environment_variables" out of a schema.ResourceData as a map[string]string.
func expandBuildEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "build_environment_variables")
}

// expandStringMap pulls the value of key out of a TerraformResourceData as a map[string]string.
func expandStringMap(d TerraformResourceData, key string) map[string]string {
	v, ok := d.GetOk(key)

	if !ok {
		return map[string]string{}
	}

	return convertStringMap(v.(map[string]interface{}))
}

func convertStringMap(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, val := range v {
		m[k] = val.(string)
	}
	return m
}

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func mapStringArr(original []string, f func(string) string) []string {
	var arr []string
	for _, v := range original {
		arr = append(arr, f(v))
	}
	return arr
}

func convertStringArrToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
}

func golangSetFromStringSlice(strings []string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, v := range strings {
		set[v] = struct{}{}
	}

	return set
}

func stringSliceFromGolangSet(sset map[string]struct{}) []string {
	ls := make([]string, 0, len(sset))
	for s := range sset {
		ls = append(ls, s)
	}
	sort.Strings(ls)

	return ls
}

func reverseStringMap(m map[string]string) map[string]string {
	o := map[string]string{}
	for k, v := range m {
		o[v] = k
	}
	return o
}

func mergeStringMaps(a, b map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func mergeResourceMaps(ms ...map[string]*schema.Resource) (map[string]*schema.Resource, error) {
	merged := make(map[string]*schema.Resource)
	duplicates := []string{}

	for _, m := range ms {
		for k, v := range m {
			if _, ok := merged[k]; ok {
				duplicates = append(duplicates, k)
			}

			merged[k] = v
		}
	}

	var err error
	if len(duplicates) > 0 {
		err = fmt.Errorf("saw duplicates in mergeResourceMaps: %v", duplicates)
	}

	return merged, err
}

func stringToFixed64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	if len(m) == 0 || m[0] == nil {
		return map[string]interface{}{}
	}

	return m[0].(map[string]interface{})
}

func lockedCall(lockKey string, f func() error) error {
	mutexKV.Lock(lockKey)
	defer mutexKV.Unlock(lockKey)

	return f()
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

// serviceAccountFQN will attempt to generate the fully qualified name in the format of:
// "projects/(-|<project>)/serviceAccounts/<service_account_id>@<project>.iam.gserviceaccount.com"
// A project is required if we are trying to build the FQN from a service account id and
// and error will be returned in this case if no project is set in the resource or the
// provider-level config
func serviceAccountFQN(serviceAccount string, d TerraformResourceData, config *Config) (string, error) {
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
	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", serviceAccount, project), nil
}

func paginatedListRequest(project, baseUrl, userAgent string, config *Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	res, err := sendRequest(config, "GET", project, baseUrl, userAgent, nil)
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
		res, err = sendRequest(config, "GET", project, url, userAgent, nil)
		if err != nil {
			return nil, err
		}
		ls = append(ls, flattener(res))
		pageToken, ok = res["pageToken"]
	}

	return ls, nil
}

func getInterconnectAttachmentLink(config *Config, project, region, ic, userAgent string) (string, error) {
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
func calcAddRemove(from []string, to []string) (add, remove []string) {
	add = make([]string, 0)
	remove = make([]string, 0)
	for _, u := range to {
		found := false
		for _, v := range from {
			if compareSelfLinkOrResourceName("", v, u, nil) {
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
			if compareSelfLinkOrResourceName("", u, v, nil) {
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

func stringInSlice(arr []string, str string) bool {
	for _, i := range arr {
		if i == str {
			return true
		}
	}

	return false
}

func migrateStateNoop(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return is, nil
}

func expandString(v interface{}, d TerraformResourceData, config *Config) (string, error) {
	return v.(string), nil
}

func changeFieldSchemaToForceNew(sch *schema.Schema) {
	sch.ForceNew = true
	switch sch.Type {
	case schema.TypeList:
	case schema.TypeSet:
		if nestedR, ok := sch.Elem.(*schema.Resource); ok {
			for _, nestedSch := range nestedR.Schema {
				changeFieldSchemaToForceNew(nestedSch)
			}
		}
	}
}

func generateUserAgentString(d TerraformResourceData, currentUserAgent string) (string, error) {
	var m providerMeta

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

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func GetCurrentUserEmail(config *Config, userAgent string) (string, error) {
	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, err := sendRequest(config, "GET", "", "https://openidconnect.googleapis.com/v1/userinfo", userAgent, nil)
	if err != nil {
		return "", fmt.Errorf("error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}
	if res["email"] == nil {
		return "", fmt.Errorf("error retrieving email from userinfo. email was nil in the response.")
	}
	return res["email"].(string), nil
}

func checkStringMap(v interface{}) map[string]string {
	m, ok := v.(map[string]string)
	if ok {
		return m
	}
	return convertStringMap(v.(map[string]interface{}))
}

// return a fake 404 so requests get retried or nested objects are considered deleted
func fake404(reasonResourceType, resourceName string) *googleapi.Error {
	return &googleapi.Error{
		Code:    404,
		Message: fmt.Sprintf("%v object %v not found", reasonResourceType, resourceName),
	}
}

// validate name of the gcs bucket. Guidelines are located at https://cloud.google.com/storage/docs/naming-buckets
// this does not attempt to check for IP addresses or close misspellings of "google"
func checkGCSName(name string) error {
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

// checkGoogleIamPolicy makes assertions about the contents of a google_iam_policy data source's policy_data attribute
func checkGoogleIamPolicy(value string) error {
	if strings.Contains(value, "\"description\":\"\"") {
		return fmt.Errorf("found an empty description field (should be omitted) in google_iam_policy data source: %s", value)
	}
	return nil
}

// Retries an operation while the canonical error code is FAILED_PRECONDTION
// which indicates there is an incompatible operation already running on the
// cluster. This error can be safely retried until the incompatible operation
// completes, and the newly requested operation can begin.
func retryWhileIncompatibleOperation(timeout time.Duration, lockKey string, f func() error) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		if err := lockedCall(lockKey, f); err != nil {
			if isFailedPreconditionError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}
