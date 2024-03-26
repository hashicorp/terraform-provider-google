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

func TestAccProviderFunction_name_from_id(t *testing.T) {
	t.Parallel()
	// Skipping due to requiring TF 1.8.0 in VCR systems : https://github.com/hashicorp/terraform-provider-google/issues/17451
	acctest.SkipIfVcr(t)

	context := map[string]interface{}{
		"function_name": "name_from_id",
		"output_name":   "name",
		"resource_name": fmt.Sprintf("tf-test-name-id-func-%s", acctest.RandString(t, 10)),
	}

	nameRegex := regexp.MustCompile(fmt.Sprintf("^%s$", context["resource_name"]))

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Can get the name from a resource's id in one step
				// Uses google_pubsub_topic resource's id attribute with format projects/{{project}}/topics/{{name}}
				Config: testProviderFunction_get_name_from_resource_id(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), nameRegex),
				),
			},
			{
				// Can get the name from a resource's self_link in one step
				// Uses google_compute_disk resource's self_link attribute
				Config: testProviderFunction_get_name_from_resource_self_link(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), nameRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_name_from_resource_id(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

resource "google_pubsub_topic" "default" {
  name = "%{resource_name}"
}

output "%{output_name}" {
  value = provider::google::%{function_name}(google_pubsub_topic.default.id)
}
`, context)
}

func testProviderFunction_get_name_from_resource_self_link(context map[string]interface{}) string {
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
