// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApphubWorkload_apphubWorkloadUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckApphubWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubWorkload_apphubWorkloadFullExample(context),
			},
			{
				ResourceName:            "google_apphub_workload.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id", "workload_id"},
			},
			{
				Config: testAccApphubWorkload_apphubWorkloadUpdate(context),
			},
			{
				ResourceName:            "google_apphub_workload.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id", "workload_id"},
			},
		},
	})
}

func testAccApphubWorkload_apphubWorkloadUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_application" "application" {
  location = "us-central1"
  application_id = "tf-test-example-application-1%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
}

resource "google_project" "service_project" {
  project_id ="tf-test-project-1%{random_suffix}"
  name = "Service Project"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
}

# Enable Compute API
resource "google_project_service" "compute_service_project" {
  project = google_project.service_project.project_id
  service = "compute.googleapis.com"
}

resource "time_sleep" "wait_120s" {
  depends_on = [google_project_service.compute_service_project]

  create_duration = "120s"
}

resource "google_apphub_service_project_attachment" "service_project_attachment" {
  service_project_attachment_id = google_project.service_project.project_id
  depends_on = [time_sleep.wait_120s]
}

# Discovered workload 
data "google_apphub_discovered_workload" "catalog-workload" {
  location = "us-central1"
  workload_uri = "${replace(google_compute_region_instance_group_manager.mig.instance_group, "https://www.googleapis.com/compute/v1", "//compute.googleapis.com")}"
  depends_on = [time_sleep.wait_120s_for_resource_ingestion]
}

resource "time_sleep" "wait_120s_for_resource_ingestion" {
  depends_on = [google_compute_region_instance_group_manager.mig]
  create_duration = "120s"
}

resource "google_apphub_workload" "example" {
  location = "us-central1"
  application_id = google_apphub_application.application.application_id
  workload_id = google_compute_region_instance_group_manager.mig.name
  discovered_workload = data.google_apphub_discovered_workload.catalog-workload.name
}

#Workload creation


# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  project                 = google_project.service_project.project_id
  auto_create_subnetworks = false
  depends_on = [time_sleep.wait_120s]
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  project       = google_project.service_project.project_id
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.ilb_network.id
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "tf-test-l7-ilb-mig-template%{random_suffix}"
  project      = google_project.service_project.project_id
  machine_type = "e2-small"
  tags         = ["http-server"]
  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id
    access_config {
      # add external ip to fetch packages
    }
  }
  disk {
    source_image = "debian-cloud/debian-12"
    auto_delete  = true
    boot         = true
  }
  # install nginx and serve a simple web page
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      set -euo pipefail
      export DEBIAN_FRONTEND=noninteractive
      apt-get update
      apt-get install -y nginx-light jq
      NAME=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/hostname")
      IP=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
      METADATA=$(curl -f -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/attributes/?recursive=True" | jq 'del(.["startup-script"])')
      cat <<EOF > /var/www/html/index.html
      <pre>
      Name: $NAME
      IP: $IP
      Metadata: $METADATA
      </pre>
      EOF
    EOF1
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_region_instance_group_manager" "mig" {
  name     = "tf-test-l7-ilb-mig1%{random_suffix}"
  project  = google_project.service_project.project_id
  region   = "us-central1"
  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
  base_instance_name = "vm"
  target_size        = 2
}
`, context)
}
