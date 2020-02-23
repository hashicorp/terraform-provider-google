package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	extraCheckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_basic(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_basicModified(
					serviceName, checkName, extraCheckName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withBackend(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 20),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withBackendAndIAP(t *testing.T) {
	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withBackendAndIAP(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:            "google_compute_backend_service.lipsum",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"iap.0.oauth2_client_secret"},
			},
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_updatePreservesOptionalParameters(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "initial-description", "GENERATED_COOKIE"),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "updated-description", "GENERATED_COOKIE"),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withConnectionDraining(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withConnectionDraining(serviceName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withConnectionDrainingAndUpdate(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withConnectionDraining(serviceName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_basic(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withHttpsHealthCheck(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withHttpsHealthCheck(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withCdnPolicy(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withCdnPolicy(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withSecurityPolicy(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	polName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSecurityPolicy(serviceName, checkName, polName, "google_compute_security_policy.policy.self_link"),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withSecurityPolicy(serviceName, checkName, polName, "\"\""),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withCDNEnabled(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withCDNEnabled(
					serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withSessionAffinity(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "description", "CLIENT_IP"),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "description", "GENERATED_COOKIE"),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withAffinityCookieTtlSec(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withAffinityCookieTtlSec(
					serviceName, checkName, "description", "GENERATED_COOKIE", 300),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withMaxConnections(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxConnections(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withMaxConnections(
					serviceName, igName, itName, checkName, 20),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withMaxConnectionsPerInstance(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerInstance(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerInstance(
					serviceName, igName, itName, checkName, 20),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withMaxRatePerEndpoint(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(10)
	service := fmt.Sprintf("tf-test-%s", randSuffix)
	instance := fmt.Sprintf("tf-test-%s", randSuffix)
	neg := fmt.Sprintf("tf-test-%s", randSuffix)
	network := fmt.Sprintf("tf-test-%s", randSuffix)
	check := fmt.Sprintf("tf-test-%s", randSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxRatePerEndpoint(
					service, instance, neg, network, check, 0.2),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withMaxRatePerEndpoint(
					service, instance, neg, network, check, 0.4),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withMaxConnectionsPerEndpoint(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(10)
	service := fmt.Sprintf("tf-test-%s", randSuffix)
	instance := fmt.Sprintf("tf-test-%s", randSuffix)
	neg := fmt.Sprintf("tf-test-%s", randSuffix)
	network := fmt.Sprintf("tf-test-%s", randSuffix)
	check := fmt.Sprintf("tf-test-%s", randSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerEndpoint(
					service, instance, neg, network, check, 5),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerEndpoint(
					service, instance, neg, network, check, 10),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeBackendService_basic(serviceName, checkName string) string {
	return fmt.Sprintf(`
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
`, serviceName, checkName)
}

func testAccComputeBackendService_withCDNEnabled(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
  enable_cdn    = true
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_basicModified(serviceName, checkOne, checkTwo string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.one.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_http_health_check" "one" {
  name               = "%s"
  request_path       = "/one"
  check_interval_sec = 30
  timeout_sec        = 30
}
`, serviceName, checkOne, checkTwo)
}

func testAccComputeBackendService_withBackend(
	serviceName, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = %v

  backend {
    group = google_compute_instance_group_manager.foobar.instance_group
  }

  health_checks = [google_compute_http_health_check.default.self_link]
}

resource "google_compute_instance_group_manager" "foobar" {
  name = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_http_health_check" "default" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, timeout, igName, itName, checkName)
}

func testAccComputeBackendService_withBackendAndIAP(
	serviceName, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = %v

  backend {
    group = google_compute_instance_group_manager.foobar.instance_group
  }

  iap {
    oauth2_client_id     = "test"
    oauth2_client_secret = "test"
  }

  health_checks = [google_compute_http_health_check.default.self_link]
}

resource "google_compute_instance_group_manager" "foobar" {
  name = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_http_health_check" "default" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, timeout, igName, itName, checkName)
}

func testAccComputeBackendService_withSessionAffinity(serviceName, checkName, description, affinityName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name             = "%s"
  description      = "%s"
  health_checks    = [google_compute_http_health_check.zero.self_link]
  session_affinity = "%s"
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, description, affinityName, checkName)
}

func testAccComputeBackendService_withAffinityCookieTtlSec(serviceName, checkName, description, affinityName string, affinityCookieTtlSec int64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name                    = "%s"
  description             = "%s"
  health_checks           = [google_compute_http_health_check.zero.self_link]
  session_affinity        = "%s"
  affinity_cookie_ttl_sec = %v
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, description, affinityName, affinityCookieTtlSec, checkName)
}

func testAccComputeBackendService_withConnectionDraining(serviceName, checkName string, drainingTimeout int64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name                            = "%s"
  health_checks                   = [google_compute_http_health_check.zero.self_link]
  connection_draining_timeout_sec = %v
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, drainingTimeout, checkName)
}

func testAccComputeBackendService_withHttpsHealthCheck(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_https_health_check.zero.self_link]
  protocol      = "HTTPS"
}

resource "google_compute_https_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_withCdnPolicy(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]

  cdn_policy {
    cache_key_policy {
      include_protocol       = true
      include_host           = true
      include_query_string   = true
      query_string_whitelist = ["foo", "bar"]
    }
  }
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_withSecurityPolicy(serviceName, checkName, polName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name            = "%s"
  health_checks   = [google_compute_http_health_check.zero.self_link]
  security_policy = %s
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
}
`, serviceName, polLink, checkName, polName)
}

func testAccComputeBackendService_withMaxConnections(
	serviceName, igName, itName, checkName string, maxConnections int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"

  backend {
    group           = google_compute_instance_group_manager.foobar.instance_group
    max_connections = %v
  }

  health_checks = [google_compute_health_check.default.self_link]
}

resource "google_compute_instance_group_manager" "foobar" {
  name = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name = "%s"
  tcp_health_check {
    port = "110"
  }
}
`, serviceName, maxConnections, igName, itName, checkName)
}

func testAccComputeBackendService_withMaxConnectionsPerInstance(
	serviceName, igName, itName, checkName string, maxConnectionsPerInstance int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"

  backend {
    group                        = google_compute_instance_group_manager.foobar.instance_group
    max_connections_per_instance = %v
  }

  health_checks = [google_compute_health_check.default.self_link]
}

resource "google_compute_instance_group_manager" "foobar" {
  name = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name = "%s"
  tcp_health_check {
    port = "110"
  }
}
`, serviceName, maxConnectionsPerInstance, igName, itName, checkName)
}

func testAccComputeBackendService_withMaxConnectionsPerEndpoint(
	service, instance, neg, network, check string, maxConnections int64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"

  backend {
    group                        = google_compute_network_endpoint_group.lb-neg.self_link
    balancing_mode               = "CONNECTION"
    max_connections_per_endpoint = %v
  }

  health_checks = [google_compute_health_check.default.self_link]
}

data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "endpoint-instance" {
  name         = "%s"
  machine_type = "n1-standard-1"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.default.self_link
    access_config {
      network_tier = "PREMIUM"
    }
  }
}

resource "google_compute_network_endpoint_group" "lb-neg" {
  name         = "%s"
  network      = google_compute_network.default.self_link
  subnetwork   = google_compute_subnetwork.default.self_link
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network_endpoint" "lb-endpoint" {
  network_endpoint_group = google_compute_network_endpoint_group.lb-neg.name

  instance   = google_compute_instance.endpoint-instance.name
  port       = google_compute_network_endpoint_group.lb-neg.default_port
  ip_address = google_compute_instance.endpoint-instance.network_interface[0].network_ip
}

resource "google_compute_network" "default" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_health_check" "default" {
  name = "%s"
  tcp_health_check {
    port = "110"
  }
}
`, service, maxConnections, instance, neg, network, network, check)
}

func testAccComputeBackendService_withMaxRatePerEndpoint(
	service, instance, neg, network, check string, maxRate float64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "https"
  protocol    = "HTTPS"

  backend {
    group                 = google_compute_network_endpoint_group.lb-neg.self_link
    balancing_mode        = "RATE"
    max_rate_per_endpoint = %v
  }

  health_checks = [google_compute_health_check.default.self_link]
}

data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "endpoint-instance" {
  name         = "%s"
  machine_type = "n1-standard-1"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.default.self_link
    access_config {
      network_tier = "PREMIUM"
    }
  }
}

resource "google_compute_network_endpoint_group" "lb-neg" {
  name         = "%s"
  network      = google_compute_network.default.self_link
  subnetwork   = google_compute_subnetwork.default.self_link
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network_endpoint" "lb-endpoint" {
  network_endpoint_group = google_compute_network_endpoint_group.lb-neg.name

  instance   = google_compute_instance.endpoint-instance.name
  port       = google_compute_network_endpoint_group.lb-neg.default_port
  ip_address = google_compute_instance.endpoint-instance.network_interface[0].network_ip
}

resource "google_compute_network" "default" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_health_check" "default" {
  name                = "%s"
  check_interval_sec  = 3
  healthy_threshold   = 3
  timeout_sec         = 2
  unhealthy_threshold = 3
  https_health_check {
    port = "443"
  }
}
`, service, maxRate, instance, neg, network, network, check)
}
