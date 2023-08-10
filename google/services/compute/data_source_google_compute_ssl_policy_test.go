// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleSslPolicy(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSslPolicy(fmt.Sprintf("test-ssl-policy-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSslPolicyCheck("data.google_compute_ssl_policy.ssl_policy", "google_compute_ssl_policy.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSslPolicyCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		ssl_policy_attrs_to_test := []string{
			"id",
			"self_link",
			"name",
			"description",
			"min_tls_version",
			"profile",
			"custom_features",
		}

		for _, attr_to_check := range ssl_policy_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceGoogleSslPolicy(policyName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "foobar" {
  name            = "%s"
  description     = "my-description"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
}

data "google_compute_ssl_policy" "ssl_policy" {
  name = google_compute_ssl_policy.foobar.name
}
`, policyName)
}
