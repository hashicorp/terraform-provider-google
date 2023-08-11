// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeTargetGrpcProxy_update(t *testing.T) {
	t.Parallel()

	proxy := fmt.Sprintf("tf-manual-proxy-%s", acctest.RandString(t, 10))
	urlmap1 := fmt.Sprintf("tf-manual-urlmap1-%s", acctest.RandString(t, 10))
	urlmap2 := fmt.Sprintf("tf-manual-urlmap2-%s", acctest.RandString(t, 10))
	backend := fmt.Sprintf("tf-manual-backend-%s", acctest.RandString(t, 10))
	healthcheck := fmt.Sprintf("tf-manual-healthcheck-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetGrpcProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetGrpcProxy_basic(proxy, urlmap1, backend, healthcheck),
			},
			{
				ResourceName:      "google_compute_target_grpc_proxy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeTargetGrpcProxy_basic(proxy, urlmap2, backend, healthcheck),
			},
			{
				ResourceName:      "google_compute_target_grpc_proxy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeTargetGrpcProxy_basic(proxy, urlmap, backend, healthcheck string) string {
	return fmt.Sprintf(`
resource "google_compute_target_grpc_proxy" "default" {
  name    = "%s"
  url_map = google_compute_url_map.urlmap.id
  validate_for_proxyless = true
}
resource "google_compute_url_map" "urlmap" {
  name        = "%s"
  description = "a description"
  default_service = google_compute_backend_service.home.id
  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }
  path_matcher {
    name = "allpaths"
    default_service = google_compute_backend_service.home.id
    route_rules {
      priority = 1
      header_action {
        request_headers_to_remove = ["RemoveMe2"]
        request_headers_to_add {
          header_name = "AddSomethingElse"
          header_value = "MyOtherValue"
          replace = true
        }
        response_headers_to_remove = ["RemoveMe3"]
        response_headers_to_add {
          header_name = "AddMe"
          header_value = "MyValue"
          replace = false
        }
      }
      match_rules {
        full_path_match = "a full path"
        header_matches {
          header_name = "someheader"
          exact_match = "match this exactly"
          invert_match = true
        }
        ignore_case = true
        metadata_filters {
          filter_match_criteria = "MATCH_ANY"
          filter_labels {
            name = "PLANET"
            value = "MARS"
          }
        }
        query_parameter_matches {
          name = "a query parameter"
          present_match = true
        }
      }
      url_redirect {
        host_redirect = "A host"
        https_redirect = false
        path_redirect = "some/path"
        redirect_response_code = "TEMPORARY_REDIRECT"
        strip_query = true
      }
    }
  }
  test {
    service = google_compute_backend_service.home.id
    host    = "hi.com"
    path    = "/home"
  }
}
resource "google_compute_backend_service" "home" {
  name        = "%s"
  port_name   = "grpc"
  protocol    = "GRPC"
  timeout_sec = 10
  health_checks = [google_compute_health_check.default.id]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}
resource "google_compute_health_check" "default" {
  name               = "%s"
  timeout_sec        = 1
  check_interval_sec = 1
  grpc_health_check {
    port_name          = "health-check-port"
    port_specification = "USE_NAMED_PORT"
    grpc_service_name  = "testservice"
  }
}
`, proxy, urlmap, backend, healthcheck)
}
