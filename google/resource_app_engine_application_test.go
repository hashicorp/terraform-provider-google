package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppEngineApplication_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := acctest.RandomWithPrefix("tf-test")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineApplication_basic(pid, org),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_app_engine_application.acceptance", "url_dispatch_rule.#"),
					resource.TestCheckResourceAttrSet("google_app_engine_application.acceptance", "name"),
					resource.TestCheckResourceAttrSet("google_app_engine_application.acceptance", "code_bucket"),
					resource.TestCheckResourceAttrSet("google_app_engine_application.acceptance", "default_hostname"),
					resource.TestCheckResourceAttrSet("google_app_engine_application.acceptance", "default_bucket"),
				),
			},
			{
				ResourceName:      "google_app_engine_application.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppEngineApplication_update(pid, org),
			},
			{
				ResourceName:      "google_app_engine_application.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppEngineApplication_basic(pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_app_engine_application" "acceptance" {
  project        = google_project.acceptance.project_id
  auth_domain    = "hashicorptest.com"
  location_id    = "us-central"
  serving_status = "SERVING"
}
`, pid, pid, org)
}

func testAccAppEngineApplication_update(pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_app_engine_application" "acceptance" {
  project        = google_project.acceptance.project_id
  auth_domain    = "tf-test.club"
  location_id    = "us-central"
  serving_status = "USER_DISABLED"
}
`, pid, pid, org)
}
