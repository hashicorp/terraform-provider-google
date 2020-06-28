package google

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const targetAudience = "https://foo.bar/"
const fakeIdToken = "eyJhbGciOiJSUzI1NiIsIm..."

func testAccCheckServiceAccountIdTokenValue(name, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		// TODO: validate the token belongs to the service account
		if rs.Value == "" {
			return fmt.Errorf("%s Cannot be empty", name)
		}

		return nil
	}
}

func TestAccDataSourceGoogleServiceAccountIdToken_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_id_token.default"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountIdToken_datasource(targetAudience),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_audience", targetAudience),
					testAccCheckServiceAccountIdTokenValue("google_service_account_id_token.default", fakeIdToken),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIdToken_datasource(targetAudience string) string {

	return fmt.Sprintf(`
data "google_service_account_id_token" "default" {
  target_audience = "%s"
}
`, targetAudience)
}

func TestAccDataSourceGoogleServiceAccountIdToken_impersonation(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_id_token.default"
	serviceAccount := getTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := BootstrapServiceAccount(t, getTestProjectFromEnv(), serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountIdToken_impersonation_datasource(targetAudience, targetServiceAccountEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_audience", targetAudience),
					testAccCheckServiceAccountIdTokenValue("google_service_account_id_token.default", fakeIdToken),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIdToken_impersonation_datasource(targetAudience string, targetServiceAccount string) string {

	return fmt.Sprintf(`
data "google_service_account_access_token" "default" {
	target_service_account = "%s"
	scopes                 = ["userinfo-email", "https://www.googleapis.com/auth/cloud-platform"]
	lifetime               = "30s"
}

provider google {
	alias  = "impersonated"
	access_token = data.google_service_account_access_token.default.access_token
}

data "google_service_account_id_token" "default" {
	provider = google.impersonated
	target_audience = "%s"
}
`, targetServiceAccount, targetAudience)
}
