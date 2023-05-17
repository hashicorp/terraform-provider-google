package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerEgressPolicy_basicTest(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	org := acctest.GetTestOrgFromEnv(t)
	projects := BootstrapServicePerimeterProjects(t, 1)
	policyTitle := RandString(t, 10)
	perimeterTitle := "perimeter"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerEgressPolicy_basic(org, policyTitle, perimeterTitle, projects[0].ProjectNumber),
			},
			{
				ResourceName:      "google_access_context_manager_egress_policy.test-access1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerEgressPolicy_destroy(org, policyTitle, perimeterTitle),
				Check:  testAccCheckAccessContextManagerEgressPolicyDestroyProducer(t),
			},
		},
	})
}

func testAccCheckAccessContextManagerEgressPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_egress_policy" {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{egress_policy_name}}")
			if err != nil {
				return err
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
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
}

func testAccAccessContextManagerEgressPolicy_basic(org, policyTitle, perimeterTitleName string, projectNumber1 int64) string {
	return fmt.Sprintf(`
%s

resource "google_access_context_manager_egress_policy" "test-access1" {
  egress_policy_name = google_access_context_manager_service_perimeter.test-access.name
  resource           = "projects/%d"
}

`, testAccAccessContextManagerEgressPolicy_destroy(org, policyTitle, perimeterTitleName), projectNumber1)
}

func testAccAccessContextManagerEgressPolicy_destroy(org, policyTitle, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  status {
    restricted_services = ["storage.googleapis.com"]
	egress_policies {
		egress_from {
			identity_type = "ANY_USER_ACCOUNT"
		}
	}
  }

  lifecycle {
  	ignore_changes = [status[0].resources]
  }
}
`, org, policyTitle, perimeterTitleName, perimeterTitleName)
}
