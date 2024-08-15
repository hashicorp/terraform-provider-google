// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareengineExternalAddress_vmwareEngineExternalAddressUpdate(t *testing.T) {
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
		CheckDestroy: testAccCheckVmwareengineExternalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareengineExternalAddressCreateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_external_address.ds-primary", "google_vmwareengine_external_address.vmw-engine-external-address-primary", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address-primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareengineExternalAddressUpdateConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address-primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareengineExternalAccessRuleCreateConfig(context),
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
				Config: testVmwareengineExternalAccessRuleUpdateConfig(context),
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

func testVmwareengineExternalAddressCreateConfig(context map[string]interface{}) string {
	return testVmwareengineBaseConfig(context) + testVmwareengineExternalAddressConfig(context, "primary", "sample description", "192.168.0.73")
}

func testVmwareengineExternalAddressUpdateConfig(context map[string]interface{}) string {
	return testVmwareengineBaseConfig(context) + testVmwareengineExternalAddressConfig(context, "primary", "sample updated description", "192.168.0.74")
}

func testVmwareengineExternalAccessRuleCreateConfig(context map[string]interface{}) string {
	return testVmwareengineBaseConfig(context) +
		testVmwareengineExternalAddressConfig(context, "primary", "sample updated description", "192.168.0.74") +
		testVmwareengineExternalAddressConfig(context, "secondary", "sample description", "192.168.0.75") +
		testVmwareengineExternalAccessRuleAllowConfig(context)
}

func testVmwareengineExternalAccessRuleUpdateConfig(context map[string]interface{}) string {
	return testVmwareengineBaseConfig(context) +
		testVmwareengineExternalAddressConfig(context, "primary", "sample updated description", "192.168.0.74") +
		testVmwareengineExternalAddressConfig(context, "secondary", "sample description", "192.168.0.75") +
		testVmwareengineExternalAccessRuleDenyConfig(context)
}

func testVmwareengineBaseConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-ea-ear-nw" {
  name        = "tf-test-sample-ea-ear-nw%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "PC network description."
}
resource "google_vmwareengine_private_cloud" "vmw-engine-ea-ear-pc" {
  location    = "%{region}-b"
  name        = "tf-test-sample-ea-ear-pc%{random_suffix}"
  type        = "TIME_LIMITED"
  description = "Sample test PC."
  deletion_delay_hours = 0
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr       = "192.168.1.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-ea-ear-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-external-address-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_network_policy" "vmw-engine-ea-ear-np" {
  location = "%{region}"
  name = "tf-test-sample-ea-ear-np%{random_suffix}"
  edge_services_cidr = "192.168.0.0/26"
  vmware_engine_network = google_vmwareengine_network.vmw-engine-ea-ear-nw.id
  internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }
}
	`, context)

}

func testVmwareengineExternalAddressConfig(context map[string]interface{}, id, description, internalIp string) string {
	context["id"] = id
	context["description"] = description
	context["internal_ip"] = internalIp
	return acctest.Nprintf(`
resource "google_vmwareengine_external_address" "vmw-engine-external-address-%{id}" {
  name = "tf-test-sample-external-address-%{id}%{random_suffix}"
  parent = google_vmwareengine_private_cloud.vmw-engine-ea-ear-pc.id
  internal_ip = "%{internal_ip}"
  description = "%{description}"
  depends_on = [
    google_vmwareengine_network_policy.vmw-engine-ea-ear-np,
  ]
}

data "google_vmwareengine_external_address" "ds-%{id}" {
  name = google_vmwareengine_external_address.vmw-engine-external-address-%{id}.name
  parent = google_vmwareengine_private_cloud.vmw-engine-ea-ear-pc.id
  depends_on = [
    google_vmwareengine_external_address.vmw-engine-external-address-%{id},
  ]
}
`, context)
}

func testVmwareengineExternalAccessRuleAllowConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
  name = "tf-test-sample-external-access-rule%{random_suffix}"
  parent =  google_vmwareengine_network_policy.vmw-engine-ea-ear-np.id
  description = "description1"
  priority = "101"
  action = "ALLOW"
  ip_protocol = "TCP"
  source_ip_ranges {
    ip_address_range = "0.0.0.0/0"
  }
  source_ports = ["80"]
  destination_ip_ranges {
    external_address = google_vmwareengine_external_address.vmw-engine-external-address-primary.id
  }
  destination_ip_ranges {
    external_address = google_vmwareengine_external_address.vmw-engine-external-address-secondary.id
  }
  destination_ports = ["433"]
  depends_on = [
    google_vmwareengine_external_address.vmw-engine-external-address-primary,
    google_vmwareengine_external_address.vmw-engine-external-address-secondary,
  ]
}

data "google_vmwareengine_external_access_rule" "ds" {
  name = google_vmwareengine_external_access_rule.vmw-engine-external-access-rule.name
  parent = google_vmwareengine_network_policy.vmw-engine-ea-ear-np.id
  depends_on = [
    google_vmwareengine_external_access_rule.vmw-engine-external-access-rule,
  ]
}
`, context)
}

func testVmwareengineExternalAccessRuleDenyConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
  name = "tf-test-sample-external-access-rule%{random_suffix}"
  parent =  google_vmwareengine_network_policy.vmw-engine-ea-ear-np.id
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
    google_vmwareengine_external_address.vmw-engine-external-address-primary,
    google_vmwareengine_external_address.vmw-engine-external-address-secondary,
  ]
}

data "google_vmwareengine_external_access_rule" "ds" {
  name = google_vmwareengine_external_access_rule.vmw-engine-external-access-rule.name
  parent = google_vmwareengine_network_policy.vmw-engine-ea-ear-np.id
  depends_on = [
    google_vmwareengine_external_access_rule.vmw-engine-external-access-rule,
  ]
}
`, context)
}

func testAccCheckVmwareengineExternalAddressDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_external_address" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VmwareengineBasePath}}{{parent}}/externalAddresses/{{name}}")
			if err != nil {
				return err
			}
			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VmwareengineExternalAddress still exists at %s", url)
			}
		}
		return nil
	}
}
