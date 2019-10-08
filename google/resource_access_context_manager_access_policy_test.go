package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be ran serially
func TestAccAccessContextManager(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"access_policy":            testAccAccessContextManagerAccessPolicy_basicTest,
		"service_perimeter":        testAccAccessContextManagerServicePerimeter_basicTest,
		"service_perimeter_update": testAccAccessContextManagerServicePerimeter_updateTest,
		"access_level":             testAccAccessContextManagerAccessLevel_basicTest,
		"access_level_full":        testAccAccessContextManagerAccessLevel_fullTest,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccAccessContextManagerAccessPolicy_basicTest(t *testing.T) {
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my new policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_access_context_manager_access_policy" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}accessPolicies/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", "", url, nil)
		if err == nil {
			return fmt.Errorf("AccessPolicy still exists at %s", url)
		}
	}

	return nil
}

func testAccAccessContextManagerAccessPolicy_basic(org, title string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}
`, org, title)
}
