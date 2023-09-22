// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecurityCenterProjectCustomModule_sccProjectCustomModuleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterProjectCustomModuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterProjectCustomModule_sccProjectCustomModuleFullExample(context),
			},
			{
				ResourceName:      "google_scc_project_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterProjectCustomModule_sccProjectCustomModuleUpdate(context),
			},
			{
				ResourceName:      "google_scc_project_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterProjectCustomModule_sccProjectCustomModuleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_project_custom_module" "example" {
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
