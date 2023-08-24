// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformOauthIdpConfig_identityPlatformOauthIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformOauthIdpConfigDestroyProducer(t),
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
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
