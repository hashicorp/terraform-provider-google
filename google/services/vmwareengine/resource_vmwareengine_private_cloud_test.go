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

func TestAccVmwareenginePrivateCloud_vmwareEnginePrivateCloudUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":        "southamerica-west1",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareenginePrivateCloudDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testPrivateCloudUpdateConfig(context, "description1", 1),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_private_cloud.ds", "google_vmwareengine_private_cloud.vmw-engine-pc", map[string]struct{}{"type": {}}),
					testAccCheckGoogleVmwareengineNsxCredentialsMeta("data.google_vmwareengine_nsx_credentials.nsx-ds"),
					testAccCheckGoogleVmwareengineVcenterCredentialsMeta("data.google_vmwareengine_vcenter_credentials.vcenter-ds"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "type"},
			},
			{
				Config: testPrivateCloudUpdateConfig(context, "description2", 4), // Expand PC
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "type"},
			},
			{
				Config: testPrivateCloudUpdateConfig(context, "description2", 3), // Shrink PC
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "type"},
			},
		},
	})
}

func testPrivateCloudUpdateConfig(context map[string]interface{}, description string, nodeCount int) string {
	context["node_count"] = nodeCount
	context["description"] = description

	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  location = "%{region}-a"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "%{description}"
  type = "TIME_LIMITED"
  network_config {
    management_cidr = "192.168.30.0/24"
    vmware_engine_network = google_vmwareengine_network.default-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = "%{node_count}"
      custom_core_count = 32
    }
  }
}

data "google_vmwareengine_private_cloud" "ds" {
	location = "%{region}-a"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
}

# NSX and Vcenter Credentials are child datasources of PC and are included in the PC test due to the high deployment time involved in the Creation and deletion of a PC
data "google_vmwareengine_nsx_credentials" "nsx-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}

data "google_vmwareengine_vcenter_credentials" "vcenter-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}

`, context)
}

func testAccCheckGoogleVmwareengineNsxCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nsx credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckGoogleVmwareengineVcenterCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vcenter credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckVmwareenginePrivateCloudDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_private_cloud" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VmwareengineBasePath}}projects/{{project}}/locations/{{location}}/privateClouds/{{name}}")
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
				return fmt.Errorf("VmwareenginePrivateCloud still exists at %s", url)
			}
		}
		return nil
	}
}
