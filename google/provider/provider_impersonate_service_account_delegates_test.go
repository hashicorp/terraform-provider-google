// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSdkProvider_impersonate_service_account_delegates(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		//     There are no environment variables for this field
		"impersonate_service_account_delegates can be set in config": testAccSdkProvider_impersonate_service_account_delegates_setInConfig,

		// Schema-level validation
		"when impersonate_service_account_delegates is set to an empty list in the config the value IS ignored": testAccSdkProvider_impersonate_service_account_delegates_emptyListUsage,

		// Usage
		"impersonate_service_account_delegates controls which service account is used for actions": testAccSdkProvider_impersonate_service_account_delegates_usage,
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

func testAccSdkProvider_impersonate_service_account_delegates_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	delegates := []string{
		"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
		"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
	}
	delegatesString := fmt.Sprintf(`["%s","%s"]`, delegates[0], delegates[1])

	// There are no ENVs for this provider argument

	context := map[string]interface{}{
		"random_suffix":                         acctest.RandString(t, 10),
		"impersonate_service_account_delegates": delegatesString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account_delegates.#", "2"),
				),
			},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_delegates_emptyListUsage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"random_suffix":                         acctest.RandString(t, 10),
		"impersonate_service_account_delegates": "[]",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account_delegates.#", "0"),
				),
			},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_delegates_usage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context),
			},
			{
				// This needs to be split into a second step as impersonate_service_account_delegates does
				// not tolerate unknown values
				Config:      testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_2(context),
				ExpectError: regexp.MustCompile("Error creating Topic: googleapi: Error 403: User not authorized"),
			},
		},
	})
}

// testAccSdkProvider_impersonate_service_account_delegates_testProvisioning allows setting the impersonate_service_account_delegates argument in a provider block
// and testing its impact on provisioning a resource
func testAccSdkProvider_impersonate_service_account_delegates_testProvisioning(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account_delegates = %{impersonate_service_account_delegates}
}

data "google_provider_config_sdk" "default" {}

resource "google_pubsub_topic" "example" {
  name = "tf-test-%{random_suffix}"
}
`, context)
}

func testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context map[string]interface{}) string {

	// This test config sets up the ability to use impersonate_service_account_delegates
	//    The 'base service account' is the service account that credentials supplied via ENVs is linked to.
	//    The 'delegate service account' is google_service_account.delegate
	//        The base SA is given roles/iam.serviceAccountTokenCreator on the delegate SA via google_service_account_iam_member.base_create_delegate_token
	//    The 'target service account' is google_service_account.target
	//        The delegate SA is given roles/iam.serviceAccountTokenCreator on the target SA via google_service_account_iam_member.delegate_create_target_token

	return acctest.Nprintf(`
// This will succeed due to the Terraform identity/base service account having necessary permissions
resource "google_pubsub_topic" "ok" {
  name = "tf-test-%{random_suffix}-ok"
}

//  Create a delegate service account and ensure the Terraform identity/base service account
//  can make tokens for it
resource "google_service_account" "delegate" {
  account_id   = "tf-test-%{random_suffix}-1"
  display_name = "Acceptance test impersonated service account"
}

data "google_client_openid_userinfo" "me" {
}

resource "google_service_account_iam_member" "base_create_delegate_token" {
  service_account_id = google_service_account.delegate.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${data.google_client_openid_userinfo.me.email}"
}

//  Create a second service account and ensure the first service account can make tokens for it
resource "google_service_account" "target" {
  account_id   = "tf-test-%{random_suffix}-2"
  display_name = "Acceptance test impersonated service account"
}

resource "google_service_account_iam_member" "delegate_create_target_token" {
  service_account_id = google_service_account.target.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.delegate.email}"
}

# Despite provisioning all the needed service accounts and permissions above
# this test sometimes fails with "Permission 'iam.serviceAccounts.getAccessToken' denied on resource (or it may not exist)"
# This error can be caused by either of:
#   - the IAM Service Account Credentials API not being enabled
#   - the service account not existing
#   - eventual consistency affecting IAM policies set on the service accounts
# Splitting this test into 2 steps is not sufficient to help with timing issues, so we add this sleep
resource "time_sleep" "wait_5_minutes" {
  depends_on = [
    google_service_account_iam_member.base_create_delegate_token,
    google_service_account_iam_member.delegate_create_target_token
  ]

  create_duration = "300s"	
}
`, context)
}

func testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_2(context map[string]interface{}) string {
	// See comments in testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1, about how the config
	// sets up the ability to use impersonate_service_account_delegates.

	// Here in testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_2 we:
	//    Pass the base service account to google.impersonation implicitly via `credentials` (ENVs in the test environment)
	//    Set the target service account as `impersonate_service_account`
	//    Set the delegate service account(s) in `impersonate_service_account_delegates`
	return testAccSdkProvider_impersonate_service_account_delegates_testViaFailure_1(context) + acctest.Nprintf(`
provider "google" {
  alias = "impersonation"
  impersonate_service_account = google_service_account.target.email
  impersonate_service_account_delegates = [
    google_service_account.delegate.email,
  ]
}

// This will fail due to the impersonated service account not having any permissions
resource "google_pubsub_topic" "fail" {
  provider = google.impersonation
  name = "tf-test-%{random_suffix}-fail"
}
`, context)
}
