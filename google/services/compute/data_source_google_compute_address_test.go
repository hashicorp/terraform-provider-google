// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestComputeAddressIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedCanonicalId string
		Config              *transport_tpg.Config
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
			Config:              &transport_tpg.Config{Project: "default-project"},
		},
		"id is address": {
			ImportId:            "test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/regions/us-east1/addresses/test-address",
			Config:              &transport_tpg.Config{Project: "default-project", Region: "us-east1"},
		},
		"id has invalid format": {
			ImportId:      "i/n/v/a/l/i/d",
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		addressId, err := compute.ParseComputeAddressId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if addressId.CanonicalId() != tc.ExpectedCanonicalId {
			t.Fatalf("bad: %s, expected canonical id to be `%s` but is `%s`", tn, tc.ExpectedCanonicalId, addressId.CanonicalId())
		}
	}
}

func TestAccDataSourceComputeAddress(t *testing.T) {
	t.Parallel()

	addressName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	rsName := "foobar"
	rsFullName := fmt.Sprintf("google_compute_address.%s", rsName)
	dsName := "my_address"
	dsFullName := fmt.Sprintf("data.google_compute_address.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataSourceComputeAddressDestroy(t, rsFullName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeAddressConfig(addressName, rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeAddressCheck(t, dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeAddressCheck(t *testing.T, data_source_name string, resource_name string) resource.TestCheckFunc {
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

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		if ds_attr["status"] != "RESERVED" {
			return fmt.Errorf("status is %s; want RESERVED", ds_attr["status"])
		}

		return nil
	}
}

func testAccCheckDataSourceComputeAddressDestroy(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_address" {
				continue
			}

			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			addressId, err := compute.ParseComputeAddressId(rs.Primary.ID, nil)
			if err != nil {
				return err
			}

			_, err = config.NewComputeClient(config.UserAgent).Addresses.Get(
				config.Project, addressId.Region, addressId.Name).Do()
			if err == nil {
				return fmt.Errorf("Address still exists")
			}
		}

		return nil
	}
}

func testAccDataSourceComputeAddressConfig(addressName, rsName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "%s" {
  name = "%s"
}

data "google_compute_address" "%s" {
  name = google_compute_address.%s.name
}
`, rsName, addressName, dsName, rsName)
}
