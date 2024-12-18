// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineNetworkPeering_update(t *testing.T) {
	t.Skip("https://github.com/hashicorp/terraform-provider-google/issues/20719")
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPeering_config(context, "Sample description."),
			},
			{
				ResourceName:            "google_vmwareengine_network_peering.vmw-engine-network-peering",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccVmwareengineNetworkPeering_config(context, "Updated description."),
			},
			{
				ResourceName:            "google_vmwareengine_network_peering.vmw-engine-network-peering",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccVmwareengineNetworkPeering_config(context map[string]interface{}, description string) string {
	context["description"] = description
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
  description = "%{description}"
  vmware_engine_network = google_vmwareengine_network.network-peering-nw.id
  peer_network = google_vmwareengine_network.network-peering-peer-nw.id
  peer_network_type = "VMWARE_ENGINE_NETWORK"
}
`, context)
}
