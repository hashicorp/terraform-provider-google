package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSPolicy_update(t *testing.T) {
	t.Parallel()

	policySuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsPolicy_privateUpdate(policySuffix, "true", "172.16.1.10", "network-1"),
			},
			{
				ResourceName:      "google_dns_policy.example-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsPolicy_privateUpdate(policySuffix, "false", "172.16.1.20", "network-2"),
			},
			{
				ResourceName:      "google_dns_policy.example-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsPolicy_privateUpdate(suffix, forwarding, nameserver, network string) string {
	return fmt.Sprintf(`
resource "google_dns_policy" "example-policy" {
  name                      = "example-policy-%s"
  enable_inbound_forwarding = %s

  alternative_name_server_config {
    target_name_servers {
      ipv4_address = "%s"
    }
  }

  networks {
    network_url = google_compute_network.%s.self_link
  }
}

resource "google_compute_network" "network-1" {
  name                    = "network-1-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  name                    = "network-2-%s"
  auto_create_subnetworks = false
}
`, suffix, forwarding, nameserver, network, suffix, suffix)
}
