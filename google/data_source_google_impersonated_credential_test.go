package google

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleImpersonatedCredential_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_impersonated_credential.current"

	sourceAccessToken := "foo"
	targetServiceAccount := getTestServiceAccountFromEnv(t)
	scopes := []string{"https://www.googleapis.com/auth/cloud-platform"}
	delegates := []string{"projects/-/serviceAccounts/impersonated-account@some-project-111.iam.gserviceaccount.com"}
	lifetime := "30s"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleImpersonatedCredential_datasource(sourceAccessToken, targetServiceAccount, scopes, delegates, lifetime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_access_token", sourceAccessToken),
					resource.TestCheckResourceAttr(resourceName, "target_service_account", targetServiceAccount),
					resource.TestCheckResourceAttrSet(resourceName, "scopes"),
					resource.TestCheckResourceAttrSet(resourceName, "delegates"),
					resource.TestCheckResourceAttr(resourceName, "lifetime", lifetime),
				),
			},
		},
	})
}

func testAccCheckGoogleImpersonatedCredential_datasource(sourceAccessToken string, targetServiceAccount string, scopes []string, delegates []string, lifetime string) string {
	return fmt.Sprintf(`
	data "google_impersonated_credential" "current" {
		source_access_token = "%s"
		target_service_account = "%s"
		scopes = "%s"
		delegates = "%s"
		lifetime = "%s"
}
	`, sourceAccessToken, targetServiceAccount, scopes, delegates, lifetime)
}
