// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMonitoringService_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringServiceDestroyProducer(t),
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
