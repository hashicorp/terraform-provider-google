package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test that a service account key can be created and destroyed
func TestAccServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + randString(t, 10)
	displayName := "Terraform Test"
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountKey(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(t, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func TestAccServiceAccountKey_fromEmail(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + randString(t, 10)
	displayName := "Terraform Test"
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountKey_fromEmail(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(t, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func TestAccServiceAccountKey_fromCertificate(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + randString(t, 10)
	displayName := "Terraform Test"
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountKey_fromCertificate(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(t, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountKeyExists(t *testing.T, r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)

		_, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccServiceAccountKey(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = google_service_account.acceptance.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}
`, account, name)
}

func testAccServiceAccountKey_fromEmail(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = google_service_account.acceptance.email
  public_key_type    = "TYPE_X509_PEM_FILE"
}
`, account, name)
}

func testAccServiceAccountKey_fromCertificate(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = google_service_account.acceptance.email
  public_key_data    = filebase64("test-fixtures/serviceaccount/public_key.pem")
}
`, account, name)
}
