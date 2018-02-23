package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test that a service account resource can be created, updated, and destroyed
func TestAccServiceAccount_basic(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(10)
	uniqueId := ""
	displayName := "Terraform Test"
	displayName2 := "Terraform Test Update"
	project := getTestProjectFromEnv()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			resource.TestStep{
				Config: testAccServiceAccountBasic(accountId, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountExists("google_service_account.acceptance"),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
				),
			},
			// The second step updates the service account
			resource.TestStep{
				Config: testAccServiceAccountBasic(accountId, displayName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountNameModified("google_service_account.acceptance", displayName2),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			// The third step explicitely adds the same default project to the service account configuration
			// and ensure the service account is not recreated by comparing the value of its unique_id with the one from the previous step
			resource.TestStep{
				Config: testAccServiceAccountWithProject(project, accountId, displayName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountNameModified("google_service_account.acceptance", displayName2),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttrPtr(
						"google_service_account.acceptance", "unique_id", &uniqueId),
				),
			},
		},
	})
}

// Test that a service account resource can be created with a policy, updated,
// and destroyed.
func TestAccServiceAccount_createPolicy(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic service account with an IAM policy
			resource.TestStep{
				Config: testAccServiceAccountPolicy(accountId, getTestProjectFromEnv()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountPolicyCount("google_service_account.acceptance", 1),
				),
			},
			// The second step updates the service account with no IAM policy
			resource.TestStep{
				Config: testAccServiceAccountBasic(accountId, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountPolicyCount("google_service_account.acceptance", 0),
				),
			},
			// The final step re-applies the IAM policy
			resource.TestStep{
				Config: testAccServiceAccountPolicy(accountId, getTestProjectFromEnv()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountPolicyCount("google_service_account.acceptance", 1),
				),
			},
		},
	})
}

func testAccStoreServiceAccountUniqueId(uniqueId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*uniqueId = s.RootModule().Resources["google_service_account.acceptance"].Primary.Attributes["unique_id"]
		return nil
	}
}

func testAccCheckGoogleServiceAccountPolicyCount(r string, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Config)
		p, err := getServiceAccountIamPolicy(s.RootModule().Resources[r].Primary.ID, c)
		if err != nil {
			return fmt.Errorf("Failed to retrieve IAM Policy for service account: %s", err)
		}
		if len(p.Bindings) != n {
			return fmt.Errorf("The service account has %v bindings but %v were expected", len(p.Bindings), n)
		}
		return nil
	}
}

func testAccCheckGoogleServiceAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckGoogleServiceAccountNameModified(r, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.Attributes["display_name"] != n {
			return fmt.Errorf("display_name is %q expected %q", rs.Primary.Attributes["display_name"], n)
		}

		return nil
	}
}

func testAccServiceAccountBasic(account, name string) string {
	t := `resource "google_service_account" "acceptance" {
    account_id = "%v"
	display_name = "%v"
 }`
	return fmt.Sprintf(t, account, name)
}

func testAccServiceAccountWithProject(project, account, name string) string {
	t := `resource "google_service_account" "acceptance" {
    project = "%v"
    account_id = "%v"
    display_name = "%v"
 }`
	return fmt.Sprintf(t, project, account, name)
}

func testAccServiceAccountPolicy(account, project string) string {

	t := `resource "google_service_account" "acceptance" {
    account_id = "%v"
    display_name = "%v"
    policy_data = "${data.google_iam_policy.service_account.policy_data}"
}

data "google_iam_policy" "service_account" {
  binding {
    role = "roles/iam.serviceAccountActor"
    members = [
      "serviceAccount:%v@%v.iam.gserviceaccount.com",
    ]
  }
}`

	return fmt.Sprintf(t, account, account, account, project)
}
