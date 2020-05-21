package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMonitoringService_AppEngine(t *testing.T) {
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitoringService_AppEngine(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_monitoring_app_engine_service.default", "name"),
					resource.TestCheckResourceAttrSet("data.google_monitoring_app_engine_service.default", "display_name"),
					resource.TestCheckResourceAttr(
						"data.google_monitoring_app_engine_service.default",
						"telemetry.0.resource_name",
						fmt.Sprintf("//appengine.googleapis.com/apps/%s/services/default", getTestProjectFromEnv()),
					),
				),
			},
		},
	})
}

// This does not create an app engine service - instead, it uses the
// base App Engine service "default" that cannot be deleted
func testAccDataSourceMonitoringService_AppEngine() string {
	return fmt.Sprintf(`
data "google_monitoring_app_engine_service" "default" {
	module_id = "default"
}`)
}
