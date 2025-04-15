// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeStoragePool_computeStoragePool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeStoragePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeStoragePool_computeStoragePoolFullExample(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "labels", "terraform_labels", "zone"},
			},
			{
				Config: testAccComputeStoragePool_computeStoragePool_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_storage_pool.test-storage-pool-full", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "labels", "terraform_labels", "zone"},
			},
		},
	})
}

func testAccComputeStoragePool_computeStoragePool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "test-storage-pool-full" {
  name = "tf-test-storage-pool-full%{random_suffix}"

  description = "Hyperdisk Balanced storage pool"

  capacity_provisioning_type   = "STANDARD"
  pool_provisioned_capacity_gb = "11264"

  performance_provisioning_type = "STANDARD"
  pool_provisioned_iops         = "20000"
  pool_provisioned_throughput   = "2048"

	storage_pool_type = "hyperdisk-balanced"

	deletion_protection = false

	zone = "us-central1-a"
}

data "google_project" "project" {}

data "google_compute_storage_pool_types" "balanced" {
  zone = "us-central1-a"
	storage_pool_type = "hyperdisk-balanced"
}
`, context)
}

func TestAccComputeStoragePool_computeStoragePoolBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeStoragePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeStoragePool_computeStoragePoolBasicExample(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "zone"},
			},
		},
	})
}

func testAccComputeStoragePool_computeStoragePoolBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "test-storage-pool-basic" {
  name = "tf-test-storage-pool-basic%{random_suffix}"

  pool_provisioned_capacity_gb = "10240"

  pool_provisioned_throughput = 100

  storage_pool_type = "hyperdisk-throughput"

  zone = "us-central1-a"

  deletion_protection = false
}

data "google_project" "project" {}

`, context)
}

func TestAccComputeStoragePool_computeStoragePoolFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeStoragePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeStoragePool_computeStoragePoolFullExample(context),
			},
			{
				ResourceName:            "google_compute_storage_pool.test-storage-pool-full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "zone"},
			},
		},
	})
}

func testAccComputeStoragePool_computeStoragePoolFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "test-storage-pool-full" {
  name = "tf-test-storage-pool-full%{random_suffix}"

  description = "Hyperdisk Balanced storage pool"

  capacity_provisioning_type   = "STANDARD"
  pool_provisioned_capacity_gb = "10240"

  performance_provisioning_type = "STANDARD"
  pool_provisioned_iops         = "10000"
  pool_provisioned_throughput   = "1024"

  storage_pool_type = data.google_compute_storage_pool_types.balanced.self_link

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_storage_pool_types" "balanced" {
  zone = "us-central1-a"
	storage_pool_type = "hyperdisk-balanced"
}
`, context)
}

func testAccCheckComputeStoragePoolDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_storage_pool" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/storagePools/{{name}}")
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
				return fmt.Errorf("ComputeStoragePool still exists at %s", url)
			}
		}

		return nil
	}
}
