package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerGcpUserAccessBinding_basicTest(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
  "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}



resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  group_key       = trimprefix(google_cloud_identity_group.group.id, "groups/")
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
}
`, context)
}

func testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_gcp_user_access_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}organizations/{{organization_id}}/gcpUserAccessBindings/{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("AccessContextManagerGcpUserAccessBinding still exists at %s", url)
			}
		}

		return nil
	}
}
