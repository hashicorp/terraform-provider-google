package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDatasourceGoogleServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_key.acceptance"
	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatasourceGoogleServiceAccountKey(account),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "key_algorithm"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
				),
			},
		},
	})
}

func testAccDatasourceGoogleServiceAccountKey(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
	account_id = "%s"
}

resource "google_service_account_key" "acceptance" {
	service_account_id = "${google_service_account.acceptance.name}"
	public_key_type = "TYPE_X509_PEM_FILE"
}

data "google_service_account_key" "acceptance" {
	service_account_id = "${google_service_account_key.acceptance.id}"
}`, account)
}
