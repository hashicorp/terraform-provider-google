package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataCatalogTagTemplate_dataCatalogTagTemplate_updateFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"force_delete":  true,
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataCatalogTagTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogTagTemplate_dataCatalogTagTemplateBasicExample(context),
			},
			{
				ResourceName:            "google_data_catalog_tag_template.basic_tag_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "tag_template_id", "force_delete"},
			},
			{
				Config: testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdateFields(context),
			},
			{
				ResourceName:            "google_data_catalog_tag_template.basic_tag_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "tag_template_id", "force_delete"},
			},
			{
				Config:      testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdatePrimitiveTypeOfFieldsWithRequired(context),
				ExpectError: regexp.MustCompile("Updating the primitive type for a required field on an existing tag template is not supported"),
			},
			{
				Config: testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdatePrimitiveTypeOfFieldsWithOptional(context),
			},
			{
				ResourceName:            "google_data_catalog_tag_template.basic_tag_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "tag_template_id", "force_delete"},
			},
		},
	})
}

func testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdateFields(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_catalog_tag_template" "basic_tag_template" {
  tag_template_id = "tf_test_my_template%{random_suffix}"
  region = "us-central1"
  display_name = "Demo Tag Template Test Update"

  fields {
    field_id = "source"
    display_name = "Source of data asset test update"
    type {
      primitive_type = "STRING"
    }
    is_required = true
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
`, context)
}

func testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdatePrimitiveTypeOfFieldsWithRequired(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_catalog_tag_template" "basic_tag_template" {
  tag_template_id = "tf_test_my_template%{random_suffix}"
  region = "us-central1"
  display_name = "Demo Tag Template Test Update"

  fields {
    field_id = "source"
    display_name = "Source of data asset test update"
    type {
      primitive_type = "DOUBLE"
    }
    is_required = true
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
`, context)
}

func testAccDataCatalogTagTemplate_dataCatalogTagTemplateUpdatePrimitiveTypeOfFieldsWithOptional(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_catalog_tag_template" "basic_tag_template" {
  tag_template_id = "tf_test_my_template%{random_suffix}"
  region = "us-central1"
  display_name = "Demo Tag Template Test Update"

  fields {
    field_id = "source"
    display_name = "Source of data asset test update"
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
`, context)
}
