// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/googleapi"
)

type TerraformResourceData interface {
	HasChange(string) bool
	GetOkExists(string) (interface{}, bool)
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
	Set(string, interface{}) error
	SetId(string)
	Id() string
}

type TerraformResourceDiff interface {
	GetChange(string) (interface{}, interface{})
	Get(string) interface{}
	Clear(string) error
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
	if e, ok := err.(*googleapi.Error); ok && e.Code == 409 {
		return true
	} else if !ok && errwrap.ContainsType(err, &googleapi.Error{}) {
		e := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
		if e.Code == 409 {
			return true
		}
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

func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	if len(m) == 0 {
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

func paginatedListRequest(project, baseUrl string, config *Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	res, err := sendRequest(config, "GET", project, baseUrl, nil)
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
		res, err = sendRequest(config, "GET", project, url, nil)
		if err != nil {
			return nil, err
		}
		ls = append(ls, flattener(res))
		pageToken, ok = res["pageToken"]
	}

	return ls, nil
}

func getInterconnectAttachmentLink(config *Config, project, region, ic string) (string, error) {
	if !strings.Contains(ic, "/") {
		icData, err := config.clientCompute.InterconnectAttachments.Get(
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
