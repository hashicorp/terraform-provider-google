package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformTenantDefaultSupportedIdpConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigBasic(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_default_supported_idp_config.idp_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
			{
				Config: testAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigUpdate(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_default_supported_idp_config.idp_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
		},
	})
}

func testAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_default_supported_idp_config" "idp_config" {
  enabled       = true
  tenant        = google_identity_platform_tenant.tenant.name
  idp_id        = "playgames.google.com"
  client_id     = "client-id"
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_default_supported_idp_config" "idp_config" {
  enabled       = false
  tenant        = google_identity_platform_tenant.tenant.name
  idp_id        = "playgames.google.com"
  client_id     = "client-id2"
  client_secret = "differentsecret"
}
`, context)
}
