package google

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBinaryAuthorizationPolicy_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-test-" + acctest.RandString(10)
	billingId := getTestBillingAccountFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(pid),
			},
		},
	})
}

// Because Container Analysis is still in beta, we can't run any of the tests that call that
// resource without vendoring in the full beta provider.

func testAccCheckBinaryAuthorizationPolicyDefault(pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		url := fmt.Sprintf("https://binaryauthorization.googleapis.com/v1beta1/projects/%s/policy", pid)
		pol, err := sendRequest(config, "GET", "", url, nil)
		if err != nil {
			return err
		}
		delete(pol, "updateTime")

		defaultPol := defaultBinaryAuthorizationPolicy(pid)
		if !reflect.DeepEqual(pol, defaultPol) {
			return fmt.Errorf("Policy for project %s was %v, expected default policy %v", pid, pol, defaultPol)
		}
		return nil
	}
}

func testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billing string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}
`, pid, pname, org, billing)
}

func testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billing string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}

resource "google_binary_authorization_policy" "policy" {
  project = google_project.project.project_id

  admission_whitelist_patterns {
    name_pattern = "gcr.io/google_containers/*"
  }

  default_admission_rule {
    evaluation_mode  = "ALWAYS_DENY"
    enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"
  }

  depends_on = [google_project_service.binauthz]
}
`, pid, pname, org, billing)
}
