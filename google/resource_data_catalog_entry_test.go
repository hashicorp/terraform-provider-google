package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataCatalogEntry_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataCatalogEntryDestroyProducer(t),
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
