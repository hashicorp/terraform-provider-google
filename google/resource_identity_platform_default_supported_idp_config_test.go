package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityPlatformDefaultSupportedIdpConfigDestroyProducer(t),
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

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{IdentityPlatformBasePath}}projects/{{project}}/defaultSupportedIdpConfigs/{{client_id}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("IdentityPlatformDefaultSupportedIdpConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = true
  idp_id  = "playgames.google.com"
  client_id = "client-id"
  client_secret = "secret"
}
`, context)
}

func testAccIdentityPlatformDefaultSupportedIdpConfig_defaultSupportedIdpConfigUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_identity_platform_default_supported_idp_config" "idp_config" {
  enabled = false
  idp_id  = "playgames.google.com"
  client_id = "client-id"
  client_secret = "anothersecret"
}
`, context)
}
