// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAppEngineFlexibleAppVersion_appEngineFlexibleAppVersionExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineFlexibleAppVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineFlexibleAppVersion_appEngineFlexibleAppVersionExample(context),
			},
			{
				ResourceName:            "google_app_engine_flexible_app_version.myapp_v1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"beta_settings", "env_variables", "deployment", "entrypoint", "service", "delete_service_on_destroy"},
			},
		},
	})
}

func testAccAppEngineFlexibleAppVersion_appEngineFlexibleAppVersionExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_app_engine_flexible_app_version" "myapp_v1" {
  version_id = "v1"
  service    = "tf-test-service-%{random_suffix}"
  runtime    = "nodejs"

  entrypoint {
    shell = "node ./app.js"
  }

  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    }
  }

  liveness_check {
    path = "/"
  }

  readiness_check {
    path = "/"
  }

  env_variables = {
    port = "8080"
  }

  automatic_scaling {
    cool_down_period = "120s"
    cpu_utilization {
      target_utilization = 0.5
    }
  }

  delete_service_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  name = "tf-test-appengine-static-content%{random_suffix}"
}

resource "google_storage_bucket_object" "object" {
  name   = "hello-world.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world.zip"
}
`, context)
}

func testAccCheckAppEngineFlexibleAppVersionDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_app_engine_flexible_app_version" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		log.Printf("[DEBUG] Ignoring destroy during test")
	}

	return nil
}
