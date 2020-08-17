package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformOauthIdpConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigBasic(context),
			},
			{
				ResourceName:      "google_identity_platform_oauth_idp_config.oauth_idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigUpdate(context),
			},
			{
				ResourceName:      "google_identity_platform_oauth_idp_config.oauth_idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_oauth_idp_config" "oauth_idp_config" {
  name          = "oidc.oauth-idp-config%{random_suffix}"
  display_name  = "Display Name"
  client_id     = "client-id"
  issuer        = "issuer"
  enabled       = true
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_oauth_idp_config" "oauth_idp_config" {
  name          = "oidc.oauth-idp-config%{random_suffix}"
  display_name  = "Another display name"
  client_id     = "different"
  issuer        = "different-issuer"
  enabled       = false
  client_secret = "secret2"
}
`, context)
}
