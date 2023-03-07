package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccAccessContextManagerAuthorizedOrgsDesc_basicTest(t *testing.T) {
	context := map[string]interface{}{
		"org_id": GetTestOrgFromEnv(t),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAuthorizedOrgsDescDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAuthorizedOrgsDesc_accessContextManagerAuthorizedOrgsDescBasicExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_authorized_orgs_desc.authorized-orgs-desc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccAccessContextManagerAuthorizedOrgsDesc_accessContextManagerAuthorizedOrgsDescBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_access_context_manager_authorized_orgs_desc" "authorized-orgs-desc" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/authorizedOrgsDescs/fakeDescName"
  authorization_type = "AUTHORIZATION_TYPE_TRUST"
  asset_type = "ASSET_TYPE_CREDENTIAL_STRENGTH"
  authorization_direction = "AUTHORIZATION_DIRECTION_TO"
  orgs = ["organizations/12345", "organizations/98765"]
}

resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}
`, context)
}

func testAccCheckAccessContextManagerAuthorizedOrgsDescDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_authorized_orgs_desc" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("AccessContextManagerAuthorizedOrgsDesc still exists at %s", url)
			}
		}

		return nil
	}
}
