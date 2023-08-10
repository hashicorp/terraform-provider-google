// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIdentityPlatformDefaultSupportedIdpConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigBasic(context),
			},
			{
				ResourceName:      "google_identity_platform_default_supported_idp_config.idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(context),
			},
			{
				ResourceName:      "google_identity_platform_default_supported_idp_config.idp_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIdentityPlatformDefaultSupportedIdpConfigDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_identity_platform_default_supported_idp_config" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{IdentityPlatformBasePath}}projects/{{project}}/defaultSupportedIdpConfigs/{{client_id}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("IdentityPlatformDefaultSupportedIdpConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = true
  idp_id  = "playgames.google.com"
  client_id = "client-id"
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = false
  idp_id  = "playgames.google.com"
  client_id = "client-id"
  client_secret = "anothersecret"
}
`, context)
}
