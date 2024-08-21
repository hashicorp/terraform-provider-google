// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeAddress_networkTier(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "effective_labels.%"),
				),
			},
		},
	})
}

func TestAccComputeAddress_internal(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_internal(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.internal",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet_and_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_networkTier_withLabels(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeAddress_networkTier_withLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.default_expiration_ms", "3600000"),
				),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_networkTier_withLabelsUpdate(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.default_expiration_ms", "7200000"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_networkTier_withProvider5(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "effective_labels.%"),
				),
			},
			{
				Config: testAccComputeAddress_networkTier_withLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.default_expiration_ms", "3600000"),
				),
			},
		},
	})
}

func TestAccComputeAddress_withProviderDefaultLabels(t *testing.T) {
	// The test failed if VCR testing is enabled, because the cached provider config is used.
	// With the cached provider config, any changes in the provider default labels will not be applied.
	acctest.SkipIfVcr(t)
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_withProviderDefaultLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_resourceLabelsOverridesProviderDefaultLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_moveResourceLabelToProviderDefaultLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_resourceLabelsOverridesProviderDefaultLabels(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_withCreationOnlyAttribution(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Creating with two user supplied labels should result in those labels + the attribution label.
				Config: testAccComputeAddress_networkTier_withAttribution(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// Updating the user supplied labels should leave the attribution label intact.
				Config: testAccComputeAddress_networkTier_withAttributionUpdate(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				// Removing the user supplied labels should leave the attribution label intact.
				Config: testAccComputeAddress_networkTier_withAttributionClear(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "1"),
				),
			},
		},
	})
}

func TestAccComputeAddress_withCreationOnlyAttributionSetOnUpdate(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create the initial resource without the attribution label.
				Config: testAccComputeAddress_networkTier_withSkipAttribution(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
				),
			},
			{
				// Updating with attribution label set to "CREATION_ONLY" should not add the label.
				Config: testAccComputeAddress_networkTier_withAttributionUpdate(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
				),
			},
		},
	})
}

func TestAccComputeAddress_withProactiveAttribution(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Creating with two user supplied labels should result in those labels + the attribution label.
				Config: testAccComputeAddress_networkTier_withAttribution(suffix, "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_compute_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// Updating the user supplied labels should leave the attribution label intact.
				Config: testAccComputeAddress_networkTier_withAttributionUpdate(suffix, "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				// Removing the user supplied labels should leave the attribution label intact.
				Config: testAccComputeAddress_networkTier_withAttributionClear(suffix, "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_address.foobar", "labels.%"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "1"),
				),
			},
		},
	})
}

func TestAccComputeAddress_withProactiveAttributionSetOnUpdate(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create the initial resource without the attribution label.
				Config: testAccComputeAddress_networkTier_withSkipAttribution(suffix, "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
				),
			},
			{
				// Updating with attribution label set to "PROACTIVE" should add the label.
				Config: testAccComputeAddress_networkTier_withAttributionUpdate(suffix, "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
		},
	})
}

func TestAccComputeAddress_withAttributionRemoved(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Creating with two user supplied labels should result in those labels + the attribution label.
				Config: testAccComputeAddress_networkTier_withAttribution(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "3"),
				),
			},
			{
				// Skipping attribution on resources that already have attribution removes the previous attribution.
				Config: testAccComputeAddress_networkTier_withSkipAttributionUpdate(suffix, "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_compute_address.foobar", "terraform_labels.default_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_compute_address.foobar", "effective_labels.%", "2"),
				),
			},
		},
	})
}

func testAccComputeAddress_networkTier_withLabels(i string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
  }
}
`, i)
}

func testAccComputeAddress_networkTier_withLabelsUpdate(i string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "bar"
    default_expiration_ms = 7200000
  }
}
`, i)
}

func testAccComputeAddress_withProviderDefaultLabels(i string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
  }
}
`, i)
}

func testAccComputeAddress_resourceLabelsOverridesProviderDefaultLabels(i string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
    default_key1          = "value1"
  }
}
`, i)
}

func testAccComputeAddress_moveResourceLabelToProviderDefaultLabels(i string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
    env          = "foo"
  }
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    default_expiration_ms = 3600000
    default_key1          = "value1"
  }
}
`, i)
}

func testAccComputeAddress_networkTier_withAttribution(suffix, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label               = true
  terraform_attribution_label_addition_strategy = %q
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
  }
}
`, strategy, suffix)
}

func testAccComputeAddress_networkTier_withSkipAttribution(suffix, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label               = false
  terraform_attribution_label_addition_strategy = %q
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
  }
}
`, strategy, suffix)
}

func testAccComputeAddress_networkTier_withAttributionUpdate(suffix, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label               = true
  terraform_attribution_label_addition_strategy = %q
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "bar"
    default_expiration_ms = 7200000
  }
}
`, strategy, suffix)
}

func testAccComputeAddress_networkTier_withSkipAttributionUpdate(suffix, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label               = false
  terraform_attribution_label_addition_strategy = %q
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"

  labels = {
    env                   = "bar"
    default_expiration_ms = 7200000
  }
}
`, strategy, suffix)
}

func testAccComputeAddress_networkTier_withAttributionClear(suffix, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label               = true
  terraform_attribution_label_addition_strategy = %q
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"
}
`, strategy, suffix)
}

func testAccComputeAddress_internal(i string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "internal" {
  name         = "tf-test-address-internal-%s"
  address_type = "INTERNAL"
  region       = "us-east1"
}

resource "google_compute_network" "default" {
  name = "tf-test-network-test-%s"
}

resource "google_compute_subnetwork" "foo" {
  name          = "subnetwork-test-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_address" "internal_with_subnet" {
  name         = "tf-test-address-internal-with-subnet-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  region       = "us-east1"
}

// We can't test the address alone, because we don't know what IP range the
// default subnetwork uses.
resource "google_compute_address" "internal_with_subnet_and_address" {
  name         = "tf-test-address-internal-with-subnet-and-address-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  address      = "10.0.42.42"
  region       = "us-east1"
}
`,
		i, // google_compute_address.internal name
		i, // google_compute_network.default name
		i, // google_compute_subnetwork.foo name
		i, // google_compute_address.internal_with_subnet_name
		i, // google_compute_address.internal_with_subnet_and_address name
	)
}

func testAccComputeAddress_networkTier(i string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"
}
`, i)
}

func TestAccComputeAddress_internalIpv6(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_internalIpv6(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.ipv6",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeAddress_internalIpv6(i string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name                     = "tf-test-network-test-%s"
  enable_ula_internal_ipv6 = true
  auto_create_subnetworks  = false
}
resource "google_compute_subnetwork" "foo" {
  name             = "subnetwork-test-%s"
  ip_cidr_range    = "10.0.0.0/16"
  region           = "us-east1"
  network          = google_compute_network.default.self_link
  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "INTERNAL"
}
resource "google_compute_address" "ipv6" {
  name         = "tf-test-address-internal-ipv6-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  region       = "us-east1"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  ip_version   = "IPV6"
}
`,
		i, // google_compute_network.default name
		i, // google_compute_subnetwork.foo name
		i, // google_compute_address.ipv6
	)
}
