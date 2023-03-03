package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataprocMetastoreServiceDatasource_basic(t *testing.T) {
	t.Parallel()

	name := "tf-test-" + RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreServiceDatasource_basic(name, "DEVELOPER"),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_dataproc_metastore_service.my_metastore", "google_dataproc_metastore_service.my_metastore"),
				),
			},
		},
	})
}

func testAccDataprocMetastoreServiceDatasource_basic(name, tier string) string {
	return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
	service_id = "%s"
	location   = "us-central1"
	tier       = "%s"

	hive_metastore_config {
		version = "2.3.6"
	}
}

data "google_dataproc_metastore_service" "my_metastore" {
	service_id = google_dataproc_metastore_service.my_metastore.service_id
	location = google_dataproc_metastore_service.my_metastore.location
}
`, name, tier)
}
