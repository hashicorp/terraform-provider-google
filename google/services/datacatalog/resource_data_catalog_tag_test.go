// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datacatalog_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataCatalogTag_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"force_delete":  true,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataCatalogEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogTag_dataCatalogEntryTagBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_tag.basic_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogTag_dataCatalogEntryTag_update(context),
			},
			{
				ResourceName:      "google_data_catalog_tag.basic_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogTag_dataCatalogEntryTagBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_tag.basic_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataCatalogTag_dataCatalogEntryTag_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_entry" "entry" {
  entry_group = google_data_catalog_entry_group.entry_group.id
  entry_id = "tf_test_my_entry%{random_suffix}"

  user_specified_type = "my_custom_type"
  user_specified_system = "SomethingExternal"
}

resource "google_data_catalog_entry_group" "entry_group" {
  entry_group_id = "tf_test_my_entry_group%{random_suffix}"
}

resource "google_data_catalog_tag_template" "tag_template" {
  tag_template_id = "tf_test_my_template%{random_suffix}"
  region = "us-central1"
  display_name = "Demo Tag Template"

  fields {
    field_id = "source"
    display_name = "Source of data asset"
    type {
      primitive_type = "STRING"
    }
    is_required = true
  }

  fields {
    field_id = "num_rows"
    display_name = "Number of rows in the data asset"
    type {
      primitive_type = "DOUBLE"
    }
  }

  fields {
    field_id = "pii_type"
    display_name = "PII type"
    type {
      enum_type {
        allowed_values {
          display_name = "EMAIL"
        }
        allowed_values {
          display_name = "SOCIAL SECURITY NUMBER"
        }
        allowed_values {
          display_name = "NONE"
        }
      }
    }
  }

  force_delete = "%{force_delete}"
}

resource "google_data_catalog_tag" "basic_tag" {
  parent   = google_data_catalog_entry.entry.id
  template = google_data_catalog_tag_template.tag_template.id

  fields {
    field_name   = "source"
    string_value = "my-new-string"
  }

  fields {
    field_name   = "num_rows"
    double_value = 5
  }
}
`, context)
}
