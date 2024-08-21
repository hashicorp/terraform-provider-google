// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccComputeGlobalAddress_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_update1(context),
			},
			{
				ResourceName:            "google_compute_global_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeGlobalAddress_update2(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_compute_global_address.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

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

func testAccComputeGlobalAddress_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-address-%{random_suffix}"
}

resource "google_compute_global_address" "foobar" {
  address       = "172.20.181.0"
  description   = "Description"
  name          = "tf-test-address-%{random_suffix}"
  labels        = {
  	foo = "bar"
  }
  ip_version     = "IPV4"
  prefix_length = 24
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  network       = google_compute_network.foobar.self_link
}
`, context)
}

func testAccComputeGlobalAddress_update2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-address-%{random_suffix}"
}

resource "google_compute_global_address" "foobar" {
  address       = "172.20.181.0"
  description   = "Description"
  name          = "tf-test-address-%{random_suffix}"
  labels        = {
  	foo = "baz"
  }
  ip_version     = "IPV4"
  prefix_length = 24
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  network       = google_compute_network.foobar.self_link
}
`, context)
}
