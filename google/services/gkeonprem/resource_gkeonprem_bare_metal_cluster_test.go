// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkeonprem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLbStart(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_cluster.cluster-metallb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLb(context),
			},
			{
				ResourceName:            "google_gkeonprem_bare_metal_cluster.cluster-metallb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-manuallb",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-manuallb",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremBareMetalClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-bgplb",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_bare_metal_cluster.cluster-bgplb",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-metallb" {
    name = "cluster-metallb%{random_suffix}"
    location = "us-west1"
    annotations = {
      env = "test"
    }
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/16"]
        pod_address_cidr_blocks = ["10.240.0.0/13"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.9"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 443
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      metal_lb_config {
        address_pools {
          pool = "pool1"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "10.200.0.18/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128",
            "fd00:1::12/128"
          ]
        }
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
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateMetalLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-metallb" {
    name = "cluster-metallb%{random_suffix}"
    location = "us-west1"
    annotations = {
      env = "test-update"
    }
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      metal_lb_config {
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128"
          ]
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-manuallb" {
    name = "cluster-manuallb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      metal_lb_config {
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.14/32",
            "10.200.0.15/32",
            "10.200.0.16/32",
            "10.200.0.17/32",
            "fd00:1::f/128",
            "fd00:1::10/128",
            "fd00:1::11/128"
          ]
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share"
          storage_class = "local-shared"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk"
        storage_class = "local-disks"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
    binary_authorization {
      evaluation_mode = "DISABLED"
    }
    upgrade_policy {
      policy = "SERIAL"
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateManualLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-manuallb" {
    name = "cluster-manuallb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      manual_lb_config {
        enabled = true
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
    binary_authorization {
      evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
    }
    upgrade_policy {
      policy = "CONCURRENT"
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-bgplb" {
    name = "cluster-bgplb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.13"
        ingress_vip = "10.200.0.14"
      }
      bgp_lb_config {
        asn = 123456
        bgp_peer_configs {
          asn = 123457
          ip_address = "10.0.0.1"
          control_plane_nodes = ["test-node"]
        }
        address_pools {
          pool = "pool1"
          addresses = [
            "10.200.0.14/32",
            "fd00:1::12/128"
          ]
        }
        load_balancer_node_pool_config {
          node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.9"
            }
          }
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share"
          storage_class = "local-shared"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk"
        storage_class = "local-disks"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin@hashicorptest.com"
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremBareMetalCluster_bareMetalClusterUpdateBgpLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_bare_metal_cluster" "cluster-bgplb" {
    name = "cluster-bgplb%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    bare_metal_version = "1.12.3"
    network_config {
      island_mode_cidr {
        service_address_cidr_blocks = ["172.26.0.0/20"]
        pod_address_cidr_blocks = ["10.240.0.0/14"]
      }
    }
    control_plane {
      control_plane_node_pool_config {
        node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.10"
            }
        }
      }
    }
    load_balancer {
      port_config {
        control_plane_load_balancer_port = 80
      }
      vip_config {
        control_plane_vip = "10.200.0.14"
        ingress_vip = "10.200.0.15"
      }
      bgp_lb_config {
        asn = 123457
        bgp_peer_configs {
          asn = 123458
          ip_address = "10.0.0.2"
          control_plane_nodes = ["test-node-updated"]
        }
        address_pools {
          pool = "pool2"
          addresses = [
            "10.200.0.15/32",
            "fd00:1::16/128"
          ]
        }
        load_balancer_node_pool_config {
          node_pool_config {
            labels = {}
            operating_system = "LINUX"
            node_configs {
              labels = {}
              node_ip = "10.200.0.11"
            }
          }
        }
      }
    }
    storage {
      lvp_share_config {
        lvp_config {
          path = "/mnt/localpv-share-updated"
          storage_class = "local-shared-updated"
        }
        shared_path_pv_count = 6
      }
      lvp_node_mounts_config {
        path = "/mnt/localpv-disk-updated"
        storage_class = "local-disks-updated"
      }
    }
    security_config {
      authorization {
        admin_users {
          username = "admin-updated@hashicorptest.com"
        }
      }
    }
  }
`, context)
}
