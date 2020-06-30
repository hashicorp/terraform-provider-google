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

		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("can't find %s in state", value)
		}

		_, ok = rs.Primary.Attributes["id_token"]
		if !ok {
			return fmt.Errorf("id_token not found")
		}

		// TODO, validate id_token
		//v, ok := rs.Primary.Attributes["id_token"]
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
				Config: testAccCheckGoogleServiceAccountIdToken_basic(targetAudience),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_audience", targetAudience),
					testAccCheckServiceAccountIdTokenValue(resourceName, fakeIdToken),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIdToken_basic(targetAudience string) string {

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
					testAccCheckServiceAccountIdTokenValue(resourceName, fakeIdToken),
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
	target_service_account = "%s"
	target_audience = "%s"
}
`, targetServiceAccount, targetServiceAccount, targetAudience)
}
