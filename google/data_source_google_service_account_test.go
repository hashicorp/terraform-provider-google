package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDatasourceGoogleServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account.acceptance"
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccount_basic(account),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "id", fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", getTestProjectFromEnv(), account, getTestProjectFromEnv())),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "unique_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccount_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "Testing Account"
}

data "google_service_account" "acceptance" {
  account_id = google_service_account.acceptance.account_id
}
`, account)
}
