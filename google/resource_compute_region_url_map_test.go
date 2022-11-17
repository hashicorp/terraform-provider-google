package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRegionUrlMap_update_path_matcher(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_basic1(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_basic2(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_advanced(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_advanced1(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_advanced2(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_noPathRulesWithUpdate(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_noPathRules(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_basic1(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_ilbPathUpdate(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_ilbPath(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_ilbPathUpdate(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_ilbRouteUpdate(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_ilbRoute(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_ilbRouteUpdate(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_defaultUrlRedirect(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_defaultUrlRedirectConfig(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionUrlMap_defaultUrlRedirectWithinPathMatcher(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_defaultUrlRedirectWithinPathMatcherConfig(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Set all fields nested within `defaultRouteAction`, test import, then test updating all fields
func TestAccComputeRegionUrlMap_defaultRouteAction_full_update(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionUrlMap_defaultRouteAction_full(randomSuffix),
			},
			{
				ResourceName:      "google_compute_region_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionUrlMap_defaultRouteAction_full_update(randomSuffix),
			},
		},
	})
}

func testAccComputeRegionUrlMap_basic1(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  region        = "us-central1"
  name          = "regionurlmap-test-%s"
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  region   = "us-central1"
  name     = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}

resource "google_compute_region_url_map" "foobar" {
  region          = "us-central1"
  name            = "regionurlmap-test-%s"
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
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_basic2(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  region        = "us-central1"
  name          = "regionurlmap-test-%s"
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  region   = "us-central1"
  name     = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}

resource "google_compute_region_url_map" "foobar" {
  region          = "us-central1"
  name            = "regionurlmap-test-%s"
  default_service = google_compute_region_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "blip"
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }

  test {
    host    = "mysite.com"
    path    = "/test"
    service = google_compute_region_backend_service.foobar.self_link
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_advanced1(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  region        = "us-central1"
  name          = "regionurlmap-test-%s"
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  region   = "us-central1"
  name     = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}

resource "google_compute_region_url_map" "foobar" {
  region          = "us-central1"
  name            = "regionurlmap-test-%s"
  default_service = google_compute_region_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "blop"
  }

  host_rule {
    hosts        = ["myfavoritesite.com"]
    path_matcher = "blip"
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blop"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_advanced2(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  region        = "us-central1"
  name          = "regionurlmap-test-%s"
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  region   = "us-central1"
  name     = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}

resource "google_compute_region_url_map" "foobar" {
  region          = "us-central1"
  name            = "regionurlmap-test-%s"
  default_service = google_compute_region_backend_service.foobar.self_link

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
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blep"

    path_rule {
      paths   = ["/home"]
      service = google_compute_region_backend_service.foobar.self_link
    }

    path_rule {
      paths   = ["/login"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blub"

    path_rule {
      paths   = ["/*", "/blub"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "blip"

    path_rule {
      paths   = ["/*", "/home"]
      service = google_compute_region_backend_service.foobar.self_link
    }
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_noPathRules(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  region        = "us-central1"
  name          = "regionurlmap-test-%s"
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  health_checks = [google_compute_region_health_check.zero.self_link]
}

resource "google_compute_region_health_check" "zero" {
  region   = "us-central1"
  name     = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}

resource "google_compute_region_url_map" "foobar" {
  region          = "us-central1"
  name            = "regionurlmap-test-%s"
  default_service = google_compute_region_backend_service.foobar.self_link

  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }

  path_matcher {
    default_service = google_compute_region_backend_service.foobar.self_link
    name            = "boop"
  }

  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar.self_link
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_ilbPath(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name          = "regionurlmap-test-%s"
  description = "a description"
  default_service = google_compute_region_backend_service.home.self_link

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name = "allpaths"
    default_service = google_compute_region_backend_service.home.self_link

    path_rule {
      paths   = ["/home"]
      route_action {
        cors_policy {
          allow_credentials = true
          allow_headers = ["Allowed content"]
          allow_methods = ["GET"]
          allow_origins = ["Allowed origin"]
          expose_headers = ["Exposed header"]
          max_age = 30
          disabled = false
        }
        fault_injection_policy {
          abort {
            http_status = 234
            percentage = 5.6
          }
          delay {
            fixed_delay {
              seconds = 0
              nanos = 50000
            }
            percentage = 7.8
          }
        }
        request_mirror_policy {
          backend_service = google_compute_region_backend_service.home.self_link
        }
        retry_policy {
          num_retries = 4
          per_try_timeout {
            seconds = 30
          }
          retry_conditions = ["5xx", "deadline-exceeded"]
        }
        timeout {
          seconds = 20
          nanos = 750000000
        }
        url_rewrite {
          host_rewrite = "dev.example.com"
          path_prefix_rewrite = "/v1/api/"
        }
        weighted_backend_services {
          backend_service = google_compute_region_backend_service.home.self_link
          weight = 400
          header_action {
            request_headers_to_remove = ["RemoveMe"]
            request_headers_to_add {
              header_name = "AddMe"
              header_value = "MyValue"
              replace = true
            }
            response_headers_to_remove = ["RemoveMe"]
            response_headers_to_add {
              header_name = "AddMe"
              header_value = "MyValue"
              replace = false
            }
          }
        }
      }
    }
  }

  test {
    service = google_compute_region_backend_service.home.self_link
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_region_backend_service" "home" {
  name          = "regionurlmap-test-%s"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "default" {
  name          = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_ilbPathUpdate(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name          = "regionurlmap-test-%s"
  description = "a description"
  default_service = google_compute_region_backend_service.home2.self_link

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths2"
  }

  path_matcher {
    name = "allpaths2"
    default_service = google_compute_region_backend_service.home.self_link

    path_rule {
      paths   = ["/home2"]
      route_action {
        cors_policy {
          allow_credentials = true
          allow_headers = ["Allowed content again"]
          allow_methods = ["PUT"]
          allow_origins = ["Allowed origin again"]
          expose_headers = ["Exposed header again"]
          max_age = 31
          disabled = true
        }
        fault_injection_policy {
          abort {
            http_status = 345
            percentage = 6.7
          }
          delay {
            fixed_delay {
              seconds = 1
              nanos = 51000
            }
            percentage = 8.9
          }
        }
        request_mirror_policy {
          backend_service = google_compute_region_backend_service.home.self_link
        }
        retry_policy {
          num_retries = 6
          per_try_timeout {
            seconds = 31
          }
          retry_conditions = ["5xx"]
        }
        timeout {
          seconds = 21
          nanos = 760000000
        }
        url_rewrite {
          host_rewrite = "stage.example.com" # updated
          path_prefix_rewrite = "/v2/api/" # updated
        }
        weighted_backend_services {
          backend_service = google_compute_region_backend_service.home.self_link
          weight = 401
          header_action {
            request_headers_to_remove = ["RemoveMe2"]
            request_headers_to_add {
              header_name = "AddMe2"
              header_value = "MyValue2"
              replace = false
            }
            response_headers_to_remove = ["RemoveMe2"]
            response_headers_to_add {
              header_name = "AddMe2"
              header_value = "MyValue2"
              replace = true
            }
          }
        }
      }
    }
  }

  test {
    service = google_compute_region_backend_service.home.self_link
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_region_backend_service" "home" {
  name          = "regionurlmap-test-%s"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "home2" {
  name          = "regionurlmap-test-%s-2"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "default" {
  name          = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}
`, randomSuffix, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_ilbRoute(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name          = "regionurlmap-test-%s"
  description = "a description"
  default_service = google_compute_region_backend_service.home.self_link

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name = "allpaths"
    default_service = google_compute_region_backend_service.home.self_link

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
    service = google_compute_region_backend_service.home.self_link
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_region_backend_service" "home" {
  name          = "regionurlmap-test-%s"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "default" {
  name          = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_ilbRouteUpdate(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name          = "regionurlmap-test-%s"
  description = "a description"
  default_service = google_compute_region_backend_service.home.self_link

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths2"
  }

  path_matcher {
    name = "allpaths2"
    default_service = google_compute_region_backend_service.home2.self_link

    route_rules {
      priority = 2
      header_action {
        request_headers_to_remove = ["RemoveMe2Again"]
        request_headers_to_add {
          header_name = "AddSomethingElseAgain"
          header_value = "MyOtherValueAgain"
          replace = false
        }
        response_headers_to_remove = ["RemoveMe3Again"]
        response_headers_to_add {
          header_name = "AddMeAgain"
          header_value = "MyValueAgain"
          replace = true
        }
      }
      match_rules {
        full_path_match = "a full path again"
        header_matches {
          header_name = "someheaderagain"
          exact_match = "match this exactly again"
          invert_match = false
        }
        ignore_case = false
        metadata_filters {
          filter_match_criteria = "MATCH_ALL"
          filter_labels {
            name = "PLANET"
            value = "JUPITER"
          }
        }
      }
      url_redirect {
        host_redirect = "A hosti again"
        https_redirect = true
        path_redirect = "some/path/again"
        redirect_response_code = "TEMPORARY_REDIRECT"
        strip_query = false
      }
    }
  }

  test {
    service = google_compute_region_backend_service.home.self_link
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_region_backend_service" "home" {
  name          = "regionurlmap-test-%s"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "home2" {
  name          = "regionurlmap-test-%s-2"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "default" {
  name          = "regionurlmap-test-%s"
  http_health_check {
    port = 80
  }
}
`, randomSuffix, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_defaultUrlRedirectConfig(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_url_redirect {
    https_redirect = true
    strip_query    = false
  }
}
`, randomSuffix)
}

func testAccComputeRegionUrlMap_defaultUrlRedirectWithinPathMatcherConfig(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  name            = "urlmap-test-%s"
  default_url_redirect {
    https_redirect = true
    strip_query    = false
  }
  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name             = "allpaths"
    default_url_redirect {
      https_redirect = true
      strip_query    = false
    }
  }
}
`, randomSuffix)
}

func testAccComputeRegionUrlMap_defaultRouteAction_full(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  region = "us-central1"

  name        = "regionurlmap%s"
  description = "a description"

  default_route_action {

    retry_policy {
      num_retries = 3
      per_try_timeout {
        seconds = 0
        nanos = 500
      }
    }

    timeout {
      seconds = 3
      nanos = 0
    }

    url_rewrite {
      host_rewrite = "dev.example.com"
      path_prefix_rewrite = "/v1/api/"
    }

    request_mirror_policy {
      backend_service = google_compute_region_backend_service.login.id
    }

    cors_policy {
      allow_origins = [ "https://www.example.com" ]
      allow_methods = [ "GET" ]
      allow_headers = [ "Content-Type" ]
      expose_headers = [ "Authorization" ]
      max_age = 600
      allow_credentials = true
      disabled = false
    }

    weighted_backend_services {
      backend_service = google_compute_region_backend_service.login.id
      weight = 200
      header_action {
        request_headers_to_add {
          header_name = "foo-request-2"
          header_value = "bar"
          replace = true
        }
        request_headers_to_add {
          header_name = "foo-request-1"
          header_value = "bar"
          replace = true
        }
        request_headers_to_remove = [
          "fizz",
          "buzz"
        ]
        response_headers_to_add {
          header_name = "foo-response-2"
          header_value = "bar"
          replace = true
        }
        response_headers_to_add {
          header_name = "foo-response-1"
          header_value = "bar"
          replace = true
        }
        response_headers_to_remove = [
          "fizz",
          "buzz"
        ]
      }
    }
    weighted_backend_services {
      backend_service = google_compute_region_backend_service.home.id
      weight = 100
      header_action {
        request_headers_to_add {
          header_name = "foo-request-2"
          header_value = "bar"
          replace = true
        }
        request_headers_to_add {
          header_name = "foo-request-1"
          header_value = "bar"
          replace = true
        }
        request_headers_to_remove = [
          "fizz",
          "buzz"
        ]
        response_headers_to_add {
          header_name = "foo-response-2"
          header_value = "bar"
          replace = true
        }
        response_headers_to_add {
          header_name = "foo-response-1"
          header_value = "bar"
          replace = true
        }
        response_headers_to_remove = [
          "fizz",
          "buzz"
        ]
      }
    }
  }
}

resource "google_compute_region_backend_service" "login" {
  region = "us-central1"

  name        = "login%s"
  protocol    = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec = 10
}

resource "google_compute_region_backend_service" "home" {
  region = "us-central1"

  name        = "home%s"
  protocol    = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec = 10
}
`, randomSuffix, randomSuffix, randomSuffix)
}

func testAccComputeRegionUrlMap_defaultRouteAction_full_update(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_url_map" "foobar" {
  region = "us-central1"

  name        = "regionurlmap%s"
  description = "a description"

  default_route_action {

    # update all fields in retry_policy block
    retry_policy {
      num_retries = 4
      per_try_timeout {
        seconds = 1
        nanos = 0
      }
    }

    # update to be <1 second
    timeout {
      seconds = 0
      nanos = 10000000 # 0.01 seconds
    }

    # update both values
    url_rewrite {
      host_rewrite = "stage.example.com"
      path_prefix_rewrite = "/v2/api/"
    }

    # update backend_service field from 'login' to 'home'
    request_mirror_policy {
      backend_service = google_compute_region_backend_service.home.id 
    }

    # update policy and disable it
    cors_policy {
      allow_origins = [ "https://xylophone.example.com", "https://www.example.com" ]
      allow_methods = [ "PUT", "GET" ]
      allow_headers = [ "Content-Type" ]
      expose_headers = [ "Authorization" ]
      max_age = 600
      allow_credentials = true
      disabled = true
    }

    # Change various fields - marked with comments
    weighted_backend_services {
      backend_service = google_compute_region_backend_service.login.id
      weight = 150 # updated
      header_action {
        request_headers_to_add {
          header_name = "fizz-request-2" # updated
          header_value = "buzz" # updated
          replace = true
        }
        request_headers_to_add {
          header_name = "foo-request-1"
          header_value = "bar"
          replace = false # updated
        }
        request_headers_to_remove = [
          "fizz" # updated to remove element
        ]
        response_headers_to_add {
          header_name = "foo-response-2"
          header_value = "bar"
          replace = true
        }
        response_headers_to_add {
          header_name = "foo-response-1"
          header_value = "bar"
          replace = true
        }
        response_headers_to_remove = [
          "fizz",
          "buzz",
          "quack" # updated to add element
        ]
      }
    }
    weighted_backend_services {
      backend_service = google_compute_region_backend_service.home.id
      weight = 300 # updated
      header_action {
        request_headers_to_add {
          header_name = "foo-request-2"
          header_value = "bar"
          replace = true
        }
        # updated to remove a 'request_headers_to_add' block
        request_headers_to_remove = [
          "fizz",
          "buzz"
        ]
        response_headers_to_add {
          header_name = "foo-response-2"
          header_value = "bar"
          replace = true
        }
        response_headers_to_add {
          header_name = "foo-response-1"
          header_value = "bar"
          replace = true
        }
        # updated to add 'response_headers_to_add' block below
        response_headers_to_add {
          header_name = "foo-response-3"
          header_value = "bar"
          replace = true
        }
        response_headers_to_remove = [
          "fizz",
          "buzz"
        ]
      }
    }
  }
}

resource "google_compute_region_backend_service" "login" {
  region = "us-central1"

  name        = "login%s"
  protocol    = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec = 10
}

resource "google_compute_region_backend_service" "home" {
  region = "us-central1"

  name        = "home%s"
  protocol    = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec = 10
}
`, randomSuffix, randomSuffix, randomSuffix)
}
