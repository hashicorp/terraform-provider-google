// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package discoveryengine_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccDiscoveryEngineDataStore_discoveryengineDatastoreBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDiscoveryEngineDataStoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineDataStore_discoveryengineDatastoreBasicExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_store.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_store_id", "create_advanced_site_search"},
			},
			{
				Config: testAccDiscoveryEngineDataStore_discoveryengineDatastoreBasicExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_data_store.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_store_id", "create_advanced_site_search"},
			},
		},
	})
}

func testAccDiscoveryEngineDataStore_discoveryengineDatastoreBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                    = "global"
  data_store_id               = "tf-test-data-store-id%{random_suffix}"
  display_name                = "tf-test-structured-datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
}
`, context)
}

func testAccDiscoveryEngineDataStore_discoveryengineDatastoreBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_data_store" "basic" {
  location                    = "global"
  data_store_id               = "tf-test-data-store-id%{random_suffix}"
  display_name                = "updated-tf-test-structured-datastore"
  industry_vertical           = "GENERIC"
  content_config              = "NO_CONTENT"
}
`, context)
}
