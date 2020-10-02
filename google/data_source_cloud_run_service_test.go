package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleCloudRunService_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunService_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_cloud_run_service.foo", "google_cloud_run_service.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCloudRunService_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudRunServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunService_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_cloud_run_service.foo", "google_cloud_run_service.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudRunService_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "foo" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_cloud_run_service" "foo" {
  name     = google_cloud_run_service.foo.name
  location = google_cloud_run_service.foo.location
  project  = google_cloud_run_service.foo.project
}
`, context)
}

func testAccDataSourceGoogleCloudRunService_optionalProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "foo" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_cloud_run_service" "foo" {
  name     = google_cloud_run_service.foo.name
  location = google_cloud_run_service.foo.location
}
`, context)
}
