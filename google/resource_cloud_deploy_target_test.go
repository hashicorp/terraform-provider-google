package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudDeployTarget_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataCatalogEntryGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudDeployTarget_cloudDeployTargetFullExample(context),
			},
			{
				ResourceName:            "google_cloud_deploy_target.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
			{
				Config: testAccCloudDeployTarget_cloudDeployTargetFullExample_update(context),
			},
			{
				ResourceName:            "google_cloud_deploy_target.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
		},
	})
}

func testAccCloudDeployTarget_cloudDeployTargetFullExample_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_deploy_target" "pipeline" {
  name          = "tf-test-tf-test%{random_suffix}"
  description   = "Target Prod Cluster"
  annotations = {
    generated-by = "magic-modules"
	another = "one"
  }
  labels = {
    env = "prod"
  }
  gke {
    cluster = "projects/%{project}/locations/us-central1/clusters/prod"
  }
  execution_configs {
    usages = ["RENDER"]
    service_account = data.google_compute_default_service_account.default.email
  }

  execution_configs {
    usages = ["DEPLOY"]
    service_account = "%{project}@appspot.gserviceaccount.com"
  }

}

data "google_compute_default_service_account" "default" {
}
`, context)
}
