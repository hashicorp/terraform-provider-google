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

// Even though the resource has generated tests, keep this one around until we are able to generate
// checking the different import formats
func TestAccComputeSubnetworkIamPolicy(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.networkUser"
	region := envvar.GetTestRegionFromEnv()
	subnetwork := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetworkIamPolicy_basic(account, region, subnetwork, role),
			},
			// Test a few import formats
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", project, region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", project, region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     subnetwork,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSubnetworkIamPolicy_basic(account, region, subnetworkName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_compute_network" "network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "%s"
  region        = "%s"
  ip_cidr_range = "10.1.0.0/16"
  network       = google_compute_network.network.name
}

data "google_iam_policy" "foo" {
  binding {
    role = "%s"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_compute_subnetwork_iam_policy" "foo" {
  project     = google_compute_subnetwork.subnetwork.project
  region      = google_compute_subnetwork.subnetwork.region
  subnetwork  = google_compute_subnetwork.subnetwork.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, account, subnetworkName, subnetworkName, region, roleId)
}
