// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineSearchEngineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_search_engine.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location"},
			},
			{
				Config: testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_search_engine.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "collection_id", "location"},
			},
		},
	})
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
    location                    = "global"
    data_store_id               = "tf-test-example-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_data_store" "second" {
    location                    = "global"
    data_store_id               = "tf-test-example2-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore2"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.basic.location
  display_name = "Example Display Name"
  data_store_ids = [google_discovery_engine_data_store.basic.data_store_id, google_discovery_engine_data_store.second.data_store_id]
  industry_vertical = google_discovery_engine_data_store.basic.industry_vertical
  common_config {
    company_name = "Example Company Name"
  }
  search_engine_config {
    search_tier = "SEARCH_TIER_ENTERPRISE"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
}
`, context)
}

func testAccDiscoveryEngineSearchEngine_discoveryengineSearchengineBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
    location                    = "global"
    data_store_id               = "tf-test-example-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_data_store" "second" {
    location                    = "global"
    data_store_id               = "tf-test-example2-datastore%{random_suffix}"
    display_name                = "tf-test-structured-datastore2"
    industry_vertical           = "GENERIC"
    content_config              = "NO_CONTENT"
    solution_types              = ["SOLUTION_TYPE_SEARCH"]
    create_advanced_site_search = false
    }
resource "google_discovery_engine_search_engine" "basic" {
  engine_id = "tf-test-example-engine-id%{random_suffix}"
  collection_id = "default_collection"
  location = google_discovery_engine_data_store.basic.location
  display_name = "Updated Example Display Name"
  data_store_ids = [google_discovery_engine_data_store.basic.data_store_id]
  industry_vertical = google_discovery_engine_data_store.basic.industry_vertical
  common_config {
    company_name = "Updated Example Company Name"
  }
  search_engine_config {
    search_tier = "SEARCH_TIER_STANDARD"
    search_add_ons = ["SEARCH_ADD_ON_LLM"]
  }
}
`, context)
}
