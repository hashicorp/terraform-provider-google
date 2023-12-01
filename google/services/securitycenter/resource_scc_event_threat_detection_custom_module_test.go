// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccSecurityCenterEventThreatDetectionCustomModule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSecurityCenterEventThreatDetectionCustomModuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterEventThreatDetectionCustomModule_sccOrganizationCustomModuleExample(context),
			},
			{
				ResourceName:            "google_scc_event_threat_detection_custom_module.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
			{
				Config: testAccSecurityCenterEventThreatDetectionCustomModule_sccOrganizationCustomModuleUpdate(context),
			},
			{
				ResourceName:            "google_scc_event_threat_detection_custom_module.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
		},
	})
}

func testAccSecurityCenterEventThreatDetectionCustomModule_sccOrganizationCustomModuleExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_event_threat_detection_custom_module" "example" {
	organization = "%{org_id}"
	display_name = "tf_test_custom_module%{random_suffix}"
	enablement_state = "ENABLED"
	type="CONFIGURABLE_BAD_IP"
	config = <<EOF
              {"metadata": {
				"severity": "LOW",
				"description": "Flagged by Forcepoint as malicious",
				"recommendation": "Contact the owner of the relevant project."
			  },
			  "ips": [
				"192.0.2.1",
				"192.0.2.0/24"
			  ]}
            EOF
}
`, context)
}

func testAccSecurityCenterEventThreatDetectionCustomModule_sccOrganizationCustomModuleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_event_threat_detection_custom_module" "example" {
	organization = "%{org_id}"
	display_name = "tf_test_custom_module%{random_suffix}_updated"
	enablement_state = "DISABLED"
	type="CONFIGURABLE_BAD_IP"
	config = <<EOF
              {"metadata": {
				"severity": "MEDIUM",
				"description": "Flagged by Forcepoint as malicious",
				"recommendation": "Contact the owner of the relevant project."
			  },
			  "ips": [
				"192.0.2.1"
			  ]}
            EOF
}
`, context)
}

func testAccSecurityCenterEventThreatDetectionCustomModuleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_event_threat_detection_custom_module" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{SecurityCenterBasePath}}organizations/{{organization}}/eventThreatDetectionSettings/customModules/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("EventThreatDetectionCustomModule still exists at %s", url)
			}
		}

		return nil
	}
}
