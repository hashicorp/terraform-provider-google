package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformTenantOauthIdpConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigBasic(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_oauth_idp_config.tenant_oauth_idp_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
			{
				Config: testAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigUpdate(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_oauth_idp_config.tenant_oauth_idp_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
		},
	})
}

func testAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_oauth_idp_config" "tenant_oauth_idp_config" {
  name          = "oidc.oauth-idp-config%{random_suffix}"
  tenant        = google_identity_platform_tenant.tenant.name
  display_name  = "Display Name"
  client_id     = "client-id"
  issuer        = "issuer"
  enabled       = true
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_oauth_idp_config" "tenant_oauth_idp_config" {
  name          = "oidc.oauth-idp-config%{random_suffix}"
  tenant        = google_identity_platform_tenant.tenant.name
  display_name  = "My Display"
  client_id     = "different-client"
  issuer        = "issuer2"
  enabled       = false
  client_secret = "supersecret"
}
`, context)
}
