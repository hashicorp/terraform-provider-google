// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceVmwareengineNetworkPolicy_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPolicy_ds(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_network_policy.ds", "google_vmwareengine_network_policy.vmw-engine-network-policy", map[string]struct{}{}),
				),
			},
		},
	})
}

func testAccVmwareengineNetworkPolicy_ds(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-policy-ds-nw" {
	name = "tf-test-sample-nw%{random_suffix}"
	location = "global" 
	type = "STANDARD"
	description = "VMwareEngine standard network sample"
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
	location = "%{region}"
	name = "tf-test-sample-network-policy%{random_suffix}"
	internet_access {
		enabled = true
	}
	external_ip {
		enabled = true
	}
	edge_services_cidr = "192.168.30.0/26"
	vmware_engine_network = google_vmwareengine_network.network-policy-ds-nw.id
}

data "google_vmwareengine_network_policy" "ds" {
  name = google_vmwareengine_network_policy.vmw-engine-network-policy.name
  location = "%{region}"
  depends_on = [
    google_vmwareengine_network_policy.vmw-engine-network-policy,
  ]
}

`, context)
}
