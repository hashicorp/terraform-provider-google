// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestEnvVar = "TF_ACC"

// ProviderConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func ProviderConfigEnvNames() []string {

	envs := []string{}

	// Use existing collections of ENV names
	envVarsSets := [][]string{
		CredsEnvVars,   // credentials field
		ProjectEnvVars, // project field
		RegionEnvVars,  //region field
		ZoneEnvVars,    // zone field
	}
	for _, set := range envVarsSets {
		envs = append(envs, set...)
	}

	// Add remaining ENVs
	envs = append(envs, "GOOGLE_OAUTH_ACCESS_TOKEN")          // access_token field
	envs = append(envs, "GOOGLE_BILLING_PROJECT")             // billing_project field
	envs = append(envs, "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") // impersonate_service_account field
	envs = append(envs, "USER_PROJECT_OVERRIDE")              // user_project_override field
	envs = append(envs, "CLOUDSDK_CORE_REQUEST_REASON")       // request_reason field

	return envs
}

var CredsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var ProjectNumberEnvVars = []string{
	"GOOGLE_PROJECT_NUMBER",
}

var ProjectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var FirestoreProjectEnvVars = []string{
	"GOOGLE_FIRESTORE_PROJECT",
}

var RegionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

var ZoneEnvVars = []string{
	"GOOGLE_ZONE",
	"GCLOUD_ZONE",
	"CLOUDSDK_COMPUTE_ZONE",
}

var OrgEnvVars = []string{
	"GOOGLE_ORG",
}

// This value is the Customer ID of the GOOGLE_ORG_DOMAIN workspace.
// See https://admin.google.com/ac/accountsettings when logged into an org admin for the value.
var CustIdEnvVars = []string{
	"GOOGLE_CUST_ID",
}

// This value is the username of an identity account within the GOOGLE_ORG_DOMAIN workspace.
// For example in the org example.com with a user "foo@example.com", this would be set to "foo".
// See https://admin.google.com/ac/users when logged into an org admin for a list.
var IdentityUserEnvVars = []string{
	"GOOGLE_IDENTITY_USER",
}

var OrgEnvDomainVars = []string{
	"GOOGLE_ORG_DOMAIN",
}

var ServiceAccountEnvVars = []string{
	"GOOGLE_SERVICE_ACCOUNT",
}

var OrgTargetEnvVars = []string{
	"GOOGLE_ORG_2",
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
var BillingAccountEnvVars = []string{
	"GOOGLE_BILLING_ACCOUNT",
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
var MasterBillingAccountEnvVars = []string{
	"GOOGLE_MASTER_BILLING_ACCOUNT",
}

// This value is the description used for test PublicAdvertisedPrefix setup to avoid required DNS
// setup. This is only used during integration tests and would be invalid to surface to users
var PapDescriptionEnvVars = []string{
	"GOOGLE_PUBLIC_AVERTISED_PREFIX_DESCRIPTION",
}

func AccTestPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	if v := transport_tpg.MultiEnvSearch(CredsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(CredsEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ProjectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ProjectEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(RegionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(RegionEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ZoneEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ZoneEnvVars, ", "))
	}
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectNumberFromEnv() string {
	return transport_tpg.MultiEnvSearch(ProjectNumberEnvVars)
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return transport_tpg.MultiEnvSearch(ProjectEnvVars)
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	// Return empty string if GOOGLE_USE_DEFAULT_CREDENTIALS is set to true.
	if transport_tpg.MultiEnvSearch(CredsEnvVars) == "true" {
		return ""
	}
	return transport_tpg.MultiEnvSearch(CredsEnvVars)
}

// AccTestPreCheck ensures at least one of the region env variables is set.
func GetTestRegionFromEnv() string {
	return transport_tpg.MultiEnvSearch(RegionEnvVars)
}

func GetTestZoneFromEnv() string {
	return transport_tpg.MultiEnvSearch(ZoneEnvVars)
}

func GetTestCustIdFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, CustIdEnvVars...)
	return transport_tpg.MultiEnvSearch(CustIdEnvVars)
}

func GetTestIdentityUserFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, IdentityUserEnvVars...)
	return transport_tpg.MultiEnvSearch(IdentityUserEnvVars)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
func GetTestFirestoreProjectFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, FirestoreProjectEnvVars...)
	return transport_tpg.MultiEnvSearch(FirestoreProjectEnvVars)
}

// Returns the raw organization id like 1234567890, skipping the test if one is
// not found.
func GetTestOrgFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, OrgEnvVars...)
	return transport_tpg.MultiEnvSearch(OrgEnvVars)
}

// Alternative to GetTestOrgFromEnv that doesn't need *testing.T
// If using this, you need to process unset values at the call site
func UnsafeGetTestOrgFromEnv() string {
	return transport_tpg.MultiEnvSearch(OrgEnvVars)
}

func GetTestOrgDomainFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, OrgEnvDomainVars...)
	return transport_tpg.MultiEnvSearch(OrgEnvDomainVars)
}

func GetTestOrgTargetFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, OrgTargetEnvVars...)
	return transport_tpg.MultiEnvSearch(OrgTargetEnvVars)
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
func GetTestBillingAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, BillingAccountEnvVars...)
	return transport_tpg.MultiEnvSearch(BillingAccountEnvVars)
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
func GetTestMasterBillingAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, MasterBillingAccountEnvVars...)
	return transport_tpg.MultiEnvSearch(MasterBillingAccountEnvVars)
}

func GetTestServiceAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, ServiceAccountEnvVars...)
	return transport_tpg.MultiEnvSearch(ServiceAccountEnvVars)
}

func GetTestPublicAdvertisedPrefixDescriptionFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, PapDescriptionEnvVars...)
	return transport_tpg.MultiEnvSearch(PapDescriptionEnvVars)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for _, k := range envs {
		if os.Getenv(k) == "" {
			log.Printf("[DEBUG] Warning - environment variable %s is not set - skipping test %s", k, t.Name())
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}
