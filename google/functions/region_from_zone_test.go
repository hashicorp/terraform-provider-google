// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_region_from_zone(t *testing.T) {
	t.Parallel()
	// Skipping due to requiring TF 1.8.0 in VCR systems : https://github.com/hashicorp/terraform-provider-google/issues/17451
	acctest.SkipIfVcr(t)
	projectZone := "us-central1-a"
	projectRegion := "us-central1"
	projectRegionRegex := regexp.MustCompile(fmt.Sprintf("^%s$", projectRegion))

	context := map[string]interface{}{
		"function_name":     "region_from_zone",
		"output_name":       "zone",
		"resource_name":     fmt.Sprintf("tf-test-region-from-zone-func-%s", acctest.RandString(t, 10)),
		"resource_location": projectZone,
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testProviderFunction_get_region_from_zone(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), projectRegionRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_region_from_zone(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
		  source = "hashicorp/google"
		}
	}
}

resource "google_compute_disk" "default" {
	name  = "%{resource_name}"
	type  = "pd-ssd"
	zone  = "%{resource_location}"
	image = "debian-11-bullseye-v20220719"
	labels = {
	  environment = "dev"
	}
	physical_block_size_bytes = 4096
  }

output "%{output_name}" {
  value = provider::google::%{function_name}(google_compute_disk.default.zone)
}
`, context)
}
