// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareengineNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPeering_ds(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_network_peering.ds", "google_vmwareengine_network_peering.vmw-engine-network-peering", map[string]struct{}{}),
				),
			},
		},
	})
}

func testAccVmwareengineNetworkPeering_ds(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-peering-nw" {
	name              = "tf-test-sample-nw%{random_suffix}"
	location          = "global"
	type              = "STANDARD"
}

resource "google_vmwareengine_network" "network-peering-peer-nw" {
	name              = "tf-test-peer-nw%{random_suffix}"
	location          = "global"
	type              = "STANDARD"
}

resource "google_vmwareengine_network_peering" "vmw-engine-network-peering" {
	name = "tf-test-sample-network-peering%{random_suffix}"
	description = "Sample description"
	vmware_engine_network = google_vmwareengine_network.network-peering-nw.id
	peer_network = google_vmwareengine_network.network-peering-peer-nw.id
	peer_network_type = "VMWARE_ENGINE_NETWORK"
}

data "google_vmwareengine_network_peering" "ds" {
	name = google_vmwareengine_network_peering.vmw-engine-network-peering.name
	depends_on = [
		google_vmwareengine_network_peering.vmw-engine-network-peering,
	]
}
`, context)
}
