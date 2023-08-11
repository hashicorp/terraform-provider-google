// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeRegionSslPolicy_update(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSslPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSslUpdate2(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "RESTRICTED"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_2"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionSslPolicy_update_to_custom(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSslPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSslUpdate3(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "CUSTOM"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_1"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionSslPolicy_update_from_custom(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSslPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSslUpdate3(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "CUSTOM"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_1"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionSslPolicyExists(
						t, "google_compute_region_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_region_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			{
				ResourceName:      "google_compute_region_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRegionSslPolicyExists(t *testing.T, n string, sslPolicy *compute.SslPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		project, err := acctest.GetTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).RegionSslPolicies.Get(
			project, "us-central1", name).Do()
		if err != nil {
			return fmt.Errorf("Error Reading SSL Policy %s: %s", name, err)
		}

		if found.Name != name {
			return fmt.Errorf("SSL Policy not found")
		}

		*sslPolicy = *found

		return nil
	}
}

func testAccComputeRegionSslUpdate1(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_ssl_policy" "update" {
  name            = "%s"
  description     = "Generated by TF provider acceptance test"
  min_tls_version = "TLS_1_0"
  profile         = "MODERN"
  region          = "us-central1"
}
`, resourceName)
}

func testAccComputeRegionSslUpdate2(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_ssl_policy" "update" {
  name            = "%s"
  description     = "Generated by TF provider acceptance test"
  min_tls_version = "TLS_1_2"
  profile         = "RESTRICTED"
  region          = "us-central1"
}
`, resourceName)
}

func testAccComputeRegionSslUpdate3(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_ssl_policy" "update" {
  name            = "%s"
  description     = "Generated by TF provider acceptance test"
  min_tls_version = "TLS_1_1"
  profile         = "CUSTOM"
  region          = "us-central1"
  custom_features = ["TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
}
`, resourceName)
}
