package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test that a service account key can be created and destroyed
func TestAccGoogleServiceAccountKey_basic(t *testing.T) {
	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + acctest.RandString(10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleServiceAccountKey(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountKeyExists(r string) resource.TestCheckFunc {
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

func testAccGoogleServiceAccountKey(account, name string) string {
	return fmt.Sprintf(`resource "google_service_account" "acceptance" {
	account_id = "%v"
	display_name = "%v"

	}

	resource "google_service_account_key" "acceptance" {
		service_account_id = "${google_service_account.acceptance.id}"
		public_key_type = "TYPE_X509_PEM_FILE"
	}`, account, name)
}
