package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerAccessLevel_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAccessLevelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessLevel_basic(org, "my policy", "level"),
			},
			{
				ResourceName:      "google_access_context_manager_access_level.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessLevel_basicUpdated(org, "my new policy", "level"),
			},
			{
				ResourceName:      "google_access_context_manager_access_level.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAccessContextManagerAccessLevel_fullTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAccessLevelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessLevel_full(org, "my policy", "level"),
			},
			{
				ResourceName:      "google_access_context_manager_access_level.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessLevelDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_access_level" {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("AccessLevel still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerAccessLevel_customTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAccessLevelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessLevel_custom(org, "my policy", "level"),
			},
			{
				ResourceName:      "google_access_context_manager_access_level.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAccessContextManagerAccessLevel_basic(org, policyTitle, levelTitleName string) string {
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
`, org, policyTitle, levelTitleName, levelTitleName)
}

func testAccAccessContextManagerAccessLevel_custom(org, policyTitle, levelTitleName string) string {
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
    custom {
		expr {
			expression = "device.os_type == OsType.DESKTOP_MAC"
		}
  }
}
`, org, policyTitle, levelTitleName, levelTitleName)
}

func testAccAccessContextManagerAccessLevel_basicUpdated(org, policyTitle, levelTitleName string) string {
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
    combining_function = "OR"
    conditions {
      ip_subnetworks = ["192.0.2.0/24"]
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName)
}

func testAccAccessContextManagerAccessLevel_full(org, policyTitle, levelTitleName string) string {
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
      members = ["user:test@google.com", "user:test2@google.com"]
      negate = false
      device_policy {
        require_screen_lock = false
        require_admin_approval = false
        require_corp_owned = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
          require_verified_chrome_os = true
        }
      }
      regions = [
        "IT",
        "US",
      ]
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName)
}
