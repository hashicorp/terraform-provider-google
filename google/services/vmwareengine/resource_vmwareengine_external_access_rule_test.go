// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineExternalAccessRule_vmwareEngineExternalAccessRuleUpdate(t *testing.T) {
	// Temporarily skipping so that this test does not run and consume resources during PR pushes. It is bound to fail and is being fixed by PR #10992
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":          "southamerica-west1", // using region with low node utilization.
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
		Steps: []resource.TestStep{
			{
				Config: testVmwareEngineExternalAccessRuleCreateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_external_access_rule.ds", "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineExternalAccessRuleUpdateConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineExternalAccessRuleCreateConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "external-access-rule-nw" {
  name        = "tf-test-sample-external-access-rule-nw-%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
}

resource "google_vmwareengine_private_cloud" "external-access-rule-pc" {
  location    = "%{region}-a"
  name        = "tf-test-sample-external-access-rule-pc-%{random_suffix}"
  type        = "TIME_LIMITED"
  network_config {
    management_cidr       = "192.168.1.0/24"
    vmware_engine_network = google_vmwareengine_network.external-access-rule-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-sample-external-access-rule-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_network_policy" "external-access-rule-np" {
  location = "%{region}"
  name = "tf-test-sample-external-access-rule-np-%{random_suffix}"
  edge_services_cidr = "192.168.0.0/26"
  vmware_engine_network = google_vmwareengine_network.external-access-rule-nw.id
	internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }
}

resource "google_vmwareengine_external_address" "external-access-rule-ea-one" {
  name = "tf-test-sample-external-access-rule-ea-one%{random_suffix}"
  parent =  google_vmwareengine_private_cloud.external-access-rule-pc.id
  internal_ip = "192.168.0.74"
  description = "Sample description."
  depends_on = [
    google_vmwareengine_network_policy.external-access-rule-np,
  ]
}

resource "google_vmwareengine_external_address" "external-access-rule-ea-two" {
  name = "tf-test-sample-external-access-rule-ea-two%{random_suffix}"
  parent =  google_vmwareengine_private_cloud.external-access-rule-pc.id
  internal_ip = "192.168.0.75"
  description = "Sample description."
  depends_on = [
    google_vmwareengine_network_policy.external-access-rule-np,
  ]
}

resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
  name = "tf-test-sample-external-access-rule%{random_suffix}"
  parent =  google_vmwareengine_network_policy.external-access-rule-np.id
  description = "description1"
  priority = "101"
  action = "ALLOW"
  ip_protocol = "TCP"
  source_ip_ranges {
    ip_address_range = "0.0.0.0/0"
  }
  source_ports = ["80"]
  destination_ip_ranges {
    external_address = google_vmwareengine_external_address.external-access-rule-ea-one.id
  }
  destination_ip_ranges {
    external_address = google_vmwareengine_external_address.external-access-rule-ea-two.id
  }
  destination_ports = ["433"]
  depends_on = [
    google_vmwareengine_external_address.external-access-rule-ea-one,
    google_vmwareengine_external_address.external-access-rule-ea-two,
  ]
}

data "google_vmwareengine_external_access_rule" "ds" {
  name = google_vmwareengine_external_access_rule.vmw-engine-external-access-rule.name
  parent = google_vmwareengine_network_policy.external-access-rule-np.id
  depends_on = [
    google_vmwareengine_external_access_rule.vmw-engine-external-access-rule,
  ]
}
`, context)
}

func testVmwareEngineExternalAccessRuleUpdateConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "external-access-rule-nw" {
  name        = "tf-test-sample-external-access-rule-nw-%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
}

resource "google_vmwareengine_private_cloud" "external-access-rule-pc" {
  location    = "%{region}-a"
  name        = "tf-test-sample-external-access-rule-pc-%{random_suffix}"
  type        = "TIME_LIMITED"
  network_config {
    management_cidr       = "192.168.1.0/24"
    vmware_engine_network = google_vmwareengine_network.external-access-rule-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-sample-external-access-rule-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_network_policy" "external-access-rule-np" {
  location = "%{region}"
  name = "tf-test-sample-external-access-rule-np-%{random_suffix}"
  edge_services_cidr = "192.168.0.0/26"
  vmware_engine_network = google_vmwareengine_network.external-access-rule-nw.id
	internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }
}

resource "google_vmwareengine_external_address" "external-access-rule-ea-one" {
  name = "tf-test-sample-external-access-rule-ea-one%{random_suffix}"
  parent =  google_vmwareengine_private_cloud.external-access-rule-pc.id
  internal_ip = "192.168.0.74"
  description = "Sample description."
  depends_on = [
    google_vmwareengine_network_policy.external-access-rule-np,
  ]
}

resource "google_vmwareengine_external_address" "external-access-rule-ea-two" {
  name = "tf-test-sample-external-access-rule-ea-two%{random_suffix}"
  parent =  google_vmwareengine_private_cloud.external-access-rule-pc.id
  internal_ip = "192.168.0.75"
  description = "Sample description."
  depends_on = [
    google_vmwareengine_network_policy.external-access-rule-np,
  ]
}

resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
  name = "tf-test-sample-external-access-rule%{random_suffix}"
  parent =  google_vmwareengine_network_policy.external-access-rule-np.id
  description = "description2"
  priority = "102"
  action = "DENY"
  ip_protocol = "UDP"
  source_ip_ranges {
    ip_address = "192.168.40.1"
  }
  source_ip_ranges {
    ip_address = "192.168.40.0"
  }
  source_ports = ["81", "82"]
  destination_ip_ranges {
    ip_address_range = "0.0.0.0/0"
  }
  destination_ports = ["435", "434"]
  depends_on = [
    google_vmwareengine_external_address.external-access-rule-ea-one,
    google_vmwareengine_external_address.external-access-rule-ea-two,
  ]
}

data "google_vmwareengine_external_access_rule" "ds" {
  name = google_vmwareengine_external_access_rule.vmw-engine-external-access-rule.name
  parent = google_vmwareengine_network_policy.external-access-rule-np.id
  depends_on = [
    google_vmwareengine_external_access_rule.vmw-engine-external-access-rule,
  ]
}
`, context)
}
