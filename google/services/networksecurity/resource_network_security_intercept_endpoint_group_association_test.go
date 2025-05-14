// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkSecurityInterceptEndpointGroupAssociation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityInterceptEndpointGroupAssociation_basic(context),
			},
			{
				ResourceName:            "google_network_security_intercept_endpoint_group_association.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecurityInterceptEndpointGroupAssociation_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_network_security_intercept_endpoint_group_association.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_network_security_intercept_endpoint_group_association.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkSecurityInterceptEndpointGroupAssociation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "producer_network" {
  name                    = "tf-test-example-prod-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "consumer_network" {
  name                    = "tf-test-example-cons-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_intercept_deployment_group" "deployment_group" {
  intercept_deployment_group_id = "tf-test-example-dg%{random_suffix}"
  location                      = "global"
  network                       = google_compute_network.producer_network.id
}

resource "google_network_security_intercept_endpoint_group" "endpoint_group" {
  intercept_endpoint_group_id = "tf-test-example-eg%{random_suffix}"
  location                    = "global"
  intercept_deployment_group  = google_network_security_intercept_deployment_group.deployment_group.id
}

resource "google_network_security_intercept_endpoint_group_association" "default" {
  intercept_endpoint_group_association_id = "tf-test-example-ega%{random_suffix}"
  location                                = "global"
  network                                 = google_compute_network.consumer_network.id
  intercept_endpoint_group                = google_network_security_intercept_endpoint_group.endpoint_group.id
  labels = {
    foo = "bar"
  }
}
`, context)
}

func testAccNetworkSecurityInterceptEndpointGroupAssociation_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "producer_network" {
  name                    = "tf-test-example-prod-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "consumer_network" {
  name                    = "tf-test-example-cons-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_intercept_deployment_group" "deployment_group" {
  intercept_deployment_group_id = "tf-test-example-dg%{random_suffix}"
  location                      = "global"
  network                       = google_compute_network.producer_network.id
}

resource "google_network_security_intercept_endpoint_group" "endpoint_group" {
  intercept_endpoint_group_id = "tf-test-example-eg%{random_suffix}"
  location                    = "global"
  intercept_deployment_group  = google_network_security_intercept_deployment_group.deployment_group.id
}

resource "google_network_security_intercept_endpoint_group_association" "default" {
  intercept_endpoint_group_association_id = "tf-test-example-ega%{random_suffix}"
  location                                = "global"
  network                                 = google_compute_network.consumer_network.id
  intercept_endpoint_group                = google_network_security_intercept_endpoint_group.endpoint_group.id
  labels = {
    foo = "goo"
  }
}
`, context)
}
