// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBeyondcorpSecurityGateway_beyondcorpSecurityGatewayBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBeyondcorpSecurityGateway_beyondcorpSecurityGatewayBasicExample_basic(context),
			},
			{
				ResourceName:            "google_beyondcorp_security_gateway.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "security_gateway_id"},
			},
			{
				Config: testAccBeyondcorpSecurityGateway_beyondcorpSecurityGatewayBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_beyondcorp_security_gateway.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_beyondcorp_security_gateway.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "security_gateway_id"},
			},
		},
	})
}

func testAccBeyondcorpSecurityGateway_beyondcorpSecurityGatewayBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "example" {
  security_gateway_id = "default%{random_suffix}"
  location = "global"
  display_name = "My Security Gateway resource"
  hubs { region = "us-central1" }
}
`, context)
}

func testAccBeyondcorpSecurityGateway_beyondcorpSecurityGatewayBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "example" {
  security_gateway_id = "default%{random_suffix}"
  location = "global"
  display_name = "My Security Gateway resource"
  hubs { region = "us-east1" }
}
`, context)
}
