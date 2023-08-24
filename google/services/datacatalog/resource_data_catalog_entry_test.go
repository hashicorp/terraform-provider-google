// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datacatalog_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataCatalogEntry_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataCatalogEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryFullExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
