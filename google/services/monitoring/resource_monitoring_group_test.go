// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMonitoringGroup_update(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringGroupDestroyProducer(t),
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
