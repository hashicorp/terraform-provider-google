// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package discoveryengine_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccDiscoveryEngineChatEngine_discoveryengineChatengine_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineChatEngine_discoveryengineChatengine_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_chat_engine.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"chat_engine_config"},
			},
			{
				Config: testAccDiscoveryEngineChatEngine_discoveryengineChatengine_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_chat_engine.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"chat_engine_config"},
			},
		},
	})
}

func testAccDiscoveryEngineChatEngine_discoveryengineChatengine_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_discovery_engine_data_store" "test_data_store" {
		location                    = "eu"
		data_store_id               = "tf-test-data-store-id%{random_suffix}"
		display_name                = "tf-test-structured-datastore"
		industry_vertical           = "GENERIC"
		content_config              = "NO_CONTENT"
		solution_types              = ["SOLUTION_TYPE_CHAT"]
	}

	resource "google_discovery_engine_chat_engine" "primary" {
		engine_id = "tf-test-chat-engine-id%{random_suffix}"
		collection_id = "default_collection"
		location = google_discovery_engine_data_store.test_data_store.location
		display_name = "tf-test-chat-engine-name%{random_suffix}"
		industry_vertical = "GENERIC"
		data_store_ids = [google_discovery_engine_data_store.test_data_store.data_store_id]
		common_config {
		  company_name = "test-company"
		}
		chat_engine_config {
		  agent_creation_config {
			business = "test business name"
			default_language_code = "en"
			time_zone = "America/Los_Angeles"
		  }
		}
	}
	`, context)
}

func testAccDiscoveryEngineChatEngine_discoveryengineChatengine_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_discovery_engine_data_store" "test_data_store" {
		location                    = "eu"
		data_store_id               = "tf-test-data-store-id%{random_suffix}"
		display_name                = "tf-test-structured-datastore"
		industry_vertical           = "GENERIC"
		content_config              = "NO_CONTENT"
		solution_types              = ["SOLUTION_TYPE_CHAT"]
	}


	resource "google_discovery_engine_chat_engine" "primary" {
		engine_id = "tf-test-chat-engine-id%{random_suffix}"
		collection_id = "default_collection"
		location = google_discovery_engine_data_store.test_data_store.location
		display_name = "tf-test-chat-engine-name-2%{random_suffix}"
		industry_vertical = "GENERIC"
		data_store_ids = [google_discovery_engine_data_store.test_data_store.data_store_id]

		common_config {
		  company_name = "test-company"
		}
		chat_engine_config {
		  agent_creation_config {
			business = "test business name"
			default_language_code = "en"
			time_zone = "America/Los_Angeles"
		  }
		}
	}
	`, context)
}
