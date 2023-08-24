// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRegionTargetHttpProxy_update(t *testing.T) {
	t.Parallel()

	target := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	backend := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	hc := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	urlmap1 := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))
	urlmap2 := fmt.Sprintf("thttp-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionTargetHttpProxy_basic1(target, backend, hc, urlmap1, urlmap2),
			},
			{
				ResourceName:      "google_compute_region_target_http_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionTargetHttpProxy_basic2(target, backend, hc, urlmap1, urlmap2),
			},
			{
				ResourceName:      "google_compute_region_target_http_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionTargetHttpProxy_basic1(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_http_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_region_url_map.foobar1.self_link
}

resource "google_compute_region_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_region_health_check.zero.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name     = "%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_region_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar.self_link
  }
}

resource "google_compute_region_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_region_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar.self_link
  }
}
`, target, backend, hc, urlmap1, urlmap2)
}

func testAccComputeRegionTargetHttpProxy_basic2(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_http_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_region_url_map.foobar2.self_link
}

resource "google_compute_region_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_region_health_check.zero.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name     = "%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_region_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar.self_link
  }
}

resource "google_compute_region_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_region_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar.self_link
  }
}
`, target, backend, hc, urlmap1, urlmap2)
}
