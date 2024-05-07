// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package envvar

import (
	"fmt"
	"log"
	"os"
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const TestEnvVar = "TF_ACC"

var CredsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

// CredsEnvVarsExcludingAdcs returns the contents of CredsEnvVars excluding GOOGLE_APPLICATION_CREDENTIALS
func CredsEnvVarsExcludingAdcs() []string {
	envs := CredsEnvVars
	var filtered []string
	for _, e := range envs {
		if e != "GOOGLE_APPLICATION_CREDENTIALS" {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

var ProjectNumberEnvVars = []string{
	"GOOGLE_PROJECT_NUMBER",
}

var ProjectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
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

var UniverseDomainEnvVars = []string{
	"GOOGLE_UNIVERSE_DOMAIN",
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

// Returns googleapis.com if there's no universe set.
func GetTestUniverseDomainFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, IdentityUserEnvVars...)
	return transport_tpg.MultiEnvSearch(UniverseDomainEnvVars)
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

func ServiceAccountCanonicalEmail(account string) string {
	return fmt.Sprintf("%s@%s.iam.gserviceaccount.com", account, GetTestProjectFromEnv())
}
