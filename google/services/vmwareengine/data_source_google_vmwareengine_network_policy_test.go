// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceVmwareengineNetworkPolicy_basic(t *testing.T) {
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

resource "google_vmwareengine_network" "network-policy-ds-nw" {
  project = google_project.project.project_id
  name = "tf-test-sample-nw%{random_suffix}"
  location = "global" 
  type = "STANDARD"
  description = "VMwareEngine standard network sample"

  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
  project = google_project.project.project_id
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

  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
}

data "google_vmwareengine_network_policy" "ds" {
  project = google_project.project.project_id
  name = google_vmwareengine_network_policy.vmw-engine-network-policy.name
  location = "%{region}"
  depends_on = [
    google_vmwareengine_network_policy.vmw-engine-network-policy,
  ]
}

`, context)
}
