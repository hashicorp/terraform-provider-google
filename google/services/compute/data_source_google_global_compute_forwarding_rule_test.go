// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleGlobalForwardingRule(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleGlobalForwardingRuleConfig(poolName, ruleName),
				Check:  acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_compute_global_forwarding_rule.my_forwarding_rule", "google_compute_global_forwarding_rule.foobar-fr", map[string]struct{}{"port_range": {}, "target": {}}),
			},
		},
	})
}

func testAccDataSourceGoogleGlobalForwardingRuleConfig(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "foobar-fr" {
  name       = "%s"
  target     = google_compute_target_http_proxy.default.id
  port_range = "80"
}

resource "google_compute_target_http_proxy" "default" {
  name        = "%s"
  description = "a description"
  url_map     = google_compute_url_map.default.id
}

resource "google_compute_url_map" "default" {
  name            = "%s"
  default_url_redirect {
	https_redirect         = true
	redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
	strip_query            = false
  }
}
  
data "google_compute_global_forwarding_rule" "my_forwarding_rule" {
  name = google_compute_global_forwarding_rule.foobar-fr.name
}
`, ruleName, poolName, poolName)
}
