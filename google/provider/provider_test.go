// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
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
				ExpectError: regexp.MustCompile("got HTTP response code [4-5][0-9]{2} with body"),
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
