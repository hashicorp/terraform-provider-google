package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/discoveryengine"
)

func TestAccDiscoveryEngineSchema_discoveryengineSchemaBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineSchemaDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineSchema_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_schema.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_store_id", "location", "schema_id"},
			},
			{
				Config: testAccDiscoveryEngineSchema_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_schema.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_store_id", "location", "schema_id"},
			},
		},
	})
}

func testAccDiscoveryEngineSchema_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "tf-test-schema-ds%{random_suffix}"
  display_name                 = "tf-test-structured-datastore"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = true
}

resource "google_discovery_engine_schema" "basic" {
  location      = google_discovery_engine_data_store.basic.location
  data_store_id = google_discovery_engine_data_store.basic.data_store_id
  schema_id     = "tf-test-schema%{random_suffix}"
  json_schema   = "{\"$schema\":\"https://json-schema.org/draft/2020-12/schema\",\"type\":\"object\",\"datetime_detection\":true,\"geolocation_detection\":true}"
}
`, context)
}

func testAccDiscoveryEngineSchema_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                     = "global"
  data_store_id                = "tf-test-schema-ds%{random_suffix}"
  display_name                 = "tf-test-structured-datastore"
  industry_vertical            = "GENERIC"
  content_config               = "NO_CONTENT"
  solution_types               = ["SOLUTION_TYPE_SEARCH"]
  create_advanced_site_search  = false
  skip_default_schema_creation = true
}

resource "google_discovery_engine_schema" "basic" {
  location      = google_discovery_engine_data_store.basic.location
  data_store_id = google_discovery_engine_data_store.basic.data_store_id
  schema_id     = "tf-test-schema%{random_suffix}"
  json_schema   = "{\"$schema\":\"https://json-schema.org/draft/2020-12/schema\",\"type\":\"object\",\"datetime_detection\":true,\"geolocation_detection\":true,\"properties\":{\"title\":{\"type\":\"string\",\"retrievable\":true,\"completable\":true}}}"
}
`, context)
}
