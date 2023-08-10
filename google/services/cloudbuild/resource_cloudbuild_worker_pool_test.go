// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudbuild_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCloudbuildWorkerPool_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             funcAccTestCloudbuildWorkerPoolCheckDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildWorkerPool_basic(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_cloudbuild_worker_pool.pool",
			},
			{
				Config: testAccCloudbuildWorkerPool_updated(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_cloudbuild_worker_pool.pool",
			},
			{
				Config: testAccCloudbuildWorkerPool_noWorkerConfig(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_cloudbuild_worker_pool.pool",
			},
		},
	})
}

func testAccCloudbuildWorkerPool_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuild_worker_pool" "pool" {
	name = "pool%{random_suffix}"
	location = "europe-west1"
	worker_config {
		disk_size_gb = 100
		machine_type = "e2-standard-8"
		no_external_ip = true
	}
}
`, context)
}

func testAccCloudbuildWorkerPool_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuild_worker_pool" "pool" {
	name = "pool%{random_suffix}"
	location = "europe-west1"
	worker_config {
		disk_size_gb = 101
		machine_type = "e2-standard-4"
		no_external_ip = false
	}
}
`, context)
}

func testAccCloudbuildWorkerPool_noWorkerConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuild_worker_pool" "pool" {
	name = "pool%{random_suffix}"
	location = "europe-west1"
}
`, context)
}

func TestAccCloudbuildWorkerPool_withNetwork(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "cloudbuild-workerpool"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             funcAccTestCloudbuildWorkerPoolCheckDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildWorkerPool_withNetwork(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_cloudbuild_worker_pool.pool",
			},
		},
	})
}

func testAccCloudbuildWorkerPool_withNetwork(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_compute_network" "network" {
  name = "%{network_name}"
}

resource "google_compute_global_address" "worker_range" {
  name          = "tf-test-worker-pool-range%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.network.id
}

resource "google_service_networking_connection" "worker_pool_conn" {
  network                 = data.google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.worker_range.name]
}

resource "google_cloudbuild_worker_pool" "pool" {
	name = "pool%{random_suffix}"
	location = "europe-west1"
	worker_config {
		disk_size_gb = 101
		machine_type = "e2-standard-4"
		no_external_ip = false
	}
	network_config {
		peered_network = data.google_compute_network.network.id
		peered_network_ip_range = "/29"
	}
	depends_on = [google_service_networking_connection.worker_pool_conn]
}
`, context)
}

func funcAccTestCloudbuildWorkerPoolCheckDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloudbuild_worker_pool" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CloudBuildBasePath}}projects/{{project}}/locations/{{location}}/workerPools/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("CloudbuildWorkerPool still exists at %s", url)
			}
		}

		return nil
	}
}
