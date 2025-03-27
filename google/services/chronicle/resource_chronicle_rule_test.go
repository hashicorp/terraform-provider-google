// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChronicleRule_chronicleRuleBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"chronicle_id":  envvar.GetTestChronicleInstanceIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckChronicleRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChronicleRule_chronicleRuleBasicExample_basic(context),
			},
			{
				Config: testAccChronicleRule_chronicleRuleBasicExample_update(context),
			},
		},
	})
}

func testAccChronicleRule_chronicleRuleBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_data_access_scope" "data_access_scope_test" {
 location = "us"
 instance = "%{chronicle_id}"
 data_access_scope_id = "tf-test-scope-name%{random_suffix}"
 description = "scope-description"
 allowed_data_access_labels {
   log_type = "GCP_CLOUDAUDIT"
 }
}

resource "google_chronicle_rule" "example" {
 location = "us"
 instance = "%{chronicle_id}"
 scope = resource.google_chronicle_data_access_scope.data_access_scope_test.name
 text = <<-EOT
             rule test_rule { meta: events:  $userid = $e.principal.user.userid  match: $userid over 10m condition: $e }
         EOT
}
`, context)
}

func testAccChronicleRule_chronicleRuleBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_data_access_scope" "data_access_scope_test" {
 location = "us"
 instance = "%{chronicle_id}"
 data_access_scope_id = "tf-test-scope-name%{random_suffix}"
 description = "scope-description"
 allowed_data_access_labels {
   log_type = "GCP_CLOUDAUDIT"
 }
}

resource "google_chronicle_rule" "example" {
 location = "us"
 instance = "%{chronicle_id}"
 text = <<-EOT
             rule test_rule { meta: events:  $updated_userid = $e.principal.user.userid  match: $updated_userid over 10m condition: $e }
         EOT
}
`, context)
}
