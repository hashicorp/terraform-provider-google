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

func TestAccVmwareenginePrivateCloud_vmwareEnginePrivateCloudUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":          "me-west1",
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
		CheckDestroy: testAccCheckVmwareenginePrivateCloudDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testPrivateCloudCreateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						map[string]struct{}{
							"type":                              {},
							"deletion_delay_hours":              {},
							"send_deletion_delay_hours_if_zero": {},
						}),
					testAccCheckGoogleVmwareengineNsxCredentialsMeta("data.google_vmwareengine_nsx_credentials.nsx-ds"),
					testAccCheckGoogleVmwareengineVcenterCredentialsMeta("data.google_vmwareengine_vcenter_credentials.vcenter-ds"),
				),
			},

			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "type", "deletion_parameters", "deletion_delay_hours"},
			},
			{
				Config: testPrivateCloudUpdateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						map[string]struct{}{
							"type":                              {},
							"deletion_delay_hours":              {},
							"send_deletion_delay_hours_if_zero": {},
						}),
				),
			},

			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "type", "deletion_parameters", "deletion_delay_hours"},
			},
		},
	})
}

func testPrivateCloudCreateConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "test description"
  type = "TIME_LIMITED"
  deletion_delay_hours = 1
  network_config {
    management_cidr = "192.168.30.0/24"
    vmware_engine_network = google_vmwareengine_network.default-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = 1
      custom_core_count = 32
    }
  }
}

data "google_vmwareengine_private_cloud" "ds" {
	location = "%{region}-b"
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

func testPrivateCloudUpdateConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "default-nw" {
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description       = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "updated description"
  type = "STANDARD"
  deletion_delay_hours = 0
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr = "192.168.30.0/24"
    vmware_engine_network = google_vmwareengine_network.default-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = 3
      custom_core_count = 32
    }
  }
}

data "google_vmwareengine_private_cloud" "ds" {
	location = "%{region}-b"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
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
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				pcState, ok := res["state"]
				if !ok {
					return fmt.Errorf("Unable to fetch state for existing VmwareenginePrivateCloud %s", url)
				}
				if pcState.(string) != "DELETED" {
					return fmt.Errorf("VmwareenginePrivateCloud still exists at %s", url)
				}
			}
		}
		return nil
	}
}
