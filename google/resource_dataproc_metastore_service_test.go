package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataprocMetastoreService_updateAndImport(t *testing.T) {
	t.Parallel()

	name := "tf-metastore-" + randString(t, 10)
	tier := [2]string{"DEVELOPER", "ENTERPRISE"}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryReservationReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[0]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[1]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataprocMetastoreService_updateAndImport(name, tier string) string {
	return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
	service_id = "%s"
	location   = "us-central1"
	tier       = "%s"

	hive_metastore_config {
		version = "2.3.6"
	}
}
`, name, tier)
}
