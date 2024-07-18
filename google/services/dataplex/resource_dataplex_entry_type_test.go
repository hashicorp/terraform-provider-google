// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexEntryType_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntryType_full(context),
			},
			{
				ResourceName:            "google_dataplex_entry_type.test_entry_type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entry_type_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccDataplexEntryType_update(context),
			},
			{
				ResourceName:            "google_dataplex_entry_type.test_entry_type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entry_type_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexEntryType_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_type" "test_entry_type" {
  entry_type_id = "tf-test-entry-type%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"
}
`, context)
}

func testAccDataplexEntryType_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "test_entry_type" {
  aspect_type_id         = "tf-test-aspect-type%{random_suffix}"
  location     = "us-central1"
  project      = "%{project_name}"

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

resource "google_dataplex_entry_type" "test_entry_type" {
  entry_type_id = "tf-test-entry-type-%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"
  
  required_aspects {
    type = google_dataplex_aspect_type.test_entry_type.name
  }
}
`, context)
}
