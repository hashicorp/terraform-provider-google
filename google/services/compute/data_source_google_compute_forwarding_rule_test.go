// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleForwardingRule(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleForwardingRuleConfig(poolName, ruleName),
				Check:  acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_compute_forwarding_rule.my_forwarding_rule", "google_compute_forwarding_rule.foobar-fr", map[string]struct{}{"port_range": {}, "target": {}}),
			},
		},
	})
}

func testAccDataSourceGoogleForwardingRuleConfig(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}

resource "google_compute_forwarding_rule" "foobar-fr" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = google_compute_target_pool.foobar-tp.self_link
}

data "google_compute_forwarding_rule" "my_forwarding_rule" {
  name = google_compute_forwarding_rule.foobar-fr.name
}
`, poolName, ruleName)
}
