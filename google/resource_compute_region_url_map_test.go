package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeRegionUrlMap_update_path_matcher(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
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

	randomSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
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

	randomSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
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

func testAccComputeRegionUrlMap_basic1(randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
	region        = "us-central1"
	name          = "regionurlmap-test-%s"
	protocol      = "HTTP"
	health_checks = ["${google_compute_region_health_check.zero.self_link}"]
}

resource "google_compute_region_health_check" "zero" {
	region = "us-central1"
	name   = "regionurlmap-test-%s"
	http_health_check {
	}
}

resource "google_compute_region_url_map" "foobar" {
	region          = "us-central1"
	name            = "regionurlmap-test-%s"
	default_service = "${google_compute_region_backend_service.foobar.self_link}"

	host_rule {
		hosts        = ["mysite.com", "myothersite.com"]
		path_matcher = "boop"
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "boop"

		path_rule {
			paths   = ["/*"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}
	}

	test {
		host    = "mysite.com"
		path    = "/*"
		service = "${google_compute_region_backend_service.foobar.self_link}"
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
	health_checks = ["${google_compute_region_health_check.zero.self_link}"]
}

resource "google_compute_region_health_check" "zero" {
	region = "us-central1"
	name   = "regionurlmap-test-%s"
	http_health_check {
	}
}

resource "google_compute_region_url_map" "foobar" {
	region          = "us-central1"
	name            = "regionurlmap-test-%s"
	default_service = "${google_compute_region_backend_service.foobar.self_link}"

	host_rule {
		hosts        = ["mysite.com", "myothersite.com"]
		path_matcher = "blip"
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blip"

		path_rule {
			paths   = ["/*", "/home"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}
	}

	test {
		host    = "mysite.com"
		path    = "/test"
		service = "${google_compute_region_backend_service.foobar.self_link}"
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
	health_checks = ["${google_compute_region_health_check.zero.self_link}"]
}

resource "google_compute_region_health_check" "zero" {
	region = "us-central1"
	name   = "regionurlmap-test-%s"
	http_health_check {
	}
}

resource "google_compute_region_url_map" "foobar" {
	region          = "us-central1"
	name            = "regionurlmap-test-%s"
	default_service = "${google_compute_region_backend_service.foobar.self_link}"

	host_rule {
		hosts        = ["mysite.com", "myothersite.com"]
		path_matcher = "blop"
	}

	host_rule {
		hosts        = ["myfavoritesite.com"]
		path_matcher = "blip"
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blop"

		path_rule {
			paths   = ["/*", "/home"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blip"

		path_rule {
			paths   = ["/*", "/home"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
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
	health_checks = ["${google_compute_region_health_check.zero.self_link}"]
}

resource "google_compute_region_health_check" "zero" {
	region = "us-central1"
	name = "regionurlmap-test-%s"
	http_health_check {
	}
}

resource "google_compute_region_url_map" "foobar" {
	region          = "us-central1"
	name            = "regionurlmap-test-%s"
	default_service = "${google_compute_region_backend_service.foobar.self_link}"

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
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blep"

		path_rule {
			paths   = ["/home"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}

		path_rule {
			paths   = ["/login"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blub"

		path_rule {
			paths   = ["/*", "/blub"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
		}
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "blip"

		path_rule {
			paths   = ["/*", "/home"]
			service = "${google_compute_region_backend_service.foobar.self_link}"
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
	health_checks = ["${google_compute_region_health_check.zero.self_link}"]
}

resource "google_compute_region_health_check" "zero" {
	region = "us-central1"
	name   = "regionurlmap-test-%s"
	http_health_check {
	}
}

resource "google_compute_region_url_map" "foobar" {
	region          = "us-central1"
	name            = "regionurlmap-test-%s"
	default_service = "${google_compute_region_backend_service.foobar.self_link}"

	host_rule {
		hosts        = ["mysite.com", "myothersite.com"]
		path_matcher = "boop"
	}

	path_matcher {
		default_service = "${google_compute_region_backend_service.foobar.self_link}"
		name            = "boop"
	}

	test {
		host    = "mysite.com"
		path    = "/*"
		service = "${google_compute_region_backend_service.foobar.self_link}"
	}
}
`, randomSuffix, randomSuffix, randomSuffix)
}
