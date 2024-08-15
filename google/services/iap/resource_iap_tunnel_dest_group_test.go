// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iap_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIapTunnelDestGroup_updates(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckIapTunnelDestGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIapTunnelDestGroup_full(context),
			},
			{
				ResourceName:            "google_iap_tunnel_dest_group.dest_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "group_name"},
			},
			{
				Config: testAccIapTunnelDestGroup_updated(context),
			},
			{
				ResourceName:            "google_iap_tunnel_dest_group.dest_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "group_name"},
			},
			{
				Config: testAccIapTunnelDestGroup_updated_fqdns(context),
			},
			{
				ResourceName:            "google_iap_tunnel_dest_group.dest_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "group_name"},
			},
		},
	})
}

func testAccIapTunnelDestGroup_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iap_tunnel_dest_group" "dest_group" {
  region = "us-central1"
  group_name = "testgroup%{random_suffix}"
  cidrs = [
    "10.1.0.0/16",
    "192.168.10.0/24",
  ]
}
`, context)
}

func testAccIapTunnelDestGroup_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iap_tunnel_dest_group" "dest_group" {
  region = "us-central1"
  group_name = "testgroup%{random_suffix}"
  cidrs = [
    "10.1.0.0/16",
  ]
}
`, context)
}

func testAccIapTunnelDestGroup_updated_fqdns(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iap_tunnel_dest_group" "dest_group" {
  region = "us-central1"
  group_name = "testgroup%{random_suffix}"
  cidrs = [
    "10.1.0.0/16",
  ]
  fqdns = ["proxied.lan"]
}
`, context)
}
