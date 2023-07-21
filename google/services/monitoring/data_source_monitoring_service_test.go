// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceMonitoringService_AppEngine(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitoringService_AppEngine(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_monitoring_app_engine_service.default", "name"),
					resource.TestCheckResourceAttrSet("data.google_monitoring_app_engine_service.default", "display_name"),
					resource.TestCheckResourceAttr(
						"data.google_monitoring_app_engine_service.default",
						"telemetry.0.resource_name",
						fmt.Sprintf("//appengine.googleapis.com/apps/%s/services/default", envvar.GetTestProjectFromEnv()),
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
