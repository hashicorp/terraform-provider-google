// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComposerEnvironment_basic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComposerEnvironmentMeta("data.google_composer_environment.test"),
				),
			},
		},
	})
}

func testAccCheckGoogleComposerEnvironmentMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find environment data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("environment data source ID not set.")
		}

		configCountStr, ok := rs.Primary.Attributes["config.#"]
		if !ok {
			return errors.New("can't find 'config' attribute")
		}

		configCount, err := strconv.Atoi(configCountStr)
		if err != nil {
			return errors.New("failed to read number of valid config entries")
		}
		if configCount < 1 {
			return fmt.Errorf("expected at least 1 valid config entry, received %d, this is most likely a bug",
				configCount)
		}

		for i := 0; i < configCount; i++ {
			idx := "config." + strconv.Itoa(i)

			if v, ok := rs.Primary.Attributes[idx+".airflow_uri"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing airflow_uri", i)
			}
			if v, ok := rs.Primary.Attributes[idx+".dag_gcs_prefix"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing dag_gcs_prefix", i)
			}
			if v, ok := rs.Primary.Attributes[idx+".gke_cluster"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing gke_cluster", i)
			}
		}

		return nil
	}
}

func testAccDataSourceComposerEnvironment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_composer_environment" "test" {
	name   = "tf-test-composer-env-%{random_suffix}"
	region = "us-central1"

	config {
		node_config {
			network    = google_compute_network.test.self_link
			subnetwork = google_compute_subnetwork.test.self_link
			zone       = "us-central1-a"
		}
		software_config {
			image_version = "composer-1-airflow-2"
		}
	}
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
	name                    = "tf-test-composer-net-%{random_suffix}"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
	name          = "tf-test-composer-subnet-%{random_suffix}"
	ip_cidr_range = "10.2.0.0/16"
	region        = "us-central1"
	network       = google_compute_network.test.self_link
}

data "google_composer_environment" "test" {
	name   = google_composer_environment.test.name
	region = google_composer_environment.test.region
}
`, context)
}
