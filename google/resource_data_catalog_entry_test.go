package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataCatalogEntry_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
