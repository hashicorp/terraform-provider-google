package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test that a service account resource can be created, updated, and destroyed
func TestAccServiceAccount_basic(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(10)
	uniqueId := ""
	displayName := "Terraform Test"
	displayName2 := "Terraform Test Update"
	desc := "test description"
	desc2 := "test description update"
	project := getTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("%s/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     expectedEmail,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step updates the service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName2, desc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The third step explicitly adds the same default project to the service account configuration
			// and ensure the service account is not recreated by comparing the value of its unique_id with the one from the previous step
			{
				Config: testAccServiceAccountWithProject(project, accountId, displayName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttrPtr(
						"google_service_account.acceptance", "unique_id", &uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccServiceAccountBasic(account, name, desc string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    account_id = "%v"
    display_name = "%v"
    description = "%v"
}
`, account, name, desc)
}

func testAccServiceAccountWithProject(project, account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    project = "%v"
    account_id = "%v"
    display_name = "%v"
    description = "foo"
}
`, project, account, name)
}
