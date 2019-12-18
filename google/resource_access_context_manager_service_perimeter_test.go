package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be ran serially. See AccessPolicy for the test runner.
func testAccAccessContextManagerServicePerimeter_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerServicePerimeterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeter_basic(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAccessContextManagerServicePerimeter_updateTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerServicePerimeterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeter_basic(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_update(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimeterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_access_context_manager_service_perimeter" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", "", url, nil)
		if err == nil {
			return fmt.Errorf("ServicePerimeter still exists at %s", url)
		}
	}

	return nil
}

func testAccAccessContextManagerServicePerimeter_basic(org, policyTitle, levelTitleName, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["storage.googleapis.com"]
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}

func testAccAccessContextManagerServicePerimeter_update(org, policyTitle, levelTitleName, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["bigquery.googleapis.com"]
    access_levels       = [google_access_context_manager_access_level.test-access.name]
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}
