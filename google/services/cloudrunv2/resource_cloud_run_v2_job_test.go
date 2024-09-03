// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrunv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudRunV2Job_cloudrunv2JobFullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobFull(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobFullUpdate(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "labels", "terraform_labels", "annotations", "deletion_protection"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobFull(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "tf-test-cloudrun-job%{random_suffix}"
    location = "us-central1"
    labels = {
      label-1 = "value-1"
    }
    annotations = {
      job-annotation-1 = "job-value-1"
    }
    client = "client-1"
    client_version = "client-version-1"
    
    template {
      labels = {
        label-1 = "value-1"
      }
      annotations = {
        temp-annotation-1 = "temp-value-1"
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

    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
  resource "google_service_account" "service_account" {
    account_id   = "tf-test-my-account%{random_suffix}"
    display_name = "Test Service Account"
  }
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobFullUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_job" "default" {
  name     = "tf-test-cloudrun-job%{random_suffix}"
  location = "us-central1"
  deletion_protection = false
  binary_authorization {
    use_default = true
    breakglass_justification = "Some justification"
  }
  labels = {
    label-1 = "value-update"
  }
  annotations = {
    job-annotation-1 = "job-value-update"
  }
  client = "client-update"
  client_version = "client-version-update"
  
  template {
    labels = {
      label-1 = "value-update"
    }
    annotations = {
      temp-annotation-1 = "temp-value-update"
    }
    parallelism = 2
    task_count = 8
    template {
      timeout = "500s"
      service_account = google_service_account.service_account.email
      execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
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
      max_retries = 0
    }
  }

  lifecycle {
    ignore_changes = [
      launch_stage,
    ]
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

func TestAccCloudRunV2Job_cloudrunv2JobWithDirectVPCUpdate(t *testing.T) {
	t.Parallel()

	jobName := fmt.Sprintf("tf-test-cloudrun-service%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"job_name": jobName,
		"project":  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithDirectVPC(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithDirectVPCAndNamedBinAuthPolicyUpdate(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobWithDirectVPC(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "%{job_name}"
    location = "us-central1"
    deletion_protection = false
    launch_stage = "BETA"
    template {
      template {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/job"
        }
        vpc_access {
          network_interfaces {
            network = "default"
          }
        }
      }
    }

    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobWithDirectVPCAndNamedBinAuthPolicyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "%{job_name}"
    location = "us-central1"
    deletion_protection = false
    launch_stage = "BETA"
    binary_authorization {
      policy = "projects/%{project}/platforms/cloudRun/policies/my-policy"
      breakglass_justification = "Some justification"
    }
    template {
      template {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/job"
        }
        vpc_access {
          network_interfaces {
            network = "my-network"
            subnetwork = "my-network"
            tags = ["tag1", "tag2", "tag3"]
          }
        }
      }
    }

    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
`, context)
}

func TestAccCloudRunV2Job_cloudrunv2JobWithGcsUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	jobName := fmt.Sprintf("tf-test-cloudrun-service%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"job_name": jobName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithNoVolume(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithGcsVolume(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobWithNoVolume(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "%{job_name}"
    location = "us-central1"
    deletion_protection = false
    template {
      template {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/job"
        }
      }
    }

    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
`, context)
}

func testAccCloudRunV2Job_cloudrunv2JobWithGcsVolume(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "%{job_name}"
    location = "us-central1"
    deletion_protection = false
    template {
      template {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/job"
          volume_mounts {
            name = "gcs"
            mount_path = "/mnt/gcs"
          }
        }
        volumes {
          name = "gcs"
          gcs {
            bucket = "gcp-public-data-landsat"
            read_only = true
          }
        }
      }
    }
    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
`, context)
}

func TestAccCloudRunV2Job_cloudrunv2JobWithNfsUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	jobName := fmt.Sprintf("tf-test-cloudrun-service%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"job_name": jobName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2JobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithNoVolume(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
			{
				Config: testAccCloudRunV2Job_cloudrunv2JobWithNfsVolume(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_job.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "launch_stage", "deletion_protection"},
			},
		},
	})
}

func testAccCloudRunV2Job_cloudrunv2JobWithNfsVolume(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_cloud_run_v2_job" "default" {
    name     = "%{job_name}"
    location = "us-central1"
    deletion_protection = false
    template {
      template {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/job"
          volume_mounts {
            name = "nfs"
            mount_path = "/mnt/nfs"
          }
        }
        volumes {
          name = "nfs"
          nfs {
            server = "10.0.10.10"
            path = "/"
            read_only = true
          }
        }
      }
    }
    lifecycle {
      ignore_changes = [
        launch_stage,
      ]
    }
  }
`, context)
}
