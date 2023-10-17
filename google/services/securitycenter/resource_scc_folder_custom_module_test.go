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

// Custom Module tests cannot be run in parallel without running into 409 Conflict reponses.
// Run them as individual steps of an update test instead.
func TestAccSecurityCenterFolderCustomModule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"sleep":         true,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckSecurityCenterFolderCustomModuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleBasicExample(context),
			},
			{
				ResourceName:            "google_scc_folder_custom_module.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleFullExample(context),
			},
			{
				ResourceName:            "google_scc_folder_custom_module.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleUpdate(context),
			},
			{
				ResourceName:            "google_scc_folder_custom_module.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
		},
	})
}

func testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
}

resource "time_sleep" "wait_1_minute" {
	depends_on = [google_folder.folder]

	create_duration = "2m"
}

resource "google_scc_folder_custom_module" "example" {
	folder = google_folder.folder.folder_id
	display_name = "tf_test_basic_custom_module%{random_suffix}"
	enablement_state = "ENABLED"
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


	depends_on = [time_sleep.wait_1_minute]
}
`, context)
}

func testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
}

resource "google_scc_folder_custom_module" "example" {
	folder = google_folder.folder.folder_id
	display_name = "tf_test_full_custom_module%{random_suffix}"
	enablement_state = "ENABLED"
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

func testAccSecurityCenterFolderCustomModule_sccFolderCustomModuleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
}

resource "google_scc_folder_custom_module" "example" {
	folder = google_folder.folder.folder_id
	display_name = "tf_test_full_custom_module%{random_suffix}"
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

func testAccCheckSecurityCenterFolderCustomModuleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_folder_custom_module" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{SecurityCenterBasePath}}folders/{{folder}}/securityHealthAnalyticsSettings/customModules/{{name}}")
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
				return fmt.Errorf("SecurityCenterFolderCustomModule still exists at %s", url)
			}
		}

		return nil
	}
}
