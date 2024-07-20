// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccProviderFunction_zone_from_id(t *testing.T) {
	t.Parallel()

	zone := envvar.GetTestZoneFromEnv()
	zoneRegex := regexp.MustCompile(fmt.Sprintf("^%s$", zone))

	context := map[string]interface{}{
		"function_name": "zone_from_id",
		"output_name":   "zone",
		"resource_name": fmt.Sprintf("tf-test-zone-id-func-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Can get the zone from a resource's id in one step
				// Uses google_compute_disk resource's id attribute with format projects/{{project}}/zones/{{zone}}/disks/{{name}}
				Config: testProviderFunction_get_zone_from_resource_id(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), zoneRegex),
				),
			},
			{
				// Can get the zone from a resource's self_link in one step
				// Uses google_compute_disk resource's self_link attribute
				Config: testProviderFunction_get_zone_from_resource_self_link(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), zoneRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_zone_from_resource_id(context map[string]interface{}) string {
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
}

output "%{output_name}" {
  value = provider::google::%{function_name}(google_compute_disk.default.id)
}
`, context)
}

func testProviderFunction_get_zone_from_resource_self_link(context map[string]interface{}) string {
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
}

output "%{output_name}" {
  value = provider::google::%{function_name}(google_compute_disk.default.self_link)
}
`, context)
}
