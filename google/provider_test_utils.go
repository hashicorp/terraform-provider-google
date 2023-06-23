// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestEnvVar = acctest.TestEnvVar

var TestAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

// providerConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func providerConfigEnvNames() []string {
	return acctest.ProviderConfigEnvNames()
}

var CredsEnvVars = acctest.CredsEnvVars

var projectNumberEnvVars = acctest.ProjectNumberEnvVars

var ProjectEnvVars = acctest.ProjectEnvVars

var firestoreProjectEnvVars = acctest.FirestoreProjectEnvVars

var regionEnvVars = acctest.RegionEnvVars

var zoneEnvVars = acctest.ZoneEnvVars

var orgEnvVars = acctest.OrgEnvVars

// This value is the Customer ID of the GOOGLE_ORG_DOMAIN workspace.
// See https://admin.google.com/ac/accountsettings when logged into an org admin for the value.
var custIdEnvVars = acctest.CustIdEnvVars

// This value is the username of an identity account within the GOOGLE_ORG_DOMAIN workspace.
// For example in the org example.com with a user "foo@example.com", this would be set to "foo".
// See https://admin.google.com/ac/users when logged into an org admin for a list.
var identityUserEnvVars = acctest.IdentityUserEnvVars

var orgEnvDomainVars = acctest.OrgEnvDomainVars

var serviceAccountEnvVars = acctest.ServiceAccountEnvVars

var orgTargetEnvVars = acctest.OrgTargetEnvVars

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
var billingAccountEnvVars = acctest.BillingAccountEnvVars

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
var masterBillingAccountEnvVars = acctest.MasterBillingAccountEnvVars

// This value is the description used for test PublicAdvertisedPrefix setup to avoid required DNS
// setup. This is only used during integration tests and would be invalid to surface to users
var papDescriptionEnvVars = acctest.PapDescriptionEnvVars

func init() {
	configs = make(map[string]*transport_tpg.Config)
	fwProviders = make(map[string]*frameworkTestProvider)
	sources = make(map[string]VcrSource)
	testAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

func GoogleProviderConfig(t *testing.T) *transport_tpg.Config {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		return config
	}

	sdkProvider := provider.Provider()
	rc := terraform.ResourceConfig{}
	sdkProvider.Configure(context.Background(), &rc)
	return sdkProvider.Meta().(*transport_tpg.Config)
}

func AccTestPreCheck(t *testing.T) {
	acctest.AccTestPreCheck(t)
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	return acctest.GetTestRegion(is, config)
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	return acctest.GetTestProject(is, config)
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectNumberFromEnv() string {
	return acctest.GetTestProjectNumberFromEnv()
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return acctest.GetTestProjectFromEnv()
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	return acctest.GetTestCredsFromEnv()
}

// AccTestPreCheck ensures at least one of the region env variables is set.
func GetTestRegionFromEnv() string {
	return acctest.GetTestRegionFromEnv()
}

func GetTestZoneFromEnv() string {
	return acctest.GetTestZoneFromEnv()
}

func GetTestCustIdFromEnv(t *testing.T) string {
	return acctest.GetTestCustIdFromEnv(t)
}

func GetTestIdentityUserFromEnv(t *testing.T) string {
	return acctest.GetTestIdentityUserFromEnv(t)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
func GetTestFirestoreProjectFromEnv(t *testing.T) string {
	return acctest.GetTestFirestoreProjectFromEnv(t)
}

// Returns the raw organization id like 1234567890, skipping the test if one is
// not found.
func GetTestOrgFromEnv(t *testing.T) string {
	return acctest.GetTestOrgFromEnv(t)
}

// Alternative to GetTestOrgFromEnv that doesn't need *testing.T
// If using this, you need to process unset values at the call site
func UnsafeGetTestOrgFromEnv() string {
	return acctest.UnsafeGetTestOrgFromEnv()
}

func GetTestOrgDomainFromEnv(t *testing.T) string {
	return acctest.GetTestOrgDomainFromEnv(t)
}

func GetTestOrgTargetFromEnv(t *testing.T) string {
	return acctest.GetTestOrgTargetFromEnv(t)
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
func GetTestBillingAccountFromEnv(t *testing.T) string {
	return acctest.GetTestBillingAccountFromEnv(t)
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
func GetTestMasterBillingAccountFromEnv(t *testing.T) string {
	return acctest.GetTestMasterBillingAccountFromEnv(t)
}

func GetTestServiceAccountFromEnv(t *testing.T) string {
	return acctest.GetTestServiceAccountFromEnv(t)
}

func GetTestPublicAdvertisedPrefixDescriptionFromEnv(t *testing.T) string {
	return acctest.GetTestPublicAdvertisedPrefixDescriptionFromEnv(t)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	acctest.SkipIfVcr(t)
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return acctest.SleepInSecondsForTest(t)
}

func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	acctest.SkipIfEnvNotSet(t, envs...)
}
