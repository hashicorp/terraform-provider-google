package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRegionBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	extraCheckName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionBackendService_basic(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionBackendService_basicModified(
					serviceName, checkName, extraCheckName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionBackendService_withBackendInternal(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	igName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	itName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionBackendService_withInvalidInternalBackend(
					serviceName, igName, itName, checkName),
				ExpectError: regexp.MustCompile(`capacity_scaler" cannot be set for INTERNAL backend service`),
			},
			{
				Config: testAccComputeRegionBackendService_withBackend(
					serviceName, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_region_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionBackendService_withBackend(
					serviceName, igName, itName, checkName, 20),
			},
			{
				ResourceName:      "google_compute_region_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionBackendService_withBackendInternalManaged(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	igmName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	hcName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeRegionBackendService_internalManagedInvalidBackend(serviceName, igmName, hcName),
				ExpectError: regexp.MustCompile("non-zero capacity_scaler"),
			},
			{
				Config:      testAccComputeRegionBackendService_internalManagedInvalidMultiBackends(serviceName, igmName, hcName),
				ExpectError: regexp.MustCompile("non-zero capacity_scaler"),
			},
			{
				Config: testAccComputeRegionBackendService_internalManagedMultipleBackends(serviceName, igmName, hcName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionBackendService_withBackendMultiNic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	net1Name := fmt.Sprintf("tf-test-%s", randString(t, 10))
	net2Name := fmt.Sprintf("tf-test-%s", randString(t, 10))
	igName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	itName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionBackendService_withBackendMultiNic(
					serviceName, net1Name, net2Name, igName, itName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_region_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionBackendService_withConnectionDrainingAndUpdate(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionBackendService_withConnectionDraining(serviceName, checkName, 10),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionBackendService_basic(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionBackendService_ilbUpdateBasic(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("foo-%s", randString(t, 10))
	checkName := fmt.Sprintf("bar-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionBackendService_ilbBasic(backendName, checkName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionBackendService_ilbUpdateBasic(backendName, checkName),
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionBackendService_ilbBasic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name                  = "%s"
  health_checks         = [google_compute_health_check.health_check.self_link]
  port_name             = "http"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  locality_lb_policy    = "RING_HASH"
  circuit_breakers {
    max_connections = 10
  }
  consistent_hash {
    http_cookie {
      ttl {
        seconds = 11
        nanos   = 1234
      }
      name = "mycookie"
    }
  }
  outlier_detection {
    consecutive_errors = 2
  }
}

resource "google_compute_health_check" "health_check" {
  name     = "%s"
  http_health_check {
    port = 80
  }
}
`, serviceName, checkName)
}

func testAccComputeRegionBackendService_ilbUpdateBasic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name                  = "%s"
  health_checks         = [google_compute_health_check.health_check.self_link]
  port_name             = "https"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  locality_lb_policy    = "RANDOM"
  circuit_breakers {
    max_connections = 10
  }
  outlier_detection {
    consecutive_errors = 2
  }
}

resource "google_compute_health_check" "health_check" {
  name     = "%s"
  http_health_check {
    port = 80
  }
}
`, serviceName, checkName)
}

func testAccComputeRegionBackendService_basic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_health_check.zero.self_link]
  region        = "us-central1"
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
`, serviceName, checkName)
}

func testAccComputeRegionBackendService_basicModified(serviceName, checkOne, checkTwo string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name          = "%s"
  health_checks = [google_compute_health_check.one.self_link]
  region        = "us-central1"
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = 443
  }
}

resource "google_compute_health_check" "one" {
  name               = "%s"
  check_interval_sec = 30
  timeout_sec        = 30

  tcp_health_check {
    port = 443
  }
}
`, serviceName, checkOne, checkTwo)
}

func testAccComputeRegionBackendService_withBackend(
	serviceName, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_region_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  protocol    = "TCP"
  region      = "us-central1"
  timeout_sec = %v

  backend {
    group    = google_compute_instance_group_manager.foobar.instance_group
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
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = 443
  }
}
`, serviceName, timeout, igName, itName, checkName)
}

func testAccComputeRegionBackendService_withBackendMultiNic(
	serviceName, net1Name, net2Name, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_region_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  protocol    = "TCP"
  region      = "us-central1"
  timeout_sec = %v

  backend {
    group = google_compute_instance_group_manager.foobar.instance_group
  }

  network = google_compute_network.network2.self_link

  health_checks = [google_compute_health_check.default.self_link]
}

resource "google_compute_network" "network1" {
  name                            = "%s"
  auto_create_subnetworks         = false
}

resource "google_compute_subnetwork" "subnet1" {
  name                     = "%s"
  ip_cidr_range            = "10.0.1.0/24"
  region                   = "us-central1"
  private_ip_google_access = true
  network                  = google_compute_network.network1.self_link
}

resource "google_compute_network" "network2" {
  name                            = "%s"
  auto_create_subnetworks         = false
}

resource "google_compute_subnetwork" "subnet2" {
  name                     = "%s"
  ip_cidr_range            = "10.0.2.0/24"
  region                   = "us-central1"
  private_ip_google_access = true
  network                  = google_compute_network.network2.self_link
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "%s"
  version {
    instance_template  = google_compute_instance_template.foobar.self_link
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    subnetwork = google_compute_subnetwork.subnet1.self_link
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnet2.self_link
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = 443
  }
}
`, serviceName, timeout, net1Name, net1Name, net2Name, net2Name, igName, itName, checkName)
}

func testAccComputeRegionBackendService_withInvalidInternalBackend(
	serviceName, igName, itName, checkName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_region_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"
  region      = "us-central1"

  backend {
    group    = google_compute_instance_group_manager.foobar.instance_group
    capacity_scaler = 1.0
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
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = 443
  }
}
`, serviceName, igName, itName, checkName)
}

func testAccComputeRegionBackendService_internalManagedInvalidBackend(serviceName, igmName, hcName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "default" {
  name        = "%s"

  load_balancing_scheme = "INTERNAL_MANAGED"

  backend {
    group          = google_compute_region_instance_group_manager.rigm.instance_group
    balancing_mode = "UTILIZATION"
  }

  region      = "us-central1"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
}

data "google_compute_image" "debian_image" {
  family   = "debian-9"
  project  = "debian-cloud"
}

resource "google_compute_region_instance_group_manager" "rigm" {
  region   = "us-central1"
  name     = "%s"
  version {
    instance_template = google_compute_instance_template.instance_template.self_link
    name              = "primary"
  }
  base_instance_name = "internal-glb"
  target_size        = 1
}

resource "google_compute_instance_template" "instance_template" {
  name         = "%s-template"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = data.google_compute_image.debian_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_region_health_check" "default" {
  name   = "%s"
  region = "us-central1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}
`, serviceName, igmName, igmName, hcName)
}

func testAccComputeRegionBackendService_internalManagedInvalidMultiBackends(serviceName, igmName, hcName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "default" {
  name        = "%s"
  load_balancing_scheme = "INTERNAL_MANAGED"

  backend {
    group           = google_compute_region_instance_group_manager.rigm1.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 0.0
  }

  backend {
    group           = google_compute_region_instance_group_manager.rigm2.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 0.0
  }

  region      = "us-central1"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
}

data "google_compute_image" "debian_image" {
  family   = "debian-9"
  project  = "debian-cloud"
}

resource "google_compute_region_instance_group_manager" "rigm1" {
  name     = "%s-1"
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.self_link
    name              = "primary"
  }
  base_instance_name = "internal-glb"
  target_size        = 1
}

resource "google_compute_region_instance_group_manager" "rigm2" {
  name     = "%s-2"
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.self_link
    name              = "primary"
  }
  base_instance_name = "internal-glb"
  target_size        = 1
}

resource "google_compute_instance_template" "instance_template" {
  name         = "%s-template"
  machine_type = "n1-standard-1"

  network_interface {
    network    = "default"
  }

  disk {
    source_image = data.google_compute_image.debian_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_region_health_check" "default" {
  name   = "%s"
  region = "us-central1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}
`, serviceName, igmName, igmName, igmName, hcName)
}

func testAccComputeRegionBackendService_internalManagedMultipleBackends(serviceName, igmName, hcName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "default" {
  name        = "%s"
  load_balancing_scheme = "INTERNAL_MANAGED"

  backend {
    group          = google_compute_region_instance_group_manager.rigm1.instance_group
    balancing_mode = "UTILIZATION"
  }

  backend {
    group          = google_compute_region_instance_group_manager.rigm2.instance_group
    balancing_mode = "UTILIZATION"
    capacity_scaler = 1.0
  }

  region      = "us-central1"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_region_health_check.default.self_link]
}

data "google_compute_image" "debian_image" {
  family   = "debian-9"
  project  = "debian-cloud"
}

resource "google_compute_region_instance_group_manager" "rigm1" {
  name     = "%s-1"
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.self_link
    name              = "primary"
  }
  base_instance_name = "internal-glb"
  target_size        = 1
}

resource "google_compute_region_instance_group_manager" "rigm2" {
  name     = "%s-2"
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.self_link
    name              = "primary"
  }
  base_instance_name = "internal-glb"
  target_size        = 1
}

resource "google_compute_instance_template" "instance_template" {
  name         = "%s-template"
  machine_type = "n1-standard-1"

  network_interface {
    network    = "default"
  }

  disk {
    source_image = data.google_compute_image.debian_image.self_link
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_region_health_check" "default" {
  name   = "%s"
  region = "us-central1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}
`, serviceName, igmName, igmName, igmName, hcName)
}

func testAccComputeRegionBackendService_withConnectionDraining(serviceName, checkName string, drainingTimeout int64) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar" {
  name                            = "%s"
  health_checks                   = [google_compute_health_check.zero.self_link]
  region                          = "us-central1"
  connection_draining_timeout_sec = %v
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
`, serviceName, drainingTimeout, checkName)
}
