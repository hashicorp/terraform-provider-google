// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeHttpsHealthCheck_update(t *testing.T) {
	t.Parallel()

	hhckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHttpsHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHttpsHealthCheck_update1(hhckName),
			},
			{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeHttpsHealthCheck_update2(hhckName),
			},
			{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeHttpsHealthCheck_update1(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
  name         = "%s"
  description  = "Resource created for Terraform acceptance testing"
  request_path = "/not_default"
}
`, hhckName)
}

func testAccComputeHttpsHealthCheck_update2(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
  name                = "%s"
  description         = "Resource updated for Terraform acceptance testing"
  healthy_threshold   = 10
  unhealthy_threshold = 10
}
`, hhckName)
}
