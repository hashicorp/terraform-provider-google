// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeBackendService_basic(serviceName, checkName),
				Check:  acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_backend_service.baz", "google_compute_backend_service.foobar"),
			},
		},
	})
}

func testAccDataSourceComputeBackendService_basic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  description   = "foobar backend service"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

data "google_compute_backend_service" "baz" {
  name = google_compute_backend_service.foobar.name
}
`, serviceName, checkName)
}
