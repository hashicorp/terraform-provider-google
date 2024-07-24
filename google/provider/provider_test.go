// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = provider.Provider()
}

func TestProvider_noDuplicatesInResourceMap(t *testing.T) {
	_, err := provider.ResourceMapWithErrors()
	if err != nil {
		t.Error(err)
	}
}

func TestProvider_noDuplicatesInDatasourceMap(t *testing.T) {
	_, err := provider.DatasourceMapWithErrors()
	if err != nil {
		t.Error(err)
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", acctest.RandString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderMeta_setModuleName(t *testing.T) {
	t.Parallel()

	moduleName := "my-module"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderMeta_setModuleName(moduleName, acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billing := envvar.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	topicName := "tf-test-topic-" + acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	accessToken, err := acctest.SetupProjectsAndGetAccessToken(org, billing, pid, "pubsub", config)
	if err != nil {
		t.Error(err)
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
func TestAccProviderIndirectUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billing := envvar.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	accessToken, err := acctest.SetupProjectsAndGetAccessToken(org, billing, pid, "cloudkms", config)
	if err != nil {
		t.Error(err)
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

func TestAccProviderCredentialsEmptyString(t *testing.T) {
	// Test is not parallel because ENVs are set.
	// Need to skip VCR as this test downloads providers from the Terraform Registry
	acctest.SkipIfVcr(t)

	creds := envvar.GetTestCredsFromEnv()
	project := envvar.GetTestProjectFromEnv()
	t.Setenv("GOOGLE_CREDENTIALS", creds)
	t.Setenv("GOOGLE_PROJECT", project)

	pid := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				// This is a control for the other test steps; the provider block doesn't contain `credentials = ""`
				Config:                   testAccProviderCredentials_actWithCredsFromEnv(pid),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
			{
				// Assert that errors are expected with credentials when
				// - GOOGLE_CREDENTIALS is set
				// - provider block has credentials = ""
				// - TPG v4.60.2 is used
				// Context: this was an addidental breaking change introduced with muxing
				Config: testAccProviderCredentials_actWithCredsFromEnv_emptyString(pid),
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.60.2",
						Source:            "hashicorp/google",
					},
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`unexpected end of JSON input`),
			},
			{
				// Assert that errors are NOT expected with credentials when
				// - GOOGLE_CREDENTIALS is set
				// - provider block has credentials = ""
				// - TPG v4.84.0 is used
				// Context: this was the fix for the unintended breaking change in 4.60.2
				Config: testAccProviderCredentials_actWithCredsFromEnv_emptyString(pid),
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.84.0",
						Source:            "hashicorp/google",
					},
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Validation errors are expected in 5.0.0+
				// Context: we intentionally introduced the breaking change again in 5.0.0+
				Config:                   testAccProviderCredentials_actWithCredsFromEnv_emptyString(pid),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
				ExpectError:              regexp.MustCompile(`expected a non-empty string`),
			},
		},
	})
}

func TestAccProviderEmptyStrings(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			// When no values are set in the provider block there are no errors
			// This test case is a control to show validation doesn't accidentally flag unset fields
			// The "" argument is a lack of key = value being passed into the provider block
			{
				Config:             testAccProvider_checkPlanTimeErrors("", acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// credentials as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`credentials = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// access_token as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`access_token = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// impersonate_service_account as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`impersonate_service_account = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// project as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`project = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// billing_project as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`billing_project = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// region as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`region = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
			// zone as an empty string causes a validation error
			{
				Config:             testAccProvider_checkPlanTimeErrors(`zone = ""`, acctest.RandString(t, 10)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`expected a non-empty string`),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                   = "compute_custom_endpoint"
  compute_custom_endpoint = "%s"
}

resource "google_compute_address" "default" {
  provider = google.compute_custom_endpoint
  name     = "tf-test-address-%s"
}`, endpoint, name)
}

func testAccProviderMeta_setModuleName(key, name string) string {
	return fmt.Sprintf(`
terraform {
  provider_meta "google" {
    module_name = "%s"
  }
}

resource "google_compute_address" "default" {
	name = "tf-test-address-%s"
}`, key, name)
}

// Set up two projects. Project 1 has a service account that is used to create a
// pubsub topic in project 2. The pubsub API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.

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

// Copy the Mmv1 generated function testAccCheckComputeAddressDestroyProducer from the compute_test package to here,
// as that function is in the _test.go file and not importable.
func testAccCheckComputeAddressDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_address" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/addresses/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeAddress still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccProviderCredentials_actWithCredsFromEnv(name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias       = "testing_credentials"

}

resource "google_compute_address" "default" {
  provider = google.testing_credentials
  name     = "%s"
}`, name)
}

func testAccProviderCredentials_actWithCredsFromEnv_emptyString(name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias       = "testing_credentials"
  credentials = ""
}

resource "google_compute_address" "default" {
  provider = google.testing_credentials
  name     = "%s"
}`, name)
}

func testAccProvider_checkPlanTimeErrors(providerArgument, randString string) string {
	return fmt.Sprintf(`
provider "google" {
	%s
}

# A random resource so that the test can generate a plan (can't check validation errors when plan is empty)
resource "google_pubsub_topic" "example" {
  name = "tf-test-planned-resource-%s"
}
`, providerArgument, randString)
}
