package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitoringService_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_cloudEndpoints(randomSuffix, "an-endpoint"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringSlo_cloudEndpoints(randomSuffix, "another-endpoint"),
			},
			{
				ResourceName:      "google_monitoring_service.srv",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringSlo_cloudEndpoints(randSuffix, endpoint string) string {
	return fmt.Sprintf(`
resource "google_monitoring_service" "srv" {
	service_id = "tf-test-srv-%s"
	display_name = "My Basic CloudEnpoints Service"
	basic_service {
		service_type  = "CLOUD_ENDPOINTS"
		service_labels = {
			service = "%s"
		}
	}
}
`, randSuffix, endpoint)
}
