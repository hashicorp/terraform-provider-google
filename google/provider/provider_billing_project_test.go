// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccSdkProvider_billing_project is a series of acc tests asserting how the SDK provider handles billing_project arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_billing_project
func TestAccSdkProvider_billing_project(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                           testAccSdkProvider_billing_project_configPrecedenceOverEnvironmentVariables,
		"when billing_project is unset in the config, environment variables are used in a given order": testAccSdkProvider_billing_project_precedenceOrderEnvironmentVariables, // GOOGLE_BILLING_PROJECT

		// Schema-level validation
		"when billing_project is set to an empty string in the config the value isn't ignored and results in an error": testAccSdkProvider_billing_project_emptyStringValidation,

		// Usage
		// TODO: https://github.com/hashicorp/terraform-provider-google/issues/17882
		"GOOGLE_CLOUD_QUOTA_PROJECT environment variable interferes with the billing_account value used": testAccSdkProvider_billing_project_affectedByClientLibraryEnv,
		// 1) Usage of billing_account alone is insufficient
		// 2) Usage in combination with user_project_override changes the project where quota is used
		"using billing_account alone doesn't impact provisioning, but using together with user_project_override does": testAccSdkProvider_billing_project_useWithAndWithoutUserProjectOverride,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSdkProvider_billing_project_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all possible billing_project env vars set; show they aren't used instead
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	providerBillingProject := "foobar"

	context := map[string]interface{}{
		"billing_project": providerBillingProject,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccSdkProvider_billing_project_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "billing_project", providerBillingProject),
				)},
		},
	})
}

func testAccSdkProvider_billing_project_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for billing_project
		GOOGLE_BILLING_PROJECT

		GOOGLE_CLOUD_QUOTA_PROJECT - NOT used by provider, but is in client libraries we use
	*/

	GOOGLE_BILLING_PROJECT := "GOOGLE_BILLING_PROJECT"
	GOOGLE_CLOUD_QUOTA_PROJECT := "GOOGLE_CLOUD_QUOTA_PROJECT"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_BILLING_PROJECT is used if set
				PreConfig: func() {
					t.Setenv("GOOGLE_BILLING_PROJECT", GOOGLE_BILLING_PROJECT) //used
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", GOOGLE_CLOUD_QUOTA_PROJECT)
				},
				Config: testAccSdkProvider_billing_project_inEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "billing_project", GOOGLE_BILLING_PROJECT),
				),
			},
			{
				// GOOGLE_CLOUD_QUOTA_PROJECT is NOT used here
				PreConfig: func() {
					t.Setenv("GOOGLE_BILLING_PROJECT", "")
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", GOOGLE_CLOUD_QUOTA_PROJECT) // NOT used
				},
				Config: testAccSdkProvider_billing_project_inEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "billing_project", ""),
				),
			},
		},
	})
}

func testAccSdkProvider_billing_project_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all billing_project env vars set
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	context := map[string]interface{}{
		"billing_project": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_billing_project_inProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccSdkProvider_billing_project_useWithAndWithoutUserProjectOverride(t *testing.T) {
	// Test cannot run in VCR mode due to use of aliases
	// See: https://github.com/hashicorp/terraform-provider-google/issues/20019
	acctest.SkipIfVcr(t)

	randomString := acctest.RandString(t, 10)
	contextUserProjectOverrideFalse := map[string]interface{}{
		"org_id":                envvar.GetTestOrgFromEnv(t),
		"billing_account":       envvar.GetTestBillingAccountFromEnv(t),
		"user_project_override": "false", // Used in combo with billing_account
		"random_suffix":         randomString,
	}

	contextUserProjectOverrideTrue := map[string]interface{}{
		"org_id":                envvar.GetTestOrgFromEnv(t),
		"billing_account":       envvar.GetTestBillingAccountFromEnv(t),
		"user_project_override": "true", // Used in combo with billing_account
		"random_suffix":         randomString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Setup resources
				// Neither user_project_override nor billing_project value used here
				Config: testAccSdkProvider_billing_project_useBillingProject_setup(contextUserProjectOverrideFalse),
			},
			{
				// With user_project_override=true the PubSub topic CANNOT be provisioned because quota is consumed
				// from the newly provisioned project, and that project does not have the PubSub API enabled.
				// The billing_project is used, leading to the error occurring, because user_project_override=true
				Config:      testAccSdkProvider_billing_project_useBillingProject_scenario(contextUserProjectOverrideTrue),
				ExpectError: regexp.MustCompile(fmt.Sprintf("Error 403: Cloud Pub/Sub API has not been used in project tf-test-%s", randomString)),
			},
			{
				// With user_project_override=false the PubSub topic can be provisioned because quota is consumed
				// from the project the Terraform identity is in, and that project has PubSub API enabled.
				// The billing_project value isn't used, meaning the error doesn't happen, because user_project_override=false
				Config: testAccSdkProvider_billing_project_useBillingProject_scenario(contextUserProjectOverrideFalse),
			},
		},
	})
}

func testAccSdkProvider_billing_project_affectedByClientLibraryEnv(t *testing.T) {
	// Test cannot run in VCR mode due to use of aliases
	// See: https://github.com/hashicorp/terraform-provider-google/issues/20019
	acctest.SkipIfVcr(t)

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"org_id":                envvar.GetTestOrgFromEnv(t),
		"billing_account":       envvar.GetTestBillingAccountFromEnv(t),
		"user_project_override": "true", // Used in combo with billing_account
		"random_suffix":         randomString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Setup resources
				// Neither user_project_override nor billing_project value used here
				Config: testAccSdkProvider_billing_project_useBillingProject_setupWithApiEnabled(context),
			},
			{
				// This ENV interferes with setting the billing_project,
				// so we get an error mentioning the value
				PreConfig: func() {
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", "foobar")
				},
				Config:      testAccSdkProvider_billing_project_useBillingProject_scenarioWithApiEnabled(context),
				ExpectError: regexp.MustCompile("foobar"),
			},
			{
				// The same config without that ENV present applies without error
				PreConfig: func() {
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", "")
				},
				Config: testAccSdkProvider_billing_project_useBillingProject_scenarioWithApiEnabled(context),
			},
		},
	})
}

// testAccSdkProvider_billing_project_inProviderBlock allows setting the billing_project argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_billing_project_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	billing_project = "%{billing_project}"
}

data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_billing_project_inEnvsOnly allows testing when the billing_project argument
// is only supplied via ENVs
func testAccSdkProvider_billing_project_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}
`, context)
}

func testAccSdkProvider_billing_project_useBillingProject_setup(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {}

# Create a new project and enable service APIs in those projects
resource "google_project" "project" {
  provider = google
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "serviceusage" {
  project  = google_project.project.project_id
  service  = "serviceusage.googleapis.com"

  disable_on_destroy = false # Need it enabled in the project when the test disables services in post-test cleanup
}
`, context)
}

func testAccSdkProvider_billing_project_useBillingProject_scenario(context map[string]interface{}) string {

	// SECOND APPLY
	// This is needed as configuring the provider depends on resources provisioned in the setup step
	return testAccSdkProvider_billing_project_useBillingProject_setup(context) + acctest.Nprintf(`
# Set up the usage of
#  - user_project_override
#  - billing_project
provider "google" {
  alias                 = "user_project_override"
  user_project_override = %{user_project_override}
  billing_project       = google_project.project.project_id
  project               = google_project.project.project_id
}

# See if the impersonated SA can provision the PubSub resource in a way that uses
# the newly provisioned project as the source of consumed quota
resource "google_pubsub_topic" "example-resource-in" {
  provider = google.user_project_override
  project  = google_project.project.project_id
  name     = "tf-test-%{random_suffix}"
}
`, context)
}

// testAccSdkProvider_billing_project_useBillingProject_setupWithApiEnabled is the same setup as above but appends config to activate
// the PubSub API. This allows the second apply step to succeed in a test, if needed.
func testAccSdkProvider_billing_project_useBillingProject_setupWithApiEnabled(context map[string]interface{}) string {
	return testAccSdkProvider_billing_project_useBillingProject_setup(context) + acctest.Nprintf(`

# Needed for test steps to apply without error
resource "google_project_service" "pubsub" {
  project  = google_project.project.project_id
  service  = "pubsub.googleapis.com"
}

resource "google_project_service" "cloudresourcemanager" {
  project  = google_project.project.project_id
  service  = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false # Need it enabled in the project when the test deletes the project resource in post-test cleanup
}
`, context)
}

// testAccSdkProvider_billing_project_useBillingProject_scenarioWithApiEnabled is the same scenario as above but includes config that
// has activated the PubSub API. This allows the scenario to apply successfully in a test, if needed.
func testAccSdkProvider_billing_project_useBillingProject_scenarioWithApiEnabled(context map[string]interface{}) string {

	// SECOND APPLY
	// This is needed as configuring the provider depends on resources provisioned in the setup step
	return testAccSdkProvider_billing_project_useBillingProject_setupWithApiEnabled(context) + acctest.Nprintf(`
# Set up the usage of
#  - user_project_override
#  - billing_project
provider "google" {
  alias                 = "user_project_override"
  user_project_override = %{user_project_override}
  billing_project       = google_project.project.project_id
  project               = google_project.project.project_id
}

# See if the impersonated SA can provision the PubSub resource in a way that uses
# the newly provisioned project as the source of consumed quota
resource "google_pubsub_topic" "example-resource-in" {
  provider = google.user_project_override
  project  = google_project.project.project_id
  name     = "tf-test-%{random_suffix}"

  depends_on = [
    google_project_service.pubsub
  ]
}
`, context)
}
