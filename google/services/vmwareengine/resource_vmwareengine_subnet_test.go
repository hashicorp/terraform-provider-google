// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVmwareengineSubnet_vmwareEngineUserDefinedSubnetUpdate(t *testing.T) {
	// Temporarily skipping so that this test does not run and consume resources during PR pushes. It is bound to fail and is being fixed by PR #10992
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":        "southamerica-west1", // using region with low node utilization.
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareEngineSubnetConfig(context, "192.168.1.0/26"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_subnet.ds", "google_vmwareengine_subnet.vmw-engine-subnet", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineSubnetConfig(context, "192.168.2.0/26"),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineSubnetConfig(context map[string]interface{}, ipCidrRange string) string {
	context["ip_cidr_range"] = ipCidrRange
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "subnet-nw" {
  name        = "tf-test-subnet-nw%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "PC network description."
}

resource "google_vmwareengine_private_cloud" "subnet-pc" {
  location    = "%{region}-a"
  name        = "tf-test-subnet-pc%{random_suffix}"
  type        = "TIME_LIMITED"
  description = "Sample test PC."
  network_config {
    management_cidr       = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.subnet-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-mgmt-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_subnet" "vmw-engine-subnet" {
  name = "service-2"
  parent =  google_vmwareengine_private_cloud.subnet-pc.id
  ip_cidr_range = "%{ip_cidr_range}"
}

data "google_vmwareengine_subnet" ds {
  name = "service-2"
  parent = google_vmwareengine_private_cloud.subnet-pc.id
  depends_on = [
    google_vmwareengine_subnet.vmw-engine-subnet,
  ]
}
`, context)
}
