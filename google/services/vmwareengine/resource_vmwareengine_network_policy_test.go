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
	t.Parallel()

	context := map[string]interface{}{
		"region":          envvar.GetTestRegionFromEnv(),
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
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "vmwareengine" {
  project = google_project.project.project_id
  service = "vmwareengine.googleapis.com"
}

resource "time_sleep" "sleep" {
  create_duration = "1m"
  depends_on = [
    google_project_service.vmwareengine,
  ]
}

resource "google_vmwareengine_network" "network-policy-nw" {
  project           = google_project.project.project_id
  name              = "tf-test-sample-nw%{random_suffix}"
  location          = "global" 
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"

  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
  project           = google_project.project.project_id
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
