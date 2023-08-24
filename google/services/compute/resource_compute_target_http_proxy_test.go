// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeTargetHttpProxy_update(t *testing.T) {
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
				Config: testAccComputeTargetHttpProxy_basic1(target, backend, hc, urlmap1, urlmap2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpProxyExists(
						t, "google_compute_target_http_proxy.foobar"),
				),
			},

			{
				Config: testAccComputeTargetHttpProxy_basic2(target, backend, hc, urlmap1, urlmap2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpProxyExists(
						t, "google_compute_target_http_proxy.foobar"),
				),
			},
		},
	})
}

func testAccCheckComputeTargetHttpProxyExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).TargetHttpProxies.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("TargetHttpProxy not found")
		}

		return nil
	}
}

func testAccComputeTargetHttpProxy_basic1(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_target_http_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_url_map.foobar1.self_link
}

resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, target, backend, hc, urlmap1, urlmap2)
}

func testAccComputeTargetHttpProxy_basic2(target, backend, hc, urlmap1, urlmap2 string) string {
	return fmt.Sprintf(`
resource "google_compute_target_http_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_url_map.foobar2.self_link
}

resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar1" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_url_map" "foobar2" {
  name            = "%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, target, backend, hc, urlmap1, urlmap2)
}
