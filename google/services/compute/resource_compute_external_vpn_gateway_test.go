// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeExternalVPNGateway_updateLabels(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	resourceName := "google_compute_external_vpn_gateway.external_gateway"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeExternalVPNGateway_updateLabels(rnd, "test", "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.test", "test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeExternalVPNGateway_updateLabels(rnd, "test-updated", "test-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.test-updated", "test-updated"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccComputeExternalVPNGateway_updateLabels(suffix, key, value string) string {
	return fmt.Sprintf(`
resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "tf-test-external-gateway-%s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }

  labels = {
    %s = "%s"
  }
}
`, suffix, key, value)
}

func TestAccComputeExternalVPNGateway_insertIpv6Address(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	resourceName := "google_compute_external_vpn_gateway.external_gateway"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: computeExternalVPNGatewayIpv6AddressConfig(rnd, "2001:db8:abcd:1234:5678:90ab:cdef:1234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interface.0.ipv6_address", "2001:db8:abcd:1234:5678:90ab:cdef:1234"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func computeExternalVPNGatewayIpv6AddressConfig(suffix, ipv6_address string) string {
	return fmt.Sprintf(`
resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "tf-test-external-gateway-%s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id           = 0
    ipv6_address = "%s"
  }
}
`, suffix, ipv6_address)
}
