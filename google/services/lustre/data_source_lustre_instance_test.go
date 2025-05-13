// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package lustre_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLustreInstanceDatasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstanceDatasource_basic(context),
				Check: acctest.CheckDataSourceStateMatchesResourceState(
					"data.google_lustre_instance.default",
					"google_lustre_instance.instance",
				),
			},
			{
				ResourceName:      "google_lustre_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLustreInstanceDatasource_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id             = "my-instance-%{random_suffix}"
  location                = "us-central1-a"
  filesystem              = "testfs"
  capacity_gib            = 18000
  network                 = google_compute_network.producer_net.id
  gke_support_enabled     = false
       
 depends_on               = [ google_service_networking_connection.service_con ]
}

resource "google_compute_subnetwork" "producer_subnet" {
  name                    = "tf-test-my-subnet-%{random_suffix}"
  ip_cidr_range           = "10.0.0.248/29"
  region                  = "us-central1"
  network                 = google_compute_network.producer_net.id
}

resource "google_compute_network" "producer_net" {
  name                    = "tf-test-my-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_global_address" "private_ip_alloc" {
  name                    = "private-ip-alloc-%{random_suffix}"
  purpose                 = "VPC_PEERING"
  address_type            = "INTERNAL"
  prefix_length           = 16
  network                 = google_compute_network.producer_net.id
}

resource "google_service_networking_connection" "service_con" {
  network                 = google_compute_network.producer_net.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

data "google_lustre_instance" "default" {
  instance_id             = google_lustre_instance.instance.instance_id
  zone                    = "us-central1-a"
}
`, context)
}
