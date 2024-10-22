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

// TestAccSdkProvider_user_project_override is a series of acc tests asserting how the plugin-framework provider handles credentials arguments
// It is PF specific because the HCL used uses a PF-implemented data source
// It is a counterpart to TestAccFwProvider_user_project_override
func TestAccSdkProvider_user_project_override(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                testAccSdkProvider_user_project_override_configPrecedenceOverEnvironmentVariables,
		"when user_project_override is unset in the config, environment variables are used": testAccSdkProvider_user_project_override_precedenceOrderEnvironmentVariables,

		// Schema-level validation
		"when user_project_override is set in the config the value can be a boolean (true/false) or a string (true/false/1/0)": testAccSdkProvider_user_project_override_booleansInConfigOnly,
		"when user_project_override is set via environment variables any of these values can be used: true/false/1/0":          testAccSdkProvider_user_project_override_envStringsAccepted,

		// Usage
		"user_project_override uses a resource's project argument to control which project is used for quota and billing purposes":    testAccProviderUserProjectOverride,
		"user_project_override works for resources that don't take a project argument (provider-level default project value is used)": testAccProviderIndirectUserProjectOverride,
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

func testAccSdkProvider_user_project_override_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	override := "true"
	providerOverride := false

	// ensure all possible region env vars set; show they aren't used
	t.Setenv("USER_PROJECT_OVERRIDE", override)

	context := map[string]interface{}{
		"user_project_override": providerOverride,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", fmt.Sprintf("%v", providerOverride)),
				),
			},
		},
	})
}

func testAccSdkProvider_user_project_override_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for region, and they are in order of precedence.
		USER_PROJECT_OVERRIDE
	*/

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "") // unset
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// defaults to false when not set via config or ENVs
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "true")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
		},
	})
}

func testAccSdkProvider_user_project_override_booleansInConfigOnly(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context_true := map[string]interface{}{
		"user_project_override": true,
	}
	context_false := map[string]interface{}{
		"user_project_override": false,
	}

	context_1 := map[string]interface{}{
		"user_project_override": "1",
	}
	context_0 := map[string]interface{}{
		"user_project_override": "0",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context_true),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context_false),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_true),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_false),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_1),
				// No error expected
			},
			{
				Config: testAccSdkProvider_user_project_overrideInProviderBlock_string(context_0),
				// No error expected
			},
		},
	})
}

func testAccSdkProvider_user_project_override_envStringsAccepted(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "true")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "1")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "true"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "false")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
			{
				PreConfig: func() {
					t.Setenv("USER_PROJECT_OVERRIDE", "0")
				},
				Config: testAccSdkProvider_user_project_overrideInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "user_project_override", "false"),
				),
			},
		},
	})
}

// TestAccSdkProvider_user_project_overrideInProviderBlock allows setting the user_project_override argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the plugin-framework
func testAccSdkProvider_user_project_overrideInProviderBlock_boolean(context map[string]interface{}) string {
	v := acctest.Nprintf(`
provider "google" {
	user_project_override = %{user_project_override}
}

data "google_provider_config_sdk" "default" {}
`, context)
	return v
}

func testAccSdkProvider_user_project_overrideInProviderBlock_string(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	user_project_override = "%{user_project_override}"
}

data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_user_project_overrideInEnvsOnly allows testing when the user_project_override argument
// is only supplied via ENVs
func testAccSdkProvider_user_project_overrideInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}
`, context)
}

// Set up two projects. Project 1 has a service account that is used to create a
// pubsub topic in project 2. The pubsub API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.
// The test demonstrates how:
// - If user_project_override = false : the apply fails as the API is disabled in project-1
// - If user_project_override = true : the apply succeeds as X-Goog-User-Project will reference project-2, where API is enabled
func testAccProviderUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billing := envvar.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	topicName := "tf-test-topic-" + acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	accessToken, err := acctest.SetupProjectsAndGetAccessToken(org, billing, pid, "pubsub", config)
	if err != nil || accessToken == "" {
		if err == nil {
			t.Fatal("error when setting up projects and retrieving access token: access token is an empty string")
		}
		t.Fatalf("error when setting up projects and retrieving access token: %s", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderUserProjectOverride_step2(accessToken, pid, false, topicName),
				ExpectError: regexp.MustCompile("Cloud Pub/Sub API has not been used"),
			},
			{
				Config: testAccProviderUserProjectOverride_step2(accessToken, pid, true, topicName),
			},
			{
				ResourceName:            "google_pubsub_topic.project-2-topic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccProviderUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

// Do the same thing as TestAccProviderUserProjectOverride, but using a resource that gets its project via
// a reference to a different resource instead of a project field.
func testAccProviderIndirectUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billing := envvar.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	accessToken, err := acctest.SetupProjectsAndGetAccessToken(org, billing, pid, "cloudkms", config)
	if err != nil || accessToken == "" {
		if err == nil {
			t.Fatal("error when setting up projects and retrieving access token: access token is an empty string")
		}
		t.Fatalf("error when setting up projects and retrieving access token: %s", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, false),
				ExpectError: regexp.MustCompile(`Cloud Key Management Service \(KMS\) API has not been used`),
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, true),
			},
			{
				ResourceName:      "google_kms_crypto_key.project-2-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

func testAccProviderUserProjectOverride_step2(accessToken, pid string, override bool, topicName string) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the pubsub topic.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// pubsub topic so the whole config can be deleted.
%s

resource "google_pubsub_topic" "project-2-topic" {
	provider = google.project-1-token
	project  = "%s-2"

	name = "%s"
	labels = {
	  foo = "bar"
	}
}
`, testAccProviderUserProjectOverride_step3(accessToken, override), pid, topicName)
}

func testAccProviderUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias  = "project-1-token"
	access_token = "%s"
	user_project_override = %v
}
`, accessToken, override)
}

func testAccProviderIndirectUserProjectOverride_step2(pid, accessToken string, override bool) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the kms resources.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// kms resources so the whole config can be deleted.
%s

resource "google_kms_key_ring" "project-2-keyring" {
	provider = google.project-1-token
	project  = "%s-2"

	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "project-2-key" {
	provider = google.project-1-token
	name     = "%s"
	key_ring = google_kms_key_ring.project-2-keyring.id
}

data "google_kms_secret_ciphertext" "project-2-ciphertext" {
	provider   = google.project-1-token
	crypto_key = google_kms_crypto_key.project-2-key.id
	plaintext  = "my-secret"
}
`, testAccProviderIndirectUserProjectOverride_step3(accessToken, override), pid, pid, pid)
}

func testAccProviderIndirectUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias = "project-1-token"

	access_token          = "%s"
	user_project_override = %v
}
`, accessToken, override)
}
