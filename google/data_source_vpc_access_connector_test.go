package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPCAccessConnectorDatasource_basic(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnectorDatasourceConfig(RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vpc_access_connector.connector",
						"google_vpc_access_connector.connector",
						map[string]struct{}{
							// Ignore fields not returned in response
							"self_link": {},
							"region":    {},
						},
					),
				),
			},
		},
	})
}

func testAccVPCAccessConnectorDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "vpc-con-test-%s"
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
  region        = "us-central1"
}

data "google_vpc_access_connector" "connector" {
  name = google_vpc_access_connector.connector.name
}
`, suffix)
}
