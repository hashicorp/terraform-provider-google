// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformTenantOauthIdpConfig_identityPlatformTenantOauthIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformTenantOauthIdpConfigDestroyProducer(t),
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
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
