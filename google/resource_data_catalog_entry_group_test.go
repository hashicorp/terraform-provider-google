package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataCatalogEntryGroup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataCatalogEntryGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogEntryGroup_dataCatalogEntryGroupBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group.basic_entry_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntryGroup_dataCatalogEntryGroupFullExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group.basic_entry_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntryGroup_dataCatalogEntryGroupBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry_group.basic_entry_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
