package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudRunService_cloudRunServiceUpdate(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	name := "tftest-cloudrun-" + acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunService_cloudRunServiceUpdate(name, project, "10"),
			},
			{
				ResourceName:            "google_cloud_run_service.default",
				ImportStateId:           fmt.Sprintf("locations/us-central1/namespaces/%s/services/%s", project, name),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
			{
				Config: testAccCloudRunService_cloudRunServiceUpdate(name, project, "50"),
			},
			{
				ResourceName:            "google_cloud_run_service.default",
				ImportStateId:           fmt.Sprintf("locations/us-central1/namespaces/%s/services/%s", project, name),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
		},
	})
}

func TestAccCloudRunService_cloudRunServiceSql(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	name := "tftest-cloudrun-" + acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunService_cloudRunServiceSql(name, project),
			},
			{
				ResourceName:            "google_cloud_run_service.default",
				ImportStateId:           fmt.Sprintf("locations/us-central1/namespaces/%s/services/%s", project, name),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
		},
	})
}

func testAccCloudRunService_cloudRunServiceUpdate(name, project, concurrency string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_service" "default" {
  name     = "%s"
  location = "us-central1"

  metadata {
    namespace = "%s"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        args  = ["arrgs"]
      }
	  container_concurrency = %s
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, name, project, concurrency)
}

func testAccCloudRunService_cloudRunServiceSql(name, project string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_sql_database_instance" "instance" {
  name   = "tf-test-%s"
  region = "us-east1"
  settings {
    tier = "D0"
  }
}

resource "google_cloud_run_service" "default" {
  location = "us-east1"
  name     = "%s"

  metadata {
    namespace = "%s"
    labels = {
      "cloud.googleapis.com/location" = "us-east1"
      "foo"                           = "bar"
    }
  }

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "1000"
        "run.googleapis.com/cloudsql-instances" = "%s:us-east1:${google_sql_database_instance.instance.name}"
        "run.googleapis.com/client-name"        = "cloud-console"
      }
    }

    spec {
      service_account_name = "${data.google_project.project.number}-compute@developer.gserviceaccount.com"

      containers {
        image = "gcr.io/cloudrun/hello"
        args  = ["arrg2", "pirate"]
        resources {
          limits = {
            cpu    = "1000m"
            memory = "256Mi"
          }
        }
      }
      container_concurrency = 10
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, acctest.RandString(6), name, project, project)
}
