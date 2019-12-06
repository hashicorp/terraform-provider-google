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
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
			{
				Config: testAccCloudRunService_cloudRunServiceUpdate(name, project, "50"),
			},
			{
				ResourceName:            "google_cloud_run_service.default",
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
