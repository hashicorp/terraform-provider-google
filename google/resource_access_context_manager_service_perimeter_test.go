package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be ran serially. See AccessPolicy for the test runner.
func testAccAccessContextManagerServicePerimeter_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t),
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

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t),
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
			{
				Config: testAccAccessContextManagerServicePerimeter_updateAllowed(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_updateDryrun(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_updateAllowed(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_service_perimeter" {
				continue
			}

			config := googleProviderConfig(t)

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

func testAccAccessContextManagerServicePerimeter_updateAllowed(org, policyTitle, levelTitleName, perimeterTitleName string) string {
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

	vpc_accessible_services {
	  enable_restriction = true
	  allowed_services   = ["bigquery.googleapis.com"]
	}
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}

func testAccAccessContextManagerServicePerimeter_updateDryrun(org, policyTitle, levelTitleName, perimeterTitleName string) string {
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
  }

  spec {
    restricted_services = ["storage.googleapis.com"]
	access_levels       = [google_access_context_manager_access_level.test-access.name]
  }

  use_explicit_dry_run_spec = true
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}
