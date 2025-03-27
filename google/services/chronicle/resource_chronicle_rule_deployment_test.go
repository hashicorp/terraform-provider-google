// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChronicleRuleDeployment_chronicleRuledeploymentBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"chronicle_id":  envvar.GetTestChronicleInstanceIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChronicleRuleDeployment_chronicleRuledeploymentBasicExample_basic(context),
			},
			{
				ResourceName:            "google_chronicle_rule_deployment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "location", "rule"},
			},
			{
				Config: testAccChronicleRuleDeployment_chronicleRuledeploymentBasicExample_update(context),
			},
			{
				ResourceName:            "google_chronicle_rule_deployment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "location", "rule"},
			},
		},
	})
}

func testAccChronicleRuleDeployment_chronicleRuledeploymentBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_rule" "my-rule" {
 location = "us"
 instance = "%{chronicle_id}"
 text = <<-EOT
             rule test_rule { meta: events:  $userid = $e.principal.user.userid  match: $userid over 10m condition: $e }
         EOT
}

resource "google_chronicle_rule_deployment" "example" {
 location = "us"
 instance = "%{chronicle_id}"
 rule = element(split("/", resource.google_chronicle_rule.my-rule.name), length(split("/", resource.google_chronicle_rule.my-rule.name)) - 1)
 enabled = true
 alerting = true
 archived = false
 run_frequency = "DAILY"
}
`, context)
}

func testAccChronicleRuleDeployment_chronicleRuledeploymentBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_rule" "my-rule" {
 location = "us"
 instance = "%{chronicle_id}"
 text = <<-EOT
             rule test_rule { meta: events:  $userid = $e.principal.user.userid  match: $userid over 10m condition: $e }
         EOT
}

resource "google_chronicle_rule_deployment" "example" {
 location = "us"
 instance = "%{chronicle_id}"
 rule = element(split("/", resource.google_chronicle_rule.my-rule.name), length(split("/", resource.google_chronicle_rule.my-rule.name)) - 1)
 enabled = false
 alerting = false
 archived = false
 run_frequency = "HOURLY"
}
`, context)
}
