package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityPlatformTenantInboundSamlConfig_identityPlatformTenantInboundSamlConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformTenantInboundSamlConfigDestroy,
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
	return Nprintf(`
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
    sso_url       = "example.com"
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
	return Nprintf(`
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
    sso_url       = "example123.com"
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
