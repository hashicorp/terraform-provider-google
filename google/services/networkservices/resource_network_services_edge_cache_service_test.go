// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkServicesEdgeCacheService_updateAndImport(t *testing.T) {
	t.Parallel()
	namebkt := "tf-test-bucket-" + acctest.RandString(t, 10)
	nameorigin := "tf-test-origin-" + acctest.RandString(t, 10)
	nameservice := "tf-test-service-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesEdgeCacheServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesEdgeCacheService_update_0(namebkt, nameorigin, nameservice),
			},
			{
				ResourceName:      "google_network_services_edge_cache_service.served",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkServicesEdgeCacheService_update_1(namebkt, nameorigin, nameservice),
			},
			{
				ResourceName:      "google_network_services_edge_cache_service.served",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccNetworkServicesEdgeCacheService_update_0(bktName, originName, serviceName string) string {
	return fmt.Sprintf(`
	resource "google_storage_bucket" "dest" {
		name          = "%s"
		location      = "US"
		force_destroy = true
	}
	resource "google_network_services_edge_cache_origin" "instance" {
		name                 = "%s"
		origin_address       = google_storage_bucket.dest.url
		description          = "The default bucket for media edge test"
		max_attempts         = 2
		timeout {
			connect_timeout = "10s"
		}
	}
	resource "google_network_services_edge_cache_service" "served" {
		name                 = "%s"
		description          = "some description"
		routing {
			host_rule {
				description = "host rule description"
				hosts = ["sslcert.tf-test.club"]
				path_matcher = "routes"
			}
			path_matcher {
				name = "routes"
				route_rule {
					description = "a route rule to match against"
					priority = 1
					match_rule {
						prefix_match = "/"
					}
					origin = google_network_services_edge_cache_origin.instance.name
					route_action {
						cdn_policy {
								cache_mode = "CACHE_ALL_STATIC"
								default_ttl = "3600s"
						}
					}
					header_action {
						response_header_to_add {
							header_name = "x-cache-status"
							header_value = "{cdn_cache_status}"
						}
					}
				}
			}
		}
	}
`, bktName, originName, serviceName)
}
func testAccNetworkServicesEdgeCacheService_update_1(bktName, originName, serviceName string) string {
	return fmt.Sprintf(`
	resource "google_storage_bucket" "dest" {
		name          = "%s"
		location      = "US"
		force_destroy = true
	}
	resource "google_network_services_edge_cache_origin" "instance" {
		name                 = "%s"
		origin_address       = google_storage_bucket.dest.url
		description          = "The default bucket for media edge test"
		max_attempts         = 2
		timeout {
			connect_timeout = "10s"
		}
	}
	resource "google_network_services_edge_cache_service" "served" {
		name                 = "%s"
		description          = "some description"
		routing {
			host_rule {
				description = "host rule description"
				hosts = ["sslcert.tf-test.club"]
				path_matcher = "routes"
			}
			path_matcher {
				name = "routes"
				route_rule {
					description = "a route rule to match against"
					priority = 1
					match_rule {
						prefix_match = "/"
					}
					origin = google_network_services_edge_cache_origin.instance.name
					route_action {
						cdn_policy {
								cache_mode = "CACHE_ALL_STATIC"
								default_ttl = "3600s"
						}
					}
					header_action {
						response_header_to_add {
							header_name = "x-cache-status"
							header_value = "{cdn_cache_status}"
						}
					}
				}
			}
		}
	}
`, bktName, originName, serviceName)
}
