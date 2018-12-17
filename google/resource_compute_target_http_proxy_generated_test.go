// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeTargetHttpProxy_targetHttpProxyBasicExample(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetHttpProxy_targetHttpProxyBasicExample(acctest.RandString(10)),
			},
			{
				ResourceName:      "google_compute_target_http_proxy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeTargetHttpProxy_targetHttpProxyBasicExample(val string) string {
	return fmt.Sprintf(`
resource "google_compute_target_http_proxy" "default" {
  name        = "test-proxy-%s"
  url_map     = "${google_compute_url_map.default.self_link}"
}

resource "google_compute_url_map" "default" {
  name        = "url-map-%s"
  default_service = "${google_compute_backend_service.default.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = "${google_compute_backend_service.default.self_link}"

    path_rule {
      paths   = ["/*"]
      service = "${google_compute_backend_service.default.self_link}"
    }
  }
}

resource "google_compute_backend_service" "default" {
  name        = "backend-service-%s"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_http_health_check" "default" {
  name               = "http-health-check-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, val, val, val, val,
	)
}

func testAccCheckComputeTargetHttpProxyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_http_proxy" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(rs, "https://www.googleapis.com/compute/v1/projects/{{project}}/global/targetHttpProxies/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", url, nil)
		if err == nil {
			return fmt.Errorf("ComputeTargetHttpProxy still exists at %s", url)
		}
	}

	return nil
}
