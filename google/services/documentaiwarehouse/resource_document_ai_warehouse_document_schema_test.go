// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package documentaiwarehouse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseFull(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckDocumentAIWarehouseDocumentSchemaDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseInit(context),
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaTextExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_text",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaIntegerExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_integer",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaFloatExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_float",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaPropertyExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_property",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaPropertyEnumExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_property_enum",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaEnumExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_enum",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaMapExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_map",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaDatetimeExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_datetime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
			{
				Config: testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaTimestampExample(context),
			},
			{
				ResourceName:            "google_document_ai_warehouse_document_schema.example_timestamp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_number", "location"},
			},
		},
	})
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseInit(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "contentwarehouse" {
  project = google_project.project.project_id
  service = "contentwarehouse.googleapis.com"
  disable_on_destroy = false
}

resource "time_sleep" "wait_120s" {
  create_duration = "120s"

  depends_on = [google_project_service.contentwarehouse]
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"

  depends_on = [time_sleep.wait_120s]
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaTextExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_text" {
  project_number     = google_project.project.number
  display_name       = "test-property-text"
  location           = "us"
  document_is_folder = false

  property_definitions {
    name                 = "prop3"
    display_name         = "propdisp3"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    text_type_options {}
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaIntegerExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_integer" {
  project_number = google_project.project.number
  display_name   = "test-property-integer"
  location       = "us"

  property_definitions {
    name                 = "prop1"
    display_name         = "propdisp1"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    integer_type_options {}
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaFloatExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_float" {
  project_number = google_project.project.number
  display_name   = "test-property-float"
  location       = "us"

  property_definitions {
    name                 = "prop2"
    display_name         = "propdisp2"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    float_type_options {}
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaPropertyExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_property" {
  project_number     = google_project.project.number
  display_name       = "test-property-property"
  location           = "us"
  document_is_folder = false

  property_definitions {
    name                 = "prop8"
    display_name         = "propdisp8"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    property_type_options {
      property_definitions {
        name                 = "prop8_nested"
        display_name         = "propdisp8_nested"
        is_repeatable        = false
        is_filterable        = true
        is_searchable        = true
        is_metadata          = false
        is_required          = false
        retrieval_importance = "HIGHEST"
        schema_sources {
          name           = "dummy_source_nested"
          processor_type = "dummy_processor_nested"
        }
        text_type_options {}
      }
    }
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaPropertyEnumExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_property_enum" {
  project_number     = google_project.project.number
  display_name       = "test-property-property"
  location           = "us"
  document_is_folder = false

  property_definitions {
    name                 = "prop8"
    display_name         = "propdisp8"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    property_type_options {
      property_definitions {
        name                 = "prop8_nested"
        display_name         = "propdisp8_nested"
        is_repeatable        = false
        is_filterable        = true
        is_searchable        = true
        is_metadata          = false
        is_required          = false
        retrieval_importance = "HIGHEST"
        schema_sources {
          name           = "dummy_source_nested"
          processor_type = "dummy_processor_nested"
        }
        enum_type_options {
          possible_values = [
            "M",
            "F",
            "X"
          ]
          validation_check_disabled = false
        }
      }
    }
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaEnumExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_enum" {
  project_number = google_project.project.number
  display_name   = "test-property-enum"
  location       = "us"

  property_definitions {
    name                 = "prop6"
    display_name         = "propdisp6"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    enum_type_options {
      possible_values = [
        "M",
        "F",
        "X"
      ]
      validation_check_disabled = false
    }
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaMapExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_map" {
  project_number = google_project.project.number
  display_name   = "test-property-map"
  location       = "us"

  property_definitions {
    name                 = "prop4"
    display_name         = "propdisp4"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    map_type_options {}
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaDatetimeExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_datetime" {
  project_number = google_project.project.number
  display_name   = "test-property-date_time"
  location       = "us"

  property_definitions {
    name                 = "prop7"
    display_name         = "propdisp7"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    date_time_type_options {}
  }
}
`, context)
}

func testAccDocumentAIWarehouseDocumentSchema_documentAiWarehouseDocumentSchemaTimestampExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-%{random_suffix}"
  name            = "tf-test-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_document_ai_warehouse_location" "loc" {
  location = "us"
  project_number = google_project.project.number
  access_control_mode = "ACL_MODE_DOCUMENT_LEVEL_ACCESS_CONTROL_GCI"
  database_type = "DB_INFRA_SPANNER"
  document_creator_default_role = "DOCUMENT_ADMIN"
}

resource "google_document_ai_warehouse_document_schema" "example_timestamp" {
  project_number = google_project.project.number
  display_name   = "test-property-timestamp"
  location       = "us"

  property_definitions {
    name                 = "prop5"
    display_name         = "propdisp5"
    is_repeatable        = false
    is_filterable        = true
    is_searchable        = true
    is_metadata          = false
    is_required          = false
    retrieval_importance = "HIGHEST"
    schema_sources {
      name           = "dummy_source"
      processor_type = "dummy_processor"
    }
    timestamp_type_options {}
  }
}
`, context)
}

func testAccCheckDocumentAIWarehouseDocumentSchemaDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_document_ai_warehouse_document_schema" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DocumentAIWarehouseBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DocumentAIWarehouseDocumentSchema still exists at %s", url)
			}
		}

		return nil
	}
}
