// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineNetworkPolicy_update(t *testing.T) {
	t.Skip("https://github.com/hashicorp/terraform-provider-google/issues/20719")
	t.Parallel()

	context := map[string]interface{}{
		"region":          "me-west1", // region with allocated quota
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
		CheckDestroy: testAccCheckVmwareengineNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description1", "192.168.0.0/26", false, false),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, "description2", "192.168.1.0/26", true, true),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time"},
			},
		},
	})
}

func testAccVmwareengineNetworkPolicy_config(context map[string]interface{}, description string, edgeServicesCidr string, internetAccess bool, externalIp bool) string {
	context["internet_access"] = internetAccess
	context["external_ip"] = externalIp
	context["edge_services_cidr"] = edgeServicesCidr
	context["description"] = description

	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-policy-nw" {
  name              = "tf-test-sample-nw%{random_suffix}"
  location          = "global" 
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
  location = "%{region}"
  name = "tf-test-sample-network-policy%{random_suffix}"
  description = "%{description}" 
  internet_access {
    enabled = "%{internet_access}"
  }
  external_ip {
    enabled = "%{external_ip}"
  }
  edge_services_cidr = "%{edge_services_cidr}"
  vmware_engine_network = google_vmwareengine_network.network-policy-nw.id
}
`, context)
}
