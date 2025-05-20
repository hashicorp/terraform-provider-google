// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudbuild_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCloudbuildWorkerPool_withComputedAnnotations(t *testing.T) {
	// Skip it in VCR test because of the randomness of uuid in "annotations" field
	// which causes the replaying mode after recording mode failing in VCR test
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: funcAccTestCloudbuildWorkerPoolCheckDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildWorkerPool_updated(context),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
				ResourceName:            "google_cloudbuild_worker_pool.pool",
			},
			{
				Config: testAccCloudbuildWorkerPool_withComputedAnnotations(context),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
				ResourceName:            "google_cloudbuild_worker_pool.pool",
			},
		},
	})
}

func TestAccCloudbuildWorkerPool_basic(t *testing.T) {
	t.Parallel()

	testNetworkName := acctest.BootstrapSharedTestNetwork(t, "attachment-network")
	subnetName := acctest.BootstrapSubnet(t, "tf-test-subnet", testNetworkName)
	networkAttachmentName := acctest.BootstrapNetworkAttachment(t, "tf-test-attachment", subnetName)

	// Need to have the full network attachment name in the format project/{project_id}/regions/{region_id}/networkAttachments/{networkAttachmentName}
	fullFormNetworkAttachmentName := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), networkAttachmentName)

	context := map[string]interface{}{
		"random_suffix":      acctest.RandString(t, 10),
		"project":            envvar.GetTestProjectFromEnv(),
		"network_attachment": fullFormNetworkAttachmentName,
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
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
				ResourceName:            "google_cloudbuild_worker_pool.pool",
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

	// private_service_connect feature is not supported yet. b/394920388
	// private_service_connect {
	// 	network_attachment = "%{network_attachment}"
	// 	route_all_traffic = false
	// }
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

	annotations = {
		env                   = "foo"
		default_expiration_ms = 3600000
	}
}
`, context)
}

func testAccCloudbuildWorkerPool_withComputedAnnotations(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "random_uuid" "test" {
}

resource "google_cloudbuild_worker_pool" "pool" {
  name = "pool%{random_suffix}"
  location = "europe-west1"
  worker_config {
  disk_size_gb = 101
  machine_type = "e2-standard-4"
  no_external_ip = false
  }

  annotations = {
    env                   = "${random_uuid.test.result}"
    default_expiration_ms = 3600000
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
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "cloudbuild-workerpool-1"),
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
