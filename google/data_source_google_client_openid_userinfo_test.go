package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleClientOpenIDUserinfo_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientOpenIDUserinfo_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_client_openid_userinfo.me", "email"),
				),
			},
		},
	})
}

const testAccCheckGoogleClientOpenIDUserinfo_basic = `
provider "google" {
  alias = "google-scoped"

  # We need to add an additional scope to test this; because our tests rely on
  # every env var being set, we can just add an alias with the appropriate
  # scopes. This will fail if someone uses an access token instead of creds
  # unless they've configured the userinfo.email scope.
  scopes = [
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
    "https://www.googleapis.com/auth/devstorage.full_control",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

data "google_client_openid_userinfo" "me" {
  provider = "google.google-scoped"
}
`
