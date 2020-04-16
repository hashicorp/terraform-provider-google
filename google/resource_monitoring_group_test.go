package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMonitoringGroup_update(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringGroup_update("europe-west1"),
			},
			{
				ResourceName:      "google_monitoring_group.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringGroup_update("europe-west2"),
			},
			{
				ResourceName:      "google_monitoring_group.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringGroup_update(zone string) string {
	return fmt.Sprintf(`
resource "google_monitoring_group" "update" {
  display_name = "tf-test Integration Test Group"

  filter = "resource.metadata.region=\"%s\""
}
`, zone,
	)
}
