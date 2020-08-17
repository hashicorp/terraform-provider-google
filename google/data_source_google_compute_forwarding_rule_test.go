package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGoogleForwardingRule(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", randString(t, 10))
	ruleName := fmt.Sprintf("tf-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleForwardingRuleConfig(poolName, ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleForwardingRuleCheck("data.google_compute_forwarding_rule.my_forwarding_rule", "google_compute_forwarding_rule.foobar-fr"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleForwardingRuleCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		forwarding_rule_attrs_to_test := []string{
			"id",
			"name",
			"description",
			"region",
			"port_range",
			"ports",
			"ip_address",
			"ip_protocol",
			"load_balancing_scheme",
			"backend_service",
			"network",
			"subnetwork",
		}

		for _, attr_to_check := range forwarding_rule_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !compareSelfLinkOrResourceName("", ds_attr["target"], rs_attr["target"], nil) && ds_attr["target"] != rs_attr["target"] {
			return fmt.Errorf("target does not match: %s vs %s", ds_attr["target"], rs_attr["target"])
		}

		if !compareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		return nil
	}
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
