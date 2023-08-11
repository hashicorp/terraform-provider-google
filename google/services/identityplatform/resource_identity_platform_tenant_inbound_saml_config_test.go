// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformTenantInboundSamlConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigBasic(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_inbound_saml_config.tenant_saml_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
			{
				Config: testAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigUpdate(context),
			},
			{
				ResourceName:            "google_identity_platform_tenant_inbound_saml_config.tenant_saml_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"},
			},
		},
	})
}

func testAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_inbound_saml_config" "tenant_saml_config" {
  name         = "saml.tf-config%{random_suffix}"
  display_name = "Display Name"
  tenant       = google_identity_platform_tenant.tenant.name
  idp_config {
    idp_entity_id = "tf-idp%{random_suffix}"
    sign_request  = true
    sso_url       = "https://example.com"
    idp_certificates {
      x509_certificate = file("test-fixtures/rsa_cert.pem")
    }
  }

  sp_config {
    sp_entity_id = "tf-sp%{random_suffix}"
    callback_uri = "https://example.com"
  }
}
`, context)
}

func testAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_tenant" "tenant" {
  display_name  = "tenant"
}

resource "google_identity_platform_tenant_inbound_saml_config" "tenant_saml_config" {
  name         = "saml.tf-config%{random_suffix}"
  display_name = "Display Name2"
  tenant       = google_identity_platform_tenant.tenant.name
  idp_config {
    idp_entity_id = "tf-idp%{random_suffix}"
    sign_request  = false
    sso_url       = "https://example123.com"
    idp_certificates {
      x509_certificate = file("test-fixtures/rsa_cert.pem")
    }
  }

  sp_config {
    sp_entity_id = "tf-sp%{random_suffix}"
    callback_uri = "https://example123.com"
  }
}
`, context)
}
