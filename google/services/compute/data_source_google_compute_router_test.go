// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceComputeRouter(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-router-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRouterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_router.myrouter", "id", name),
					resource.TestCheckResourceAttr("data.google_compute_router.myrouter", "name", name),
					resource.TestCheckResourceAttr("data.google_compute_router.myrouter", "network", fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", envvar.GetTestProjectFromEnv(), name)),
				),
			},
		},
	})
}

func testAccDataSourceComputeRouterConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  network = google_compute_network.foobar.name
  bgp {
    asn = 64514
  }
}

data "google_compute_router" "myrouter" {
  name    = google_compute_router.foobar.name
  network = google_compute_network.foobar.name
}
`, name, name)
}
