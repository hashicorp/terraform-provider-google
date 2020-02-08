package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be ran serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerServicePerimeterResource_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)
	projects := BootstrapServicePerimeterProjects(t, 2)
	policyTitle := "my policy"
	perimeterTitle := "perimeter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeterResource_basic(org, policyTitle, perimeterTitle, projects[0].ProjectNumber, projects[1].ProjectNumber),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter_resource.test-access1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter_resource.test-access2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the service perimeter to still exist
			{
				Config: testAccAccessContextManagerServicePerimeterResource_destroy(org, policyTitle, perimeterTitle),
				Check:  testAccCheckAccessContextManagerServicePerimeterResourceDestroy,
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimeterResourceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_access_context_manager_service_perimeter_resource" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{perimeter_name}}")
		if err != nil {
			return err
		}

		res, err := sendRequest(config, "GET", "", url, nil)
		if err != nil {
			return err
		}

		v, ok := res["status"]
		if !ok || v == nil {
			return nil
		}

		res = v.(map[string]interface{})
		v, ok = res["resources"]
		if !ok || v == nil {
			return nil
		}

		resources := v.([]interface{})
		if len(resources) == 0 {
			return nil
		}

		return fmt.Errorf("expected 0 resources in perimeter, found %d: %v", len(resources), resources)
	}

	return nil
}

func testAccAccessContextManagerServicePerimeterResource_basic(org, policyTitle, perimeterTitleName string, projectNumber1, projectNumber2 int64) string {
	return fmt.Sprintf(`
%s

resource "google_access_context_manager_service_perimeter_resource" "test-access1" {
  perimeter_name = google_access_context_manager_service_perimeter.test-access.name
  resource = "projects/%d"
}

resource "google_access_context_manager_service_perimeter_resource" "test-access2" {
  perimeter_name = google_access_context_manager_service_perimeter.test-access.name
  resource = "projects/%d"
}
`, testAccAccessContextManagerServicePerimeterResource_destroy(org, policyTitle, perimeterTitleName), projectNumber1, projectNumber2)
}

func testAccAccessContextManagerServicePerimeterResource_destroy(org, policyTitle, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["storage.googleapis.com"]
  }

  lifecycle {
  	ignore_changes = [status[0].resources]
  }
}
`, org, policyTitle, perimeterTitleName, perimeterTitleName)
}
