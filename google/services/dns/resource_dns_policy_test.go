// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSPolicy_update(t *testing.T) {
	t.Parallel()

	policySuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsPolicy_privateUpdate(policySuffix, "true", "172.16.1.10", "172.16.1.30", "network-1"),
			},
			{
				ResourceName:      "google_dns_policy.example-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsPolicy_privateUpdate(policySuffix, "false", "172.16.1.20", "172.16.1.40", "network-2"),
			},
			{
				ResourceName:      "google_dns_policy.example-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsPolicy_privateUpdate(suffix, forwarding, first_nameserver, second_nameserver, network string) string {
	return fmt.Sprintf(`
resource "google_dns_policy" "example-policy" {
  name                      = "example-policy-%s"
  enable_inbound_forwarding = %s

  alternative_name_server_config {
    target_name_servers {
      ipv4_address = "%s"
    }
	target_name_servers {
	  ipv4_address    = "%s"
	  forwarding_path = "private"
    }
  }

  networks {
    network_url = google_compute_network.%s.self_link
  }
}

resource "google_compute_network" "network-1" {
  name                    = "tf-test-network-1-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  name                    = "tf-test-network-2-%s"
  auto_create_subnetworks = false
}
`, suffix, forwarding, first_nameserver, second_nameserver, network, suffix, suffix)
}
