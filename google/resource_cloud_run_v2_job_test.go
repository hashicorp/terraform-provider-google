package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudRunV2Job_cloudrunv2JobFullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobFull(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobFullUpdate(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobFull(context map[string]interface{}) string {
	return Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "tf-test-cloudrun-job%{random_suffix}"
    location = "us-central1"
    launch_stage = "BETA"
    labels = {
      label-1 = "value-1"
    }
    client = "client-1"
    client_version = "client-version-1"
    
    template {
      labels = {
        label-1 = "value-1"
      }
      parallelism = 4
      task_count = 4
      template {
        timeout = "300s"
        service_account = google_service_account.service_account.email
        execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
        containers {
          name = "container-1"
          image = "us-docker.pkg.dev/cloudrun/container/hello"
          args = ["https://cloud.google.com/run", "www.google.com"]
          command = ["/bin/echo"]
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
            limits = {
              cpu = "4"
              memory = "2Gi"
            }
          }
        }
        max_retries = 5
      }
    }
  }
  resource "google_service_account" "service_account" {
    account_id   = "tf-test-my-account%{random_suffix}"
    display_name = "Test Service Account"
  }
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobFullUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
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
    parallelism = 2
    task_count = 8
    template {
      timeout = "500s"
      service_account = google_service_account.service_account.email
      execution_environment = "EXECUTION_ENVIRONMENT_GEN1"
      containers {
        name = "container-update"
        image = "us-docker.pkg.dev/cloudrun/container/hello"
        args = ["https://cloud.google.com/run"]
        command = ["printenv"]
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
      max_retries = 2
    }
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

func TestAccCloudRunV2Job_cloudrunv2JobProbesUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithEmptyTCPStartupProbeAndHTTPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobUpdateWithTCPStartupProbeAndHTTPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobUpdateWithHTTPStartupProbeAndTCPLivenessProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobUpdateWithEmptyHTTPStartupProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobUpdateWithHTTPStartupProbe(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobWithEmptyTCPStartupProbeAndHTTPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
  
  template {
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
}
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobUpdateWithTCPStartupProbeAndHTTPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
  
  template{
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
}
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobUpdateWithHTTPStartupProbeAndTCPLivenessProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
  
  template{
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
          initial_delay_seconds = 2
          period_seconds = 1
          timeout_seconds = 5
          failure_threshold = 2
          tcp_socket {
            port = 8080
          }
        }
      }
    }
  } 
}
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobUpdateWithEmptyHTTPStartupProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
  
  template {
    template {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
        startup_probe {
          http_get {}
        }
      }
    }
  }
}
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobUpdateWithHTTPStartupProbe(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  launch_stage = "BETA"
  
  template{
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
}
`, context)
}
