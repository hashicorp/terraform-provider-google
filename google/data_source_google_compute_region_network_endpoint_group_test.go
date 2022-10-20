package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRegionNetworkEndpointGroup_basic(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"region":        "us-central1",
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionNetworkEndpointGroup_basic(context),
				Check:  checkDataSourceStateMatchesResourceStateWithIgnores("data.google_compute_region_network_endpoint_group.cloudrun_neg", "google_compute_region_network_endpoint_group.cloudrun_neg", map[string]struct{}{"name": {}, "region": {}}),
			},
		},
	})
}

func testAccDataSourceRegionNetworkEndpointGroup_basic(context map[string]interface{}) string {
	return Nprintf(`
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
