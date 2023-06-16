// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Contains functions that don't really belong anywhere else.

package google

import (
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	fwDiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/googleapi"
)

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
//
// Deprecated: For backward compatibility getRegion is still working,
// but all new code should use GetRegion in the tpgresource package instead.
func getRegion(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetRegion(d, config)
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
//
// Deprecated: For backward compatibility getProject is still working,
// but all new code should use GetProject in the tpgresource package instead.
func getProject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetProject(d, config)
}

// getBillingProject reads the "billing_project" field from the given resource data and falls
// back to the provider's value if not given. If no value is found, an error is returned.
//
// Deprecated: For backward compatibility getBillingProject is still working,
// but all new code should use GetBillingProject in the tpgresource package instead.
func getBillingProject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetBillingProject(d, config)
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
//
// Deprecated: For backward compatibility expandLabels is still working,
// but all new code should use ExpandLabels in the tpgresource package instead.
func expandLabels(d tpgresource.TerraformResourceData) map[string]string {
	return tpgresource.ExpandLabels(d)
}

// expandEnvironmentVariables pulls the value of "environment_variables" out of a schema.ResourceData as a map[string]string.
//
// Deprecated: For backward compatibility expandEnvironmentVariables is still working,
// but all new code should use ExpandEnvironmentVariables in the tpgresource package instead.
func expandEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return tpgresource.ExpandEnvironmentVariables(d)
}

// expandBuildEnvironmentVariables pulls the value of "build_environment_variables" out of a schema.ResourceData as a map[string]string.
//
// Deprecated: For backward compatibility expandBuildEnvironmentVariables is still working,
// but all new code should use ExpandBuildEnvironmentVariables in the tpgresource package instead.
func expandBuildEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return tpgresource.ExpandBuildEnvironmentVariables(d)
}

// expandStringMap pulls the value of key out of a TerraformResourceData as a map[string]string.
//
// Deprecated: For backward compatibility expandStringMap is still working,
// but all new code should use ExpandStringMap in the tpgresource package instead.
func expandStringMap(d tpgresource.TerraformResourceData, key string) map[string]string {
	return tpgresource.ExpandStringMap(d, key)
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

// Deprecated: For backward compatibility reverseStringMap is still working,
// but all new code should use ReverseStringMap in the tpgresource package instead.
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

// Deprecated: For backward compatibility StringToFixed64 is still working,
// but all new code should use StringToFixed64 in the tpgresource package instead.
func StringToFixed64(v string) (int64, error) {
	return tpgresource.StringToFixed64(v)
}

// Deprecated: For backward compatibility extractFirstMapConfig is still working,
// but all new code should use ExtractFirstMapConfig in the tpgresource package instead.
func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	return tpgresource.ExtractFirstMapConfig(m)
}

// Deprecated: For backward compatibility lockedCall is still working,
// but all new code should use LockedCall in the tpgresource package instead.
func lockedCall(lockKey string, f func() error) error {
	return transport_tpg.LockedCall(lockKey, f)
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
//
// Deprecated: For backward compatibility serviceAccountFQN is still working,
// but all new code should use ServiceAccountFQN in the tpgresource package instead.
func serviceAccountFQN(serviceAccount string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.ServiceAccountFQN(serviceAccount, d, config)
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

// Deprecated: For backward compatibility expandString is still working,
// but all new code should use ExpandString in the tpgresource package instead.
func expandString(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.ExpandString(v, d, config)
}

// Deprecated: For backward compatibility changeFieldSchemaToForceNew is still working,
// but all new code should use ChangeFieldSchemaToForceNew in the tpgresource package instead.
func changeFieldSchemaToForceNew(sch *schema.Schema) {
	tpgresource.ChangeFieldSchemaToForceNew(sch)
}

// Deprecated: For backward compatibility generateUserAgentString is still working,
// but all new code should use GenerateUserAgentString in the tpgresource package instead.
func generateUserAgentString(d tpgresource.TerraformResourceData, currentUserAgent string) (string, error) {
	return tpgresource.GenerateUserAgentString(d, currentUserAgent)
}

// Deprecated: For backward compatibility snakeToPascalCase is still working,
// but all new code should use SnakeToPascalCase in the tpgresource package instead.
func snakeToPascalCase(s string) string {
	return tpgresource.SnakeToPascalCase(s)
}

// Deprecated: For backward compatibility checkStringMap is still working,
// but all new code should use CheckStringMap in the tpgresource package instead.
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
	return tpgresource.RetryWhileIncompatibleOperation(timeout, lockKey, f)
}

// Deprecated: For backward compatibility frameworkDiagsToSdkDiags is still working,
// but all new code should use FrameworkDiagsToSdkDiags in the tpgresource package instead.
func frameworkDiagsToSdkDiags(fwD fwDiags.Diagnostics) *diag.Diagnostics {
	return tpgresource.FrameworkDiagsToSdkDiags(fwD)
}

// Deprecated: For backward compatibility isEmptyValue is still working,
// but all new code should use IsEmptyValue in the tpgresource package instead.
func isEmptyValue(v reflect.Value) bool {
	return tpgresource.IsEmptyValue(v)
}

// Deprecated: For backward compatibility replaceVars is still working,
// but all new code should use ReplaceVars in the tpgresource package instead.
func ReplaceVars(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return tpgresource.ReplaceVars(d, config, linkTmpl)
}

// relaceVarsForId shortens variables by running them through GetResourceNameFromSelfLink
// this allows us to use long forms of variables from configs without needing
// custom id formats. For instance:
// accessPolicies/{{access_policy}}/accessLevels/{{access_level}}
// with values:
// access_policy: accessPolicies/foo
// access_level: accessPolicies/foo/accessLevels/bar
// becomes accessPolicies/foo/accessLevels/bar
//
// Deprecated: For backward compatibility replaceVarsForId is still working,
// but all new code should use ReplaceVarsForId in the tpgresource package instead.
func replaceVarsForId(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	return tpgresource.ReplaceVarsForId(d, config, linkTmpl)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
//
// Deprecated: For backward compatibility replaceVarsRecursive is still working,
// but all new code should use ReplaceVarsRecursive in the tpgresource package instead.
func replaceVarsRecursive(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool, depth int) (string, error) {
	return tpgresource.ReplaceVarsRecursive(d, config, linkTmpl, shorten, depth)
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
//
// Deprecated: For backward compatibility buildReplacementFunc is still working,
// but all new code should use BuildReplacementFunc in the tpgresource package instead.
func buildReplacementFunc(re *regexp.Regexp, d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string, shorten bool) (func(string) string, error) {
	return tpgresource.BuildReplacementFunc(re, d, config, linkTmpl, shorten)
}
