// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_location_from_id(t *testing.T) {
	t.Parallel()

	location := "us-central1"
	locationRegex := regexp.MustCompile(fmt.Sprintf("^%s$", location))

	context := map[string]interface{}{
		"function_name":     "location_from_id",
		"output_name":       "location",
		"resource_name":     fmt.Sprintf("tf-test-location-id-func-%s", acctest.RandString(t, 10)),
		"resource_location": location,
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Can get the location from a resource's id in one step
				// Uses google_cloud_run_service resource's id attribute with format projects/{project}/locations/{location}/services/{service}.
				Config: testProviderFunction_get_location_from_resource_id(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), locationRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_location_from_resource_id(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

resource "google_cloud_run_service" "default" {
  name     = "%{resource_name}"
  location = "%{resource_location}"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

output "%{output_name}" {
  value = provider::google::%{function_name}(google_cloud_run_service.default.id)
}
`, context)
}
