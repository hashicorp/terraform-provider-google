// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformTenantDefaultSupportedIdpConfig_identityPlatformTenantDefaultSupportedIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformTenantDefaultSupportedIdpConfigDestroyProducer(t),
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
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
