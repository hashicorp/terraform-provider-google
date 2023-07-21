// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFrameworkProviderMeta_setModuleName(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	moduleName := "my-module"
	managedZoneName := fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkProviderMeta_setModuleName(moduleName, managedZoneName, acctest.RandString(t, 10)),
			},
		},
	})
}

func TestAccFrameworkProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", acctest.RandString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", acctest.RandString(t, 10)),
				ExpectError:              regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccFrameworkProviderBasePath_setBasePath(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccFrameworkProviderBasePath_setBasePath("https://www.googleapis.com/dns/v1beta2/", acctest.RandString(t, 10)),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				ResourceName:      "google_dns_managed_zone.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccFrameworkProviderBasePath_setBasePath("https://www.googleapis.com/dns/v1beta2/", acctest.RandString(t, 10)),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_dns_managed_zone.foo",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccFrameworkProviderBasePath_setBasePathstep3("https://www.googleapis.com/dns/v1beta2/", acctest.RandString(t, 10)),
			},
		},
	})
}

func testAccFrameworkProviderMeta_setModuleName(key, managedZoneName, recordSetName string) string {
	return fmt.Sprintf(`
terraform {
  provider_meta "google" {
    module_name = "%s"
  }
}

provider "google" {}

resource "google_dns_managed_zone" "zone" {
  name     = "%s-hashicorptest-com"
  dns_name = "%s.hashicorptest.com."
}

resource "google_dns_record_set" "rs" {
  managed_zone = google_dns_managed_zone.zone.name
  name         = "%s.${google_dns_managed_zone.zone.dns_name}"
  type         = "A"
  ttl          = 300
  rrdatas      = [
  "192.168.1.0",
  ]
}

data "google_dns_record_set" "rs" {
  managed_zone = google_dns_record_set.rs.managed_zone
  name         = google_dns_record_set.rs.name
  type         = google_dns_record_set.rs.type
}`, key, managedZoneName, managedZoneName, recordSetName)
}

func testAccFrameworkProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias               = "dns_custom_endpoint"
  dns_custom_endpoint = "%s"
}

resource "google_dns_managed_zone" "foo" {
  provider    = google.dns_custom_endpoint
  name        = "tf-test-zone-%s"
  dns_name    = "tf-test-zone-%s.hashicorptest.com."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  provider    = google.dns_custom_endpoint
  name = google_dns_managed_zone.foo.name
}`, endpoint, name, name)
}

func testAccFrameworkProviderBasePath_setBasePathstep3(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias               = "dns_custom_endpoint"
  dns_custom_endpoint = "%s"
}

resource "google_dns_managed_zone" "foo" {
  provider    = google.dns_custom_endpoint
  name        = "tf-test-zone-%s"
  dns_name    = "tf-test-zone-%s.hashicorptest.com."
  description = "QA DNS zone"
}
`, endpoint, name, name)
}

// Copy the function from the provider_test package to here
// as that function is in the _test.go file and not importable
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

// Copy the function testAccCheckComputeAddressDestroyProducer from the dns_test package to here,
// as that function is in the _test.go file and not importable.
//
// testAccCheckDNSManagedZoneDestroyProducerFramework is the framework version of the generated testAccCheckDNSManagedZoneDestroyProducer
// when we automate this, we'll use the automated version and can get rid of this
func testAccCheckDNSManagedZoneDestroyProducerFramework(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_managed_zone" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			p := acctest.GetFwTestProvider(t)

			url, err := acctest.ReplaceVarsForFrameworkTest(&p.FrameworkProvider.FrameworkProviderConfig, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.BillingProject.IsNull() && p.BillingProject.String() != "" {
				billingProject = p.BillingProject.String()
			}

			_, diags := fwtransport.SendFrameworkRequest(&p.FrameworkProvider.FrameworkProviderConfig, "GET", billingProject, url, p.UserAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSManagedZone still exists at %s", url)
			}
		}

		return nil
	}
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
