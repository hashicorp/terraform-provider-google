// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexAspectType_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexAspectTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexAspectType_full(context),
			},
			{
				ResourceName:            "google_dataplex_aspect_type.test_aspect_type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspect_type_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccDataplexAspectType_update(context),
			},
			{
				ResourceName:            "google_dataplex_aspect_type.test_aspect_type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspect_type_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexAspectType_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "test_aspect_type" {
  aspect_type_id = "tf-test-aspect-type%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}
`, context)
}

func testAccDataplexAspectType_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "test_aspect_type" {
  aspect_type_id = "tf-test-aspect-type%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform aspect type"
  description = "aspect type created by Terraform"
  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "updatedType",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry. This is updated."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}
`, context)
}
