package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudRunV2Service_cloudrunv2ServiceFullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceFull(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceFullUpdate(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceFull(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  description = "description creating"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"
  labels = {
    label-1 = "value-1"
  }
  client = "client-1"
  client_version = "client-version-1"
  
  template {
    labels = {
      label-1 = "value-1"
    }
    timeout = "300s"
    service_account = google_service_account.service_account.email
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
    scaling {
      max_instance_count = 3
      min_instance_count = 1
    }
    containers {
      name = "container-1"
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      env {
        name = "SOURCE"
        value = "remote"
      }
      env {
        name = "TARGET"
        value = "home"
      }
      ports {
        name = "h2c"
        container_port = 8080
      }
      resources {
        cpu_idle = true
        limits = {
          cpu = "4"
          memory = "2Gi"
        }
      }
    }
  }
  traffic {
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    tag = "traffic-tag-1"
  }
}

resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Test Service Account"
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceFullUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  description = "description updating"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"
  binary_authorization {
    use_default = true
    breakglass_justification = "Some justification"
  }
  labels = {
    label-1 = "value-update"
  }
  client = "client-update"
  client_version = "client-version-update"
  
  template {
    labels = {
      label-1 = "value-update"
    }
    timeout = "500s"
    service_account = google_service_account.service_account.email
    execution_environment = "EXECUTION_ENVIRONMENT_GEN1"
    scaling {
      max_instance_count = 2
      min_instance_count = 1
    }
    containers {
      name = "container-update"
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      env {
        name = "SOURCE_UPDATE"
        value = "remote-update"
      }
      env {
        name = "TARGET_UPDATE"
        value = "home-update"
      }
      ports {
        name = "h2c"
        container_port = 8080
      }
      resources {
        cpu_idle = true
        limits = {
          cpu = "2"
          memory = "8Gi"
        }
      }
    }
    vpc_access{
      connector = google_vpc_access_connector.connector.id
      egress = "ALL_TRAFFIC"
    }
  }
  traffic {
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
    tag = "traffic-tag-update"
  }
}

resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Test Service Account"
}

resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-run-vpc%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  region        = "us-central1"
}
resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-run-subnetwork%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}
resource "google_compute_network" "custom_test" {
  name                    = "tf-test-run-network%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceProbesUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceWithEmptyTCPStartupProbeAndHTTPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithTCPStartupProbeAndHTTPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithHTTPStartupProbeAndTCPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithEmptyHTTPStartupProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithHTTPStartupProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceWithEmptyTCPStartupProbeAndHTTPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      startup_probe {
        tcp_socket {}
      }
      liveness_probe {
        http_get {}
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithTCPStartupProbeAndHTTPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      startup_probe {
        initial_delay_seconds = 2
        period_seconds = 1
        timeout_seconds = 5
        failure_threshold = 2
        tcp_socket {
          port = 8080
        }
      }
      liveness_probe {
        initial_delay_seconds = 2
        period_seconds = 1
        timeout_seconds = 5
        failure_threshold = 2
        http_get {
          path = "/some-path"
          http_headers {
            name = "User-Agent"
            value = "magic-modules"
          }
          http_headers {
            name = "Some-Name"
          }
        }
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithHTTPStartupProbeAndTCPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      startup_probe {
        initial_delay_seconds = 3
        period_seconds = 2
        timeout_seconds = 6
        failure_threshold = 3
        http_get {
          path = "/some-path"
          http_headers {
            name = "User-Agent"
            value = "magic-modules"
          }
          http_headers {
            name = "Some-Name"
          }
        }
      }
      liveness_probe {
        initial_delay_seconds = 3
        period_seconds = 2
        timeout_seconds = 6
        failure_threshold = 3
        tcp_socket {
          port = 8080
        }
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithEmptyHTTPStartupProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      startup_probe {
        http_get {}
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithHTTPStartupProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      startup_probe {
        http_get {
          path = "/some-path"
          http_headers {
            name = "User-Agent"
            value = "magic-modules"
          }
          http_headers {
            name = "Some-Name"
          }
        }
      }
    }
  }
}
`, context)
}
