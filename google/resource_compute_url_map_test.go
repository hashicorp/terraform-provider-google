package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeUrlMap_update_path_matcher(t *testing.T) {
	t.Parallel()

	bsName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	hcName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	umName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeUrlMap_basic1(bsName, hcName, umName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},

			{
				Config: testAccComputeUrlMap_basic2(bsName, hcName, umName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},
		},
	})
}

func TestAccComputeUrlMap_advanced(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeUrlMap_advanced1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},

			{
				Config: testAccComputeUrlMap_advanced2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},
		},
	})
}

func TestAccComputeUrlMap_noPathRulesWithUpdate(t *testing.T) {
	t.Parallel()

	bsName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	hcName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	umName := fmt.Sprintf("urlmap-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeUrlMap_noPathRules(bsName, hcName, umName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},
			{
				Config: testAccComputeUrlMap_basic1(bsName, hcName, umName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeUrlMapExists(
						"google_compute_url_map.foobar"),
				),
			},
		},
	})
}

func testAccCheckComputeUrlMapExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]

		found, err := config.clientCompute.UrlMaps.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("Url map not found")
		}
		return nil
	}
}

func TestAccComputeUrlMap_trafficDirectorUpdate(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(10)

	bsName := fmt.Sprintf("urlmap-test-%s", randString)
	hcName := fmt.Sprintf("urlmap-test-%s", randString)
	umName := fmt.Sprintf("urlmap-test-%s", randString)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeUrlMap_trafficDirector(bsName, hcName, umName),
			},
			{
				ResourceName:      "google_compute_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeUrlMap_trafficDirectorUpdate(bsName, hcName, umName),
			},
			{
				ResourceName:      "google_compute_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeUrlMap_trafficDirectorRemoveRouteRule(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(10)

	bsName := fmt.Sprintf("urlmap-test-%s", randString)
	hcName := fmt.Sprintf("urlmap-test-%s", randString)
	umName := fmt.Sprintf("urlmap-test-%s", randString)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeUrlMap_trafficDirector(bsName, hcName, umName),
			},
			{
				ResourceName:      "google_compute_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeUrlMap_trafficDirectorRemoveRouteRule(bsName, hcName, umName),
			},
			{
				ResourceName:      "google_compute_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeUrlMap_basic1(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "urlmap-test-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "urlmap-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "urlmap-test-%s"
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
`, bsName, hcName, umName)
}

func testAccComputeUrlMap_basic2(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "urlmap-test-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "urlmap-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_service = google_compute_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "blip"
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_backend_service.foobar.self_link
    }
  }

  test {
    host    = "mysite.com"
    path    = "/test"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, bsName, hcName, umName)
}

func testAccComputeUrlMap_advanced1() string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "urlmap-test-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "urlmap-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_service = google_compute_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "blop"
  }

  host_rule {
    hosts        = ["myfavoritesite.com"]
    path_matcher = "blip"
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blop"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeUrlMap_advanced2() string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "urlmap-test-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "urlmap-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_service = google_compute_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "blep"
  }

  host_rule {
    hosts        = ["myfavoritesite.com"]
    path_matcher = "blip"
  }

  host_rule {
    hosts        = ["myleastfavoritesite.com"]
    path_matcher = "blub"
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blep"

    path_rule {
      paths   = ["/home"]
      service = google_compute_backend_service.foobar.self_link
    }

    path_rule {
      paths   = ["/login"]
      service = google_compute_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blub"

    path_rule {
      paths   = ["/*", "/blub"]
      service = google_compute_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeUrlMap_noPathRules(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "urlmap-test-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "urlmap-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_service = google_compute_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }

  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
  }

  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}
`, bsName, hcName, umName)
}

func testAccComputeUrlMap_trafficDirector(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_url_map" "foobar" {
  name        = "%s"
  description = "a description"
  default_service = "${google_compute_backend_service.home.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name = "allpaths"
    default_service = "${google_compute_backend_service.home.self_link}"

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
    service = "${google_compute_backend_service.home.self_link}"
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_backend_service" "home" {
  name        = "%s"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_backend_service" "home2" {
  name        = "%s-2"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  http_health_check {
    port = 80
  }
}

`, umName, bsName, bsName, hcName)
}

func testAccComputeUrlMap_trafficDirectorUpdate(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_url_map" "foobar" {
  name        = "%s"
  description = "a description"
  default_service = "${google_compute_backend_service.home2.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths2"
  }

  path_matcher {
    name = "allpaths2"
    default_service = "${google_compute_backend_service.home2.self_link}"

    route_rules {
      priority = 2
      header_action {
        request_headers_to_remove = ["RemoveMe2", "AndMe"]
        request_headers_to_add {
          header_name = "AddSomethingElseUpdated"
          header_value = "MyOtherValueUpdated"
          replace = false
        }
        response_headers_to_remove = ["RemoveMe3", "AndMe4"]
      }
      match_rules {
        full_path_match = "a full path to match"
        header_matches {
          header_name = "someheaderfoo"
          exact_match = "match this exactly again"
          invert_match = false
        }
        ignore_case = false
        metadata_filters {
          filter_match_criteria = "MATCH_ALL"
          filter_labels {
            name = "PLANET"
            value = "EARTH"
          }
        }
      }
      url_redirect {
        host_redirect = "A host again"
        https_redirect = true
        path_redirect = "some/path/twice"
        redirect_response_code = "TEMPORARY_REDIRECT"
        strip_query = false
      }
    }
  }

  test {
    service = "${google_compute_backend_service.home.self_link}"
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_backend_service" "home" {
  name        = "%s"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_backend_service" "home2" {
  name        = "%s-2"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  http_health_check {
    port = 80
  }
}
`, umName, bsName, bsName, hcName)
}

func testAccComputeUrlMap_trafficDirectorRemoveRouteRule(bsName, hcName, umName string) string {
	return fmt.Sprintf(`
resource "google_compute_url_map" "foobar" {
  name        = "%s"
  description = "a description"
  default_service = "${google_compute_backend_service.home2.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths2"
  }

  path_matcher {
    name = "allpaths2"
    default_service = "${google_compute_backend_service.home2.self_link}"
  }

  test {
    service = "${google_compute_backend_service.home.self_link}"
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_backend_service" "home" {
  name        = "%s"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_backend_service" "home2" {
  name        = "%s-2"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  http_health_check {
    port = 80
  }
}
`, umName, bsName, bsName, hcName)
}
