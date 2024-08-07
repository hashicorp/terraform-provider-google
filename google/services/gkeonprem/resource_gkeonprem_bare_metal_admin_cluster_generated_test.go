// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package gkeonprem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterBasicExample(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_admin_cluster.admin-cluster-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "name"},
			},
		},
	})
}

func testAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gkeonprem_bare_metal_admin_cluster" "admin-cluster-basic" {
  name = "tf-test-my-cluster%{random_suffix}"
  location = "us-west1"
  bare_metal_version = "1.13.4"
  network_config {
    island_mode_cidr {
      service_address_cidr_blocks = ["172.26.0.0/16"]
      pod_address_cidr_blocks = ["10.240.0.0/13"]
    }
  }
  node_config {
      max_pods_per_node = 250
  }
  control_plane {
    control_plane_node_pool_config {
      node_pool_config {
        labels = {}
        operating_system = "LINUX"
        node_configs {
            labels  = {}
            node_ip = "10.200.0.2"
        }
        node_configs {
            labels  = {}
            node_ip = "10.200.0.3"
        }
        node_configs {
            labels  = {}
            node_ip = "10.200.0.4"
        }
      }
    }
  }
  load_balancer {
    port_config {
      control_plane_load_balancer_port = 443
    }
    vip_config {
      control_plane_vip = "10.200.0.5"
    }
  }
  storage {
    lvp_share_config {
      lvp_config {
        path = "/mnt/localpv-share"
        storage_class = "local-shared"
      }
      shared_path_pv_count = 5
    }
    lvp_node_mounts_config {
      path = "/mnt/localpv-disk"
      storage_class = "local-disks"
    }
  }
  node_access_config {
    login_user = "root"
  }
}
`, context)
}

func TestAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterFullExample(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_admin_cluster.admin-cluster-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "name"},
			},
		},
	})
}

func testAccGkeonpremBareMetalAdminCluster_gkeonpremBareMetalAdminClusterFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gkeonprem_bare_metal_admin_cluster" "admin-cluster-basic" {
  name = "tf-test-my-cluster%{random_suffix}"
  location = "us-west1"
  description = "test description"
  bare_metal_version = "1.13.4"
  annotations = {
    env = "test"
  }
  network_config {
    island_mode_cidr {
      service_address_cidr_blocks = ["172.26.0.0/16"]
      pod_address_cidr_blocks = ["10.240.0.0/13"]
    }
  }
  node_config {
    max_pods_per_node = 250
  }
  control_plane {
    control_plane_node_pool_config {
      node_pool_config {
        labels = {}
        operating_system = "LINUX"
        node_configs {
          labels  = {}
          node_ip = "10.200.0.2"
        }
        node_configs {
          labels  = {}
          node_ip = "10.200.0.3"
        }
        node_configs {
          labels  = {}
          node_ip = "10.200.0.4"
        }
        taints {
          key = "test-key"
          value = "test-value"
          effect = "NO_EXECUTE"
        }
      }
    }
    api_server_args {
      argument = "test argument"
      value = "test value"
    }
  }
  load_balancer {
    port_config {
      control_plane_load_balancer_port = 443
    }
    vip_config {
      control_plane_vip = "10.200.0.5"
    }
    manual_lb_config {
      enabled = true
    }
  }
  storage {
    lvp_share_config {
      lvp_config {
        path = "/mnt/localpv-share"
        storage_class = "local-shared"
      }
      shared_path_pv_count = 5
    }
    lvp_node_mounts_config {
      path = "/mnt/localpv-disk"
      storage_class = "local-disks"
    }
  }
  node_access_config {
    login_user = "root"
  }
  security_config {
    authorization {
      admin_users {
        username = "admin@hashicorptest.com"
      }
    }
  }
  maintenance_config {
    maintenance_address_cidr_blocks = ["10.0.0.1/32", "10.0.0.2/32"]
  }
  cluster_operations {
    enable_application_logs = true
  }
  proxy {
    uri = "test proxy uri"
    no_proxy = ["127.0.0.1"]
  }
}
`, context)
}
