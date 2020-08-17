package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDatasourceGoogleServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_key.acceptance"
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	serviceAccountName := fmt.Sprintf(
		"projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com",
		getTestProjectFromEnv(),
		account,
		getTestProjectFromEnv(),
	)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceGoogleServiceAccountKey(account),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(t, resourceName),
					// Check that the 'name' starts with the service account name
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(serviceAccountName)),
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
  service_account_id = google_service_account.acceptance.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

data "google_service_account_key" "acceptance" {
  name = google_service_account_key.acceptance.name
}
`, account)
}
