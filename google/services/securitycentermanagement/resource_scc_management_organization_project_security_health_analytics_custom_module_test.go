// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycentermanagement_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Custom Module tests cannot be run in parallel without running into 409 Conflict reponses.
// Run them as individual steps of an update test instead.
func testAccSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleBasicExample(context),
			},
			{
				ResourceName:      "google_scc_management_project_security_health_analytics_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleFullExample(context),
			},
			{
				ResourceName:      "google_scc_management_project_security_health_analytics_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleUpdate(context),
			},
			{
				ResourceName:      "google_scc_management_project_security_health_analytics_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_management_project_security_health_analytics_custom_module" "example" {
	display_name = "tf_test_basic_custom_module%{random_suffix}"
	enablement_state = "ENABLED"
	location = "%{location}"
	custom_config {
		predicate {
			expression = "resource.rotationPeriod > duration(\"2592000s\")"
		}
		resource_selector {
			resource_types = [
				"cloudkms.googleapis.com/CryptoKey",
			]
		}
		description = "The rotation period of the identified cryptokey resource exceeds 30 days."
		recommendation = "Set the rotation period to at most 30 days."
		severity = "MEDIUM"
	}
}
`, context)
}

func testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_management_project_security_health_analytics_custom_module" "example" {
	display_name = "tf_test_full_custom_module%{random_suffix}"
	enablement_state = "ENABLED"
	location = "%{location}"
	custom_config {
		predicate {
			expression = "resource.rotationPeriod > duration(\"2592000s\")"
			title = "Purpose of the expression"
			description = "description of the expression"
			location = "location of the expression"
		}
		custom_output {
			properties {
				name = "duration"
				value_expression {
					expression = "resource.rotationPeriod"
					title = "Purpose of the expression"
					description = "description of the expression"
					location = "location of the expression"
				}
			}
		}
		resource_selector {
			resource_types = [
				"cloudkms.googleapis.com/CryptoKey",
			]
		}
		severity = "LOW"
		description = "Description of the custom module"
		recommendation = "Steps to resolve violation"
	}
}
`, context)
}

func testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule_sccProjectCustomModuleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_management_project_security_health_analytics_custom_module" "example" {
	location = "%{location}"
	display_name = "full_custom_module"
	enablement_state = "DISABLED"
	custom_config {
		predicate {
			expression = "resource.name == \"updated-name\""
			title = "Updated expression title"
			description = "Updated description of the expression"
			location = "Updated location of the expression"
		}
		custom_output {
			properties {
				name = "violation"
				value_expression {
					expression = "resource.name"
					title = "Updated expression title"
					description = "Updated description of the expression"
					location = "Updated location of the expression"
				}
			}
		}
		resource_selector {
			resource_types = [
				"compute.googleapis.com/Instance",
			]
		}
		severity = "CRITICAL"
		description = "Updated description of the custom module"
		recommendation = "Updated steps to resolve violation"
	}
}
`, context)
}

func testAccCheckSecurityCenterManagementProjectSecurityHealthAnalyticsCustomModuleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_management_project_security_health_analytics_custom_module" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			location := rs.Primary.Attributes["location"]

			url, err := tpgresource.ReplaceVarsForTest(config, rs, fmt.Sprintf(
				"{{SecurityCenterBasePath}}projects/{{project}}/locations/%s/securityHealthAnalyticsCustomModules/{{name}}", location))
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
				return fmt.Errorf("SecurityCenterManagementProjectSecurityHealthAnalyticsCustomModule still exists at %s", url)
			}
		}

		return nil
	}
}
