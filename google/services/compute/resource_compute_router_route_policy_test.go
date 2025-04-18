// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRouterRoutePolicy_PriorityZero(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	routePolicyName := fmt.Sprintf("route-policy-%s", acctest.RandString(t, 5))
	resourceName := "google_compute_router_route_policy.route_policy"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterRoutePolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterRoutePolicyPriorityZero(routerName, routePolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "terms.0.priority", "0"),
				),
			},
		},
	})
}

func testAccComputeRouterRoutePolicyPriorityZero(routerName, routePolicyName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "vpc" {
  name                    = "vpc-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_router" "router" {
  name    = "%[1]s"
  region  = "us-central1"
  network = google_compute_network.vpc.id
}

resource "google_compute_router_route_policy" "route_policy" {
  name    = "%[2]s"
  router  = google_compute_router.router.name
  region  = "us-central1"
  type    = "ROUTE_POLICY_TYPE_IMPORT"

  terms {
    priority = 0

    match {
      expression  = "destination == '192.168.0.0/24'"
      title       = "match-title"
      description = "test match description"
      location    = "us-central1"
    }

    actions {
      expression  = "accept()"
      title       = "actions-title"
      description = "test actions description"
      location    = "us-central1"
    }
  }
}
`, routerName, routePolicyName)
}
