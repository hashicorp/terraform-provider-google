package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.
func testAccAccessContextManagerServicePerimeters_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerServicePerimetersDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeters_basic(org, "my policy", "level", "storage_perimeter", "bigtable_perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeters.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeters_update(org, "my policy", "level", "storage_perimeter", "bigquery_perimeter", "bigtable_perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeters.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeters_empty(org, "my policy", "level"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeters.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimetersDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_service_perimeters" {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{parent}}/servicePerimeters")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("ServicePerimeters still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerServicePerimeters_basic(org, policyTitle, levelTitleName, perimeterTitleName1, perimeterTitleName2 string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["storage.googleapis.com"]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigtable.googleapis.com"]
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName1, perimeterTitleName1, perimeterTitleName2, perimeterTitleName2)
}

func testAccAccessContextManagerServicePerimeters_update(org, policyTitle, levelTitleName, perimeterTitleName1, perimeterTitleName2, perimeterTitleName3 string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["storage.googleapis.com"]
      access_levels       = [google_access_context_manager_access_level.test-access.name]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigquery.googleapis.com"]
      access_levels       = [google_access_context_manager_access_level.test-access.name]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigtable.googleapis.com"]
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName1, perimeterTitleName1, perimeterTitleName2, perimeterTitleName2, perimeterTitleName3, perimeterTitleName3)
}

func testAccAccessContextManagerServicePerimeters_empty(org, policyTitle, levelTitleName string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
}
`, org, policyTitle, levelTitleName, levelTitleName)
}
