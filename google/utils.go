// Contains functions that don't really belong anywhere else.

package google

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	fwDiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/googleapi"
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
// This is by removing the last two chars from the zone name to leave the region
// If there aren't enough characters in the input string, an empty string is returned
// e.g. southamerica-west1-a => southamerica-west1
//
// Deprecated: For backward compatibility getRegionFromZone is still working,
// but all new code should use GetRegionFromZone in the tpgresource package instead.
func getRegionFromZone(zone string) string {
	return tpgresource.GetRegionFromZone(zone)
}

// Infers the region based on the following (in order of priority):
// - `region` field in resource schema
// - region extracted from the `zone` field in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func getRegion(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getRegionFromSchema("region", "zone", d, config)
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

// getBillingProject reads the "billing_project" field from the given resource data and falls
// back to the provider's value if not given. If no value is found, an error is returned.
func getBillingProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getBillingProjectFromSchema("billing_project", d, config)
}

// getProjectFromDiff reads the "project" field from the given diff and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
//
// Deprecated: For backward compatibility getProjectFromDiff is still working,
// but all new code should use GetProjectFromDiff in the tpgresource package instead.
func getProjectFromDiff(d *schema.ResourceDiff, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetProjectFromDiff(d, config)
}

// Deprecated: For backward compatibility getRouterLockName is still working,
// but all new code should use GetRouterLockName in the tpgresource package instead.
func getRouterLockName(region string, router string) string {
	return tpgresource.GetRouterLockName(region, router)
}

// Deprecated: For backward compatibility isFailedPreconditionError is still working,
// but all new code should use IsFailedPreconditionError in the tpgresource package instead.
func isFailedPreconditionError(err error) bool {
	return tpgresource.IsFailedPreconditionError(err)
}

// Deprecated: For backward compatibility isConflictError is still working,
// but all new code should use IsConflictError in the tpgresource package instead.
func isConflictError(err error) bool {
	return tpgresource.IsConflictError(err)
}

// gRPC does not return errors of type *googleapi.Error. Instead the errors returned are *status.Error.
// See the types of codes returned here (https://pkg.go.dev/google.golang.org/grpc/codes#Code).
//
// Deprecated: For backward compatibility isNotFoundGrpcError is still working,
// but all new code should use IsNotFoundGrpcError in the tpgresource package instead.
func isNotFoundGrpcError(err error) bool {
	return tpgresource.IsNotFoundGrpcError(err)
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

// Deprecated: For backward compatibility convertStringMap is still working,
// but all new code should use ConvertStringMap in the tpgresource package instead.
func convertStringMap(v map[string]interface{}) map[string]string {
	return tpgresource.ConvertStringMap(v)
}

// Deprecated: For backward compatibility convertStringArr is still working,
// but all new code should use ConvertStringArr in the tpgresource package instead.
func convertStringArr(ifaceArr []interface{}) []string {
	return tpgresource.ConvertStringArr(ifaceArr)
}

// Deprecated: For backward compatibility convertAndMapStringArr is still working,
// but all new code should use ConvertAndMapStringArr in the tpgresource package instead.
func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	return tpgresource.ConvertAndMapStringArr(ifaceArr, f)
}

// Deprecated: For backward compatibility mapStringArr is still working,
// but all new code should use MapStringArr in the tpgresource package instead.
func mapStringArr(original []string, f func(string) string) []string {
	return tpgresource.MapStringArr(original, f)
}

// Deprecated: For backward compatibility convertStringArrToInterface is still working,
// but all new code should use ConvertStringArrToInterface in the tpgresource package instead.
func convertStringArrToInterface(strs []string) []interface{} {
	return tpgresource.ConvertStringArrToInterface(strs)
}

// Deprecated: For backward compatibility convertStringSet is still working,
// but all new code should use ConvertStringSet in the tpgresource package instead.
func convertStringSet(set *schema.Set) []string {
	return tpgresource.ConvertStringSet(set)
}

// Deprecated: For backward compatibility golangSetFromStringSlice is still working,
// but all new code should use GolangSetFromStringSlice in the tpgresource package instead.
func golangSetFromStringSlice(strings []string) map[string]struct{} {
	return tpgresource.GolangSetFromStringSlice(strings)
}

// Deprecated: For backward compatibility stringSliceFromGolangSet is still working,
// but all new code should use StringSliceFromGolangSet in the tpgresource package instead.
func stringSliceFromGolangSet(sset map[string]struct{}) []string {
	return tpgresource.StringSliceFromGolangSet(sset)
}

// Deprecated: For backward compatibility getRegionFromZone is still working,
// but all new code should use GetRegionFromZone in the tpgresource package instead.
func reverseStringMap(m map[string]string) map[string]string {
	return tpgresource.ReverseStringMap(m)
}

// Deprecated: For backward compatibility mergeStringMaps is still working,
// but all new code should use MergeStringMaps in the tpgresource package instead.
func mergeStringMaps(a, b map[string]string) map[string]string {
	return tpgresource.MergeStringMaps(a, b)
}

// Deprecated: For backward compatibility mergeSchemas is still working,
// but all new code should use MergeSchemas in the tpgresource package instead.
func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	return tpgresource.MergeSchemas(a, b)
}

// Deprecated: For backward compatibility getRegionFromZone is still working,
// but all new code should use GetRegionFromZone in the tpgresource package instead.
func StringToFixed64(v string) (int64, error) {
	return tpgresource.StringToFixed64(v)
}

// Deprecated: For backward compatibility extractFirstMapConfig is still working,
// but all new code should use ExtractFirstMapConfig in the tpgresource package instead.
func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	return tpgresource.ExtractFirstMapConfig(m)
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
//
// Deprecated: For backward compatibility Nprintf is still working,
// but all new code should use Nprintf in the tpgresource package instead.
func Nprintf(format string, params map[string]interface{}) string {
	return tpgresource.Nprintf(format, params)
}

// serviceAccountFQN will attempt to generate the fully qualified name in the format of:
// "projects/(-|<project>)/serviceAccounts/<service_account_id>@<project>.iam.gserviceaccount.com"
// A project is required if we are trying to build the FQN from a service account id and
// and error will be returned in this case if no project is set in the resource or the
// provider-level config
func serviceAccountFQN(serviceAccount string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
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

// Deprecated: For backward compatibility paginatedListRequest is still working,
// but all new code should use PaginatedListRequest in the tpgresource package instead.
func paginatedListRequest(project, baseUrl, userAgent string, config *transport_tpg.Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	return tpgresource.PaginatedListRequest(project, baseUrl, userAgent, config, flattener)
}

// Deprecated: For backward compatibility getInterconnectAttachmentLink is still working,
// but all new code should use GetInterconnectAttachmentLink in the tpgresource package instead.
func getInterconnectAttachmentLink(config *transport_tpg.Config, project, region, ic, userAgent string) (string, error) {
	return tpgresource.GetInterconnectAttachmentLink(config, project, region, ic, userAgent)
}

// Given two sets of references (with "from" values in self link form),
// determine which need to be added or removed // during an update using
// addX/removeX APIs.
//
// Deprecated: For backward compatibility calcAddRemove is still working,
// but all new code should use CalcAddRemove in the tpgresource package instead.
func calcAddRemove(from []string, to []string) (add, remove []string) {
	return tpgresource.CalcAddRemove(from, to)
}

// Deprecated: For backward compatibility stringInSlice is still working,
// but all new code should use StringInSlice in the tpgresource package instead.
func stringInSlice(arr []string, str string) bool {
	return tpgresource.StringInSlice(arr, str)
}

// Deprecated: For backward compatibility migrateStateNoop is still working,
// but all new code should use MigrateStateNoop in the tpgresource package instead.
func migrateStateNoop(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return tpgresource.MigrateStateNoop(v, is, meta)
}

func expandString(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return v.(string), nil
}

// Deprecated: For backward compatibility changeFieldSchemaToForceNew is still working,
// but all new code should use ChangeFieldSchemaToForceNew in the tpgresource package instead.
func changeFieldSchemaToForceNew(sch *schema.Schema) {
	tpgresource.ChangeFieldSchemaToForceNew(sch)
}

func generateUserAgentString(d TerraformResourceData, currentUserAgent string) (string, error) {
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
	return tpgresource.SnakeToPascalCase(s)
}

func checkStringMap(v interface{}) map[string]string {
	return tpgresource.CheckStringMap(v)
}

// return a fake 404 so requests get retried or nested objects are considered deleted
//
// Deprecated: For backward compatibility fake404 is still working,
// but all new code should use Fake404 in the tpgresource package instead.
func fake404(reasonResourceType, resourceName string) *googleapi.Error {
	return tpgresource.Fake404(reasonResourceType, resourceName)
}

// validate name of the gcs bucket. Guidelines are located at https://cloud.google.com/storage/docs/naming-buckets
// this does not attempt to check for IP addresses or close misspellings of "google"
//
// Deprecated: For backward compatibility checkGCSName is still working,
// but all new code should use CheckGCSName in the tpgresource package instead.
func checkGCSName(name string) error {
	return tpgresource.CheckGCSName(name)
}

// checkGoogleIamPolicy makes assertions about the contents of a google_iam_policy data source's policy_data attribute
//
// Deprecated: For backward compatibility checkGoogleIamPolicy is still working,
// but all new code should use CheckGoogleIamPolicy in the tpgresource package instead.
func checkGoogleIamPolicy(value string) error {
	return tpgresource.CheckGoogleIamPolicy(value)
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

// Deprecated: For backward compatibility frameworkDiagsToSdkDiags is still working,
// but all new code should use FrameworkDiagsToSdkDiags in the tpgresource package instead.
func frameworkDiagsToSdkDiags(fwD fwDiags.Diagnostics) *diag.Diagnostics {
	return tpgresource.FrameworkDiagsToSdkDiags(fwD)
}

// Deprecated: For backward compatibility isEmptyValue is still working,
// but all new code should use IsEmptyValue in the verify package instead.
//
// Deprecated: For backward compatibility isEmptyValue is still working,
// but all new code should use IsEmptyValue in the tpgresource package instead.
func isEmptyValue(v reflect.Value) bool {
	return tpgresource.IsEmptyValue(v)
}

func ReplaceVars(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
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
func replaceVarsForId(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, true, 0)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func replaceVarsRecursive(d TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool, depth int) (string, error) {
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
func buildReplacementFunc(re *regexp.Regexp, d TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool) (func(string) string, error) {
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
					return tpgresource.GetResourceNameFromSelfLink(fmt.Sprintf("%v", v))
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
