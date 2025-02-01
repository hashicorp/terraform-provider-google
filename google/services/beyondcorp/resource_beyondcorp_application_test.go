// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBeyondcorpApplication_beyondcorpSecurityGatewayApplicationBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBeyondcorpApplication_beyondcorpSecurityGatewayApplicationBasicExample_basic(context),
			},
			{
				ResourceName:            "google_beyondcorp_application.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_id", "security_gateways_id"},
			},
			{
				Config: testAccBeyondcorpApplication_beyondcorpSecurityGatewayApplicationBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_beyondcorp_application.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_beyondcorp_application.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application_id", "security_gateways_id"},
			},
		},
	})
}

func testAccBeyondcorpApplication_beyondcorpSecurityGatewayApplicationBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "default" {
  security_gateway_id = "default%{random_suffix}"
  display_name = "My Security Gateway resource"
  hubs { region = "us-central1" }
}

resource "google_beyondcorp_application" "example" {
  security_gateways_id = google_beyondcorp_security_gateway.default.security_gateway_id
  application_id = "google%{random_suffix}"
  endpoint_matchers {
    hostname = "google.com"
  }
}
`, context)
}

func testAccBeyondcorpApplication_beyondcorpSecurityGatewayApplicationBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "default" {
  security_gateway_id = "default%{random_suffix}"
  display_name = "My Security Gateway resource"
  hubs { region = "us-central1" }
}

resource "google_beyondcorp_application" "example" {
  security_gateways_id = google_beyondcorp_security_gateway.default.security_gateway_id
  display_name = "Updated Name"
  application_id = "google%{random_suffix}"
  endpoint_matchers {
    hostname = "google.com"
  }
}
`, context)
}
