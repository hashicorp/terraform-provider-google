package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGoogleSslPolicy(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSslPolicy(),
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

func testAccDataSourceGoogleSslPolicy() string {
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
`, acctest.RandomWithPrefix("test-ssl-policy"))
}
