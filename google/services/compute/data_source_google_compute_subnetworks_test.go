// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleSubnetworks_basic(t *testing.T) {
	t.Parallel()

	// Common resource configuration
	static_prefix := "tf-test"
	random_suffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()

	// Resource identifier used for content testing
	id := fmt.Sprintf(
		"projects/%s/regions/%s/subnetworks",
		project,
		region,
	)

	// Configuration of network resources
	network := static_prefix + "-network-" + random_suffix
	subnet_1 := static_prefix + "-subnet-1-" + random_suffix
	subnet_2 := static_prefix + "-subnet-2-" + random_suffix
	cidr_1 := "192.168.31.0/24"
	cidr_2 := "192.168.32.0/24"

	// Configuration map used in test deployment
	context := map[string]interface{}{
		"cidr_1":   cidr_1,
		"cidr_2":   cidr_2,
		"network":  network,
		"project":  project,
		"region":   region,
		"subnet_1": subnet_1,
		"subnet_2": subnet_2,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleSubnetworksConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.description"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.ip_cidr_range"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.name"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.network"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.network_self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.private_ip_google_access"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.0.self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.description"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.ip_cidr_range"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.name"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.network"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.network_self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.private_ip_google_access"),
					resource.TestCheckResourceAttrSet("data.google_compute_subnetworks.all", "subnetworks.1.self_link"),
					// Test content
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.all", "id", id),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.one", "subnetworks.0.ip_cidr_range", cidr_1),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.one", "subnetworks.0.name", subnet_1),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.one", "subnetworks.0.private_ip_google_access", "true"),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.two", "subnetworks.0.ip_cidr_range", cidr_2),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.two", "subnetworks.0.name", subnet_2),
					resource.TestCheckResourceAttr("data.google_compute_subnetworks.two", "subnetworks.0.private_ip_google_access", "false"),
				),
			},
		},
	})
}

func testAccCheckGoogleSubnetworksConfig(context map[string]interface{}) string {
	return fmt.Sprintf(`
locals {
  cidr_one   = "%s"
  cidr_two   = "%s"
  network    = "%s"
  project_id = "%s"
  region     = "%s"
  subnet_one = "%s"
  subnet_two = "%s"	
}

resource "google_compute_network" "this" {
  auto_create_subnetworks = false
  mtu                     = 1460
  name                    = local.network
  project                 = local.project_id
}

resource "google_compute_subnetwork" "subnet_one" {
  description              = "Test subnet one"
  ip_cidr_range            = local.cidr_one
  name                     = local.subnet_one
  network                  = google_compute_network.this.id
  private_ip_google_access = true
  project                  = local.project_id
  region                   = local.region
}

resource "google_compute_subnetwork" "subnet_two" {
  description              = "Test subnet two"
  ip_cidr_range            = local.cidr_two
  name                     = local.subnet_two
  network                  = google_compute_network.this.id
  private_ip_google_access = false
  project                  = local.project_id
  region                   = local.region
}

data "google_compute_subnetworks" "all" {
  filter = "network eq .*${google_compute_network.this.name}"

  depends_on = [
	google_compute_subnetwork.subnet_one,
	google_compute_subnetwork.subnet_two,
  ]
}

data "google_compute_subnetworks" "one" {
  filter = "name: ${google_compute_subnetwork.subnet_one.name}"
  region = local.region
}

data "google_compute_subnetworks" "two" {
  filter  = "ipCidrRange eq ${google_compute_subnetwork.subnet_two.ip_cidr_range}"
  project = local.project_id
  region  = local.region
}

data "google_compute_subnetworks" "no_attr" {
  depends_on = [
    google_compute_network.this,
    google_compute_subnetwork.subnet_one,
    google_compute_subnetwork.subnet_two,
  ]
}`,
		context["cidr_1"].(string),
		context["cidr_2"].(string),
		context["network"].(string),
		context["project"].(string),
		context["region"].(string),
		context["subnet_1"].(string),
		context["subnet_2"].(string),
	)
}
