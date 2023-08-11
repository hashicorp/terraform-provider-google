// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkServicesEdgeCacheOrigin_updateAndImport(t *testing.T) {
	t.Parallel()
	name := "tf-test-origin-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesEdgeCacheOriginDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesEdgeCacheOrigin_update_0(name),
			},
			{
				ResourceName:            "google_network_services_edge_cache_origin.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccNetworkServicesEdgeCacheOrigin_update_1(name),
			},
			{
				ResourceName:            "google_network_services_edge_cache_origin.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}
func testAccNetworkServicesEdgeCacheOrigin_update_0(name string) string {
	return fmt.Sprintf(`
	resource "google_network_services_edge_cache_origin" "instance" {
		name                 = "%s"
		origin_address       = "gs://media-edge-default"
		description          = "The default bucket for media edge test"
		max_attempts         = 2
		labels = {
			a = "b"
		}
		retry_conditions = ["NOT_FOUND"]
		timeout {
			connect_timeout = "10s"
		}
	}
`, name)
}
func testAccNetworkServicesEdgeCacheOrigin_update_1(name string) string {
	return fmt.Sprintf(`
	resource "google_network_services_edge_cache_origin" "instance" {
		name                 = "%s"
		origin_address       = "gs://media-edge-fallback"
		description          = "The default bucket for media edge test"
		max_attempts         = 3
		retry_conditions     = ["FORBIDDEN"]
		timeout {
			connect_timeout = "9s"
			max_attempts_timeout = "14s"
			response_timeout = "29s"
			read_timeout = "13s"
		}
	}
`, name)
}
