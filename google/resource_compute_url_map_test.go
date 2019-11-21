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
