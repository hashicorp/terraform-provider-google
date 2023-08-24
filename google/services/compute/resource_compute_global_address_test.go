// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_ipv6(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeGlobalAddress_internal(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_internal(acctest.RandString(t, 10), acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeGlobalAddress_ipv6(addressName string) string {
	return fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
  name        = "tf-test-address-%s"
  description = "Created for Terraform acceptance testing"
  ip_version  = "IPV6"
}
`, addressName)
}

func testAccComputeGlobalAddress_internal(networkName, addressName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-address-%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "tf-test-address-%s"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 24
  address       = "172.20.181.0"
  network       = google_compute_network.foobar.self_link
}
`, networkName, addressName)
}
