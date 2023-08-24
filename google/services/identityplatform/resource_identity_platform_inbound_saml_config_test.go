// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIdentityPlatformInboundSamlConfig_inboundSamlConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformInboundSamlConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformInboundSamlConfig_inboundSamlConfigBasic(context),
			},
			{
				ResourceName:      "google_identity_platform_inbound_saml_config.saml_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformInboundSamlConfig_inboundSamlConfigUpdate(context),
			},
			{
				ResourceName:      "google_identity_platform_inbound_saml_config.saml_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityPlatformInboundSamlConfig_inboundSamlConfigBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_inbound_saml_config" "saml_config" {
  name = "saml.tf-config%{random_suffix}"
  display_name = "Display Name"
  idp_config {
    idp_entity_id = "tf-idp%{random_suffix}"
    sso_url = "https://example.com"
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

func testAccIdentityPlatformInboundSamlConfig_inboundSamlConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_inbound_saml_config" "saml_config" {
  name = "saml.tf-config%{random_suffix}"
  display_name = "Display Name2"
  idp_config {
    idp_entity_id = "tf-idp%{random_suffix}"
    sso_url = "https://example123.com"
    sign_request = true
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
