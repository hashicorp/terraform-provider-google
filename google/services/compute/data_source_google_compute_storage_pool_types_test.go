// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStoragePoolTypes(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStoragePoolTypes("us-central1-a"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleStoragePoolTypesCheck(envvar.GetTestProjectFromEnv(), "data.google_compute_storage_pool_types.balanced"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStoragePoolTypesCheck(projectID, data_source_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		ds_attr := ds.Primary.Attributes
		ver := "v1"

		expected := map[string]string{
			"kind":                             "compute#storagePoolType",
			"id":                               "1002",
			"creation_timestamp":               "1969-12-31T16:00:00.000-08:00",
			"name":                             "hyperdisk-balanced",
			"description":                      "Hyperdisk Balanced Storage Pool",
			"zone":                             fmt.Sprintf("https://www.googleapis.com/compute/%s/projects/%s/zones/us-central1-a", ver, projectID),
			"self_link":                        fmt.Sprintf("https://www.googleapis.com/compute/%s/projects/%s/zones/us-central1-a/storagePoolTypes/hyperdisk-balanced", ver, projectID),
			"self_link_with_id":                fmt.Sprintf("https://www.googleapis.com/compute/%s/projects/%s/zones/us-central1-a/storagePoolTypes/1002", ver, projectID),
			"min_pool_provisioned_capacity_gb": "10240",
			"max_pool_provisioned_capacity_gb": "5242880",
			"min_pool_provisioned_iops":        "0",
			"max_pool_provisioned_iops":        "4190000",
			"min_pool_provisioned_throughput":  "0",
			"max_pool_provisioned_throughput":  "1048576",
		}

		for k, v := range expected {
			if ds_attr[k] != v {
				return fmt.Errorf(
					"%s is %s; want %s",
					k,
					ds_attr[k],
					v,
				)
			}
		}

		expectedDiskType := fmt.Sprintf("https://www.googleapis.com/compute/%s/projects/%s/zones/us-central1-a/diskTypes/hyperdisk-balanced", ver, projectID)

		if len(ds_attr["supported_disk_types.#"]) == 0 {
			return fmt.Errorf("supported_disk_types is empty")
		} else if ds_attr["supported_disk_types.0"] != expectedDiskType {
			return fmt.Errorf(
				"supported_disk_types is %s; want %s", ds_attr["supported_disk_types.0"], expectedDiskType)
		}

		return nil
	}
}

func testAccDataSourceGoogleStoragePoolTypes(zone string) string {
	return fmt.Sprintf(`
data "google_compute_storage_pool_types" "balanced" {
  zone = "%s"
	storage_pool_type = "hyperdisk-balanced"
}
`, zone)
}
