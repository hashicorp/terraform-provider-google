// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vpcaccess_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVPCAccessConnector_vpcAccessConnectorThroughput(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput(context),
			},
			{
				ResourceName:      "google_vpc_access_connector.connector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}
