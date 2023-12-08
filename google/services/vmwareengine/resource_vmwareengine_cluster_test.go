// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareengineCluster_vmwareEngineClusterUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"region":          "southamerica-west1", // using region with low node utilization.
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareEngineClusterConfig(context, 3),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_cluster.ds", "google_vmwareengine_cluster.vmw-engine-ext-cluster", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_cluster.vmw-engine-ext-cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineClusterConfig(context, 4), // expand the cluster
			},
			{
				ResourceName:            "google_vmwareengine_cluster.vmw-engine-ext-cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineClusterConfig(context, 3), // shrink the cluster.
			},
			{
				ResourceName:            "google_vmwareengine_cluster.vmw-engine-ext-cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineClusterConfig(context map[string]interface{}, nodeCount int) string {
	context["node_count"] = nodeCount
	return acctest.Nprintf(`

resource "google_vmwareengine_network" "cluster-nw" {
  name        = "tf-test-cluster-nw%{random_suffix}"
  location    = "global"
  type        = "STANDARD"
  description = "PC network description."
}

resource "google_vmwareengine_private_cloud" "cluster-pc" {
  location    = "%{region}-a"
  name        = "tf-test-cluster-pc%{random_suffix}"
  description = "Sample test PC."
  network_config {
    management_cidr       = "192.168.10.0/24"
    vmware_engine_network = google_vmwareengine_network.cluster-nw.id
  }

  management_cluster {
    cluster_id = "tf-test-mgmt-cluster%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count   = 3
    }
  }
}

resource "google_vmwareengine_cluster" "vmw-engine-ext-cluster" {
  name = "tf-test-ext-cluster%{random_suffix}"
  parent =  google_vmwareengine_private_cloud.cluster-pc.id
  node_type_configs {
    node_type_id = "standard-72"
    node_count   = %{node_count}
		custom_core_count = 32
  }
}

data "google_vmwareengine_cluster" ds {
  name = "tf-test-ext-cluster%{random_suffix}"
  parent = google_vmwareengine_private_cloud.cluster-pc.id
  depends_on = [
	google_vmwareengine_cluster.vmw-engine-ext-cluster,
  ]
}
`, context)
}

func testAccCheckVmwareengineClusterDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_cluster" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VmwareengineBasePath}}{{parent}}/clusters/{{name}}")
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
				return fmt.Errorf("VmwareengineCluster still exists at %s", url)
			}
		}

		return nil
	}
}
