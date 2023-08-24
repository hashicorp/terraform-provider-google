// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrunv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/cloudrunv2"
)

func TestAccCloudRunV2Service_cloudrunv2ServiceFullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
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
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  description = "description creating"
  location = "us-central1"
  annotations = {
    generated-by = "magic-modules"
  }
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
    annotations = {
      generated-by = "magic-modules"
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
        startup_cpu_boost = true
        limits = {
          cpu = "4"
          memory = "2Gi"
        }
      }
    }
    session_affinity = false
  }
}

resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Test Service Account"
}
`, context)
}

func testAccCloudRunV2Service_cloudrunv2ServiceFullUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  description = "description updating"
  location = "us-central1"
  annotations = {
    generated-by = "magic-modules-files"
  }
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
    annotations = {
      generated-by = "magic-modules"
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
        startup_cpu_boost = false
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
    session_affinity = true
  }
  traffic {
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
    tag = "tt-update"
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

func TestAccCloudRunV2Service_cloudrunv2ServiceTCPProbesUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
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
		},
	})
}

func TestAccCloudRunV2Service_cloudrunv2ServiceHTTPProbesUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
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

func TestAccCloudRunV2Service_cloudrunv2ServiceGRPCProbesUpdate(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-cloudrun-service%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"service_name": serviceName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudRunServiceUpdateWithEmptyGRPCLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			// The following test steps of gRPC startup probe are expected to fail with startup probe check failures.
			// This is because, due to the unavailability of ready-to-use container images of a gRPC service that
			// implements the standard gRPC health check protocol, we compromise and use a container image of an
			// ordinary HTTP service to deploy the gRPC service, which never passes startup probes.
			// So we only check that the `startup.grpc {}` block and its properties are accepted by the APIs.
			{
				Config:      testAccCloudRunV2Service_cloudRunServiceUpdateWithEmptyGRPCStartupProbe(context),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Revision '%s-.*' is not ready and cannot serve traffic\. The user-provided container failed the configured startup probe checks\.`, serviceName)),
			},
			{
				PreConfig:   testAccCheckCloudRunV2ServiceDestroyByNameProducer(t, serviceName),
				Config:      testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCStartupProbe(context),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Revision '%s-.*' is not ready and cannot serve traffic\. The user-provided container failed the configured startup probe checks\.`, serviceName)),
			},
			{
				PreConfig:   testAccCheckCloudRunV2ServiceDestroyByNameProducer(t, serviceName),
				Config:      testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCLivenessAndStartupProbes(context),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Revision '%s-.*' is not ready and cannot serve traffic\. The user-provided container failed the configured startup probe checks\.`, serviceName)),
			},
			{
				PreConfig:          testAccCheckCloudRunV2ServiceDestroyByNameProducer(t, serviceName),
				Config:             testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCLivenessAndStartupProbes(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCloudRunV2ServiceDestroyByNameProducer(t *testing.T, serviceName string) func() {
	return func() {
		config := acctest.GoogleProviderConfig(t)
		service := config.NewCloudRunV2Client(config.UserAgent).Projects.Locations.Services
		qualifiedServiceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", config.Project, config.Region, serviceName)
		op, err := service.Delete(qualifiedServiceName).Do()
		if err != nil {
			t.Errorf("Error while deleting the Cloud Run service: %s", err)
			return
		}
		err = cloudrunv2.RunAdminV2OperationWaitTime(config, op, config.Project, "Waiting for Cloud Run service to be deleted", config.UserAgent, 5*time.Minute)
		if err != nil {
			t.Errorf("Error while waiting for Cloud Run service delete operation to complete: %s", err.Error())
		}
	}
}

func testAccCloudRunV2Service_cloudrunv2ServiceWithEmptyTCPStartupProbeAndHTTPLivenessProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
          port = 8080
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

func testAccCloudRunV2Service_cloudrunv2ServiceUpdateWithEmptyHTTPStartupProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      startup_probe {
        initial_delay_seconds = 3
        period_seconds = 2
        timeout_seconds = 6
        failure_threshold = 3
        http_get {
          path = "/some-path"
          port = 8080
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

func testAccCloudRunV2Service_cloudRunServiceUpdateWithEmptyGRPCLivenessProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     ="%{service_name}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      liveness_probe {
        grpc {}
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCLivenessProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "%{service_name}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      liveness_probe {
        grpc {
          port = 8080
          service = "grpc.health.v1.Health"
        }
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudRunServiceUpdateWithEmptyGRPCStartupProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "%{service_name}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      startup_probe {
        grpc {}
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCStartupProbe(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "%{service_name}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      startup_probe {
        grpc {
          port = 8080
          service = "grpc.health.v1.Health"
        }
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Service_cloudRunServiceUpdateWithGRPCLivenessAndStartupProbes(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "%{service_name}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8080
      }
      liveness_probe {
        grpc {
          port = 8080
          service = "grpc.health.v1.Health"
        }
      }
      startup_probe {
        grpc {
          port = 8080
          service = "grpc.health.v1.Health"
        }
      }
    }
  }
}
`, context)
}
