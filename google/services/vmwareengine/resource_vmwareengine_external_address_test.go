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
	// Temporarily skipping so that this test does not run and consume resources during PR pushes. It is bound to fail and is being fixed by PR #10992
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":          "southamerica-east1", // using region with low node utilization.
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

resource "google_vmwareengine_network" "external-address-nw" {
  project = google_project.project.project_id
  name        = "tf-test-sample-external-address-nw%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "PC network description."

  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
}

resource "google_vmwareengine_private_cloud" "external-address-pc" {
  project = google_project.project.project_id
  location    = "%{region}-a"
  name        = "tf-test-sample-external-address-pc%{random_suffix}"
  type        = "TIME_LIMITED"
  description = "Sample test PC."
  network_config {
    management_cidr       = "192.168.1.0/24"
    vmware_engine_network = google_vmwareengine_network.external-address-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-sample-external-address-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 1
    }
  }
}

resource "google_vmwareengine_network_policy" "external-address-np" {
  project = google_project.project.project_id
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
