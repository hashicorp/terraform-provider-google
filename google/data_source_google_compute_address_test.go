package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestComputeAddressIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedCanonicalId string
		Config              *Config
	}{
		"id is a full self link": {
			ImportId:            "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/addresses/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is a partial self link": {
			ImportId:            "projects/test-project/regions/us-central1/addresses/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is project/region/address": {
			ImportId:            "test-project/us-central1/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is region/address": {
			ImportId:            "us-central1/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/regions/us-central1/addresses/test-address",
			Config:              &Config{Project: "default-project"},
		},
		"id is address": {
			ImportId:            "test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/regions/us-east1/addresses/test-address",
			Config:              &Config{Project: "default-project", Region: "us-east1"},
		},
		"id has invalid format": {
			ImportId:      "i/n/v/a/l/i/d",
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		addressId, err := parseComputeAddressId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if addressId.canonicalId() != tc.ExpectedCanonicalId {
			t.Fatalf("bad: %s, expected canonical id to be `%s` but is `%s`", tn, tc.ExpectedCanonicalId, addressId.canonicalId())
		}
	}
}

func TestAccDataSourceComputeAddress(t *testing.T) {
	t.Parallel()

	rsName := "foobar"
	rsFullName := fmt.Sprintf("google_compute_address.%s", rsName)
	dsName := "my_address"
	dsFullName := fmt.Sprintf("data.google_compute_address.%s", dsName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceComputeAddressDestroy(rsFullName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeAddressConfig(rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeAddressCheck(dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeAddressCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		address_attrs_to_test := []string{
			"name",
			"address",
		}

		for _, attr_to_check := range address_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !compareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		if ds_attr["status"] != "RESERVED" {
			return fmt.Errorf("status is %s; want RESERVED", ds_attr["status"])
		}

		return nil
	}
}

func testAccCheckDataSourceComputeAddressDestroy(resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		addressId, err := parseComputeAddressId(rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		_, err = config.clientCompute.Addresses.Get(
			config.Project, addressId.Region, addressId.Name).Do()
		if err == nil {
			return fmt.Errorf("Address still exists")
		}

		return nil
	}
}

func testAccDataSourceComputeAddressConfig(rsName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "%s" {
  name = "address-test"
}

data "google_compute_address" "%s" {
  name = google_compute_address.%s.name
}
`, rsName, dsName, rsName)
}
