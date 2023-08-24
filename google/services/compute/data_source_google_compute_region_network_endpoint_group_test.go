// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceRegionNetworkEndpointGroup_basic(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"region":        "us-central1",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionNetworkEndpointGroup_basic(context),
				Check:  acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_compute_region_network_endpoint_group.cloudrun_neg", "google_compute_region_network_endpoint_group.cloudrun_neg", map[string]struct{}{"name": {}, "region": {}}),
			},
		},
	})
}

func testAccDataSourceRegionNetworkEndpointGroup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_compute_region_network_endpoint_group" "cloudrun_neg" {
    name                  = "cloud-run-rneg-%{random_suffix}"
    network_endpoint_type = "SERVERLESS"
    region                = "%{region}"
    project     = "%{project}"
    cloud_run {
      service = google_cloud_run_service.cloudrun_neg.name
    }
  }

  resource "google_cloud_run_service" "cloudrun_neg" {
    name     = "tf-test-cloudrun-neg%{random_suffix}"
    location = "us-central1"
    template {
      spec {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/hello"
        }
      }
    }

    traffic {
      percent         = 100
      latest_revision = true
    }
  }

  data "google_compute_region_network_endpoint_group" "cloudrun_neg" {
      name = google_compute_region_network_endpoint_group.cloudrun_neg.name
      region = "%{region}"
  }
`, context)
}
