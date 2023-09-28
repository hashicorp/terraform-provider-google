// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
