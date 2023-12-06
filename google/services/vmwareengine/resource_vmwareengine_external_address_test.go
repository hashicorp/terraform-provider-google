// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareengineExternalAddress_vmwareEngineExternalAddressUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        "southamerica-east1", // using region with low node utilization.
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineExternalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareEngineExternalAddressConfig(context, "description1", "192.168.0.66"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_external_address.ds", "google_vmwareengine_external_address.vmw-engine-external-address", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineExternalAddressConfig(context, "description2", "192.168.0.67"),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineExternalAddressConfig(context map[string]interface{}, description string, internalIp string) string {
	context["internal_ip"] = internalIp
	context["description"] = description
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "external-address-nw" {
  name        = "tf-test-sample-external-address-nw%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "PC network description."
}

resource "google_vmwareengine_private_cloud" "external-address-pc" {
  location    = "%{region}-a"
  name        = "tf-test-sample-external-address-pc%{random_suffix}"
  description = "Sample test PC."
  network_config {
    management_cidr       = "192.168.1.0/24"
    vmware_engine_network = google_vmwareengine_network.external-address-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-sample-external-address-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 3
    }
  }
}

resource "google_vmwareengine_network_policy" "external-address-np" {
  location = "%{region}"
  name = "tf-test-sample-external-address-np%{random_suffix}"
  edge_services_cidr = "192.168.0.0/26"
  vmware_engine_network = google_vmwareengine_network.external-address-nw.id
	internet_access {
    enabled = true
  }
  external_ip {
    enabled = true
  }
}

resource "google_vmwareengine_external_address" "vmw-engine-external-address" {
  name = "tf-test-sample-external-address%{random_suffix}"
  parent = google_vmwareengine_private_cloud.external-address-pc.id
  internal_ip = "%{internal_ip}"
  description = "%{description}"
  depends_on = [
    google_vmwareengine_network_policy.external-address-np,
  ]
}

data "google_vmwareengine_external_address" "ds" {
  name = google_vmwareengine_external_address.vmw-engine-external-address.name
  parent = google_vmwareengine_private_cloud.external-address-pc.id
  depends_on = [
    google_vmwareengine_external_address.vmw-engine-external-address,
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
