// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerAccessLevelCondition_basicTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	project := envvar.GetTestProjectFromEnv()

	serviceAccountName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"ipSubnetworks": []interface{}{"192.0.4.0/24"},
		"members":       []interface{}{"user:test@google.com", "user:test2@google.com", fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", serviceAccountName, project)},
		"devicePolicy": map[string]interface{}{
			"requireCorpOwned": true,
			"osConstraints": []interface{}{
				map[string]interface{}{
					"osType": "DESKTOP_CHROME_OS",
				},
			},
		},
		"regions": []interface{}{"IT", "US"},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerAccessLevelConditionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessLevelCondition_basic(org, "my policy", "level", serviceAccountName),
				Check:  testAccCheckAccessContextManagerAccessLevelConditionPresent(t, "google_access_context_manager_access_level_condition.access-level-condition", expected),
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessLevelConditionPresent(t *testing.T, n string, expected map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{access_level}}")
		if err != nil {
			return err
		}

		al, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}
		conditions := al["basic"].(map[string]interface{})["conditions"].([]interface{})
		for _, c := range conditions {
			if reflect.DeepEqual(c, expected) {
				return nil
			}
		}
		return fmt.Errorf("Did not find condition %+v", expected)
	}
}

func testAccCheckAccessContextManagerAccessLevelConditionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_access_level_condition" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{access_level}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AccessLevel still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerAccessLevelCondition_basic(org, policyTitle, levelTitleName, saName string) string {
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
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
          require_verified_chrome_os = true
        }
      }
      regions = [
  "CH",
  "IT",
  "US",
      ]
    }

    conditions {
      ip_subnetworks = ["176.0.4.0/24"]
    }
  }

  lifecycle {
    ignore_changes = [basic.0.conditions]
  }
}

resource "google_service_account" "created-later" {
  account_id = "%s"
}

resource "google_access_context_manager_access_level_condition" "access-level-condition" {
  access_level = google_access_context_manager_access_level.test-access.name
  ip_subnetworks = ["192.0.4.0/24"]
  members = ["user:test@google.com", "user:test2@google.com", "serviceAccount:${google_service_account.created-later.email}"]
  negate = false
  device_policy {
    require_screen_lock = false
    require_admin_approval = false
    require_corp_owned = true
    os_constraints {
      os_type = "DESKTOP_CHROME_OS"
    }
  }
  regions = [
    "IT",
    "US",
  ]
}
`, org, policyTitle, levelTitleName, levelTitleName, saName)
}
