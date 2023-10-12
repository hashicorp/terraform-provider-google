// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkeonprem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLbStart(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLb(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateF5Lb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5LbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5lb(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLbStart(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(context),
			},
			{
				ResourceName:      "google_gkeonprem_vmware_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {
      env = "test"
    }
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      metal_lb_config {
        address_pools {
          pool = "ingress-ip"
          manual_assign = "true"
          addresses = ["10.251.135.19"]
          avoid_buggy_ips = true
        }
        address_pools {
          pool = "lb-test-ip"
          manual_assign = "true"
          addresses = ["10.251.135.19"]
          avoid_buggy_ips = true
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateMetalLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster updated"
    on_prem_version = "1.13.1-gke.36"
    annotations = {
      env = "test-update"
    }
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/16"]
      pod_address_cidr_blocks = ["192.168.0.0/20"]
      dhcp_ip_config {
        enabled = true
      }
    }
    control_plane_node {
       cpus = 5
       memory = 4098
       replicas = 3
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.6"
        ingress_vip = "10.251.135.20"
      }
      metal_lb_config {
        address_pools {
          pool = "ingress-ip-updated"
          manual_assign = "false"
          addresses = ["10.251.135.20"]
          avoid_buggy_ips = false
        }
        address_pools {
          pool = "lb-test-ip-updated"
          manual_assign = "false"
          addresses = ["10.251.135.20"]
          avoid_buggy_ips = false
        }
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5LbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
      control_plane_v2_config {
        control_plane_ip_block {
          ips {
            hostname = "test-hostname"
            ip = "10.0.0.1"
          }
          netmask="10.0.0.1/32"
          gateway="test-gateway"
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      f5_config {
        address = "10.0.0.1"
        partition = "test-partition"
        snat_pool = "test-snap-pool"
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateF5lb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      dhcp_ip_config {
        enabled = true
      }
      control_plane_v2_config {
        control_plane_ip_block {
          ips {
            hostname = "test-hostname-updated"
            ip = "10.0.0.2"
          }
          netmask="10.0.0.2/32"
          gateway="test-gateway-updated"
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      f5_config {
        address = "10.0.0.2"
        partition = "test-partition-updated"
        snat_pool = "test-snap-pool-updated"
      }
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLbStart(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      host_config {
        dns_servers = ["10.254.41.1"]
        ntp_servers = ["216.239.35.8"]
        dns_search_domains = ["test-domain"]
      }
      static_ip_config {
        ip_blocks {
          netmask = "255.255.252.0"
          gateway = "10.251.31.254"
          ips {
            ip = "10.251.30.153"
            hostname = "test-hostname1"
          }
          ips {
            ip = "10.251.31.206"
            hostname = "test-hostname2"
          }
          ips {
            ip = "10.251.31.193"
            hostname = "test-hostname3"
          }
          ips { 
            ip = "10.251.30.230"
            hostname = "test-hostname4"
          }
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      manual_lb_config {
        ingress_http_node_port = 30005
        ingress_https_node_port = 30006
        control_plane_node_port = 30007
        konnectivity_server_node_port = 30008
      }
    }
    vcenter {
      resource_pool = "test-resource-pool"
      datastore = "test-datastore"
      datacenter = "test-datacenter"
      cluster = "test-cluster"
      folder = "test-folder"
      ca_cert_data = "test-ca-cert-data"
      storage_policy_name = "test-storage-policy-name"
    }
    dataplane_v2 {
      dataplane_v2_enabled = true
      windows_dataplane_v2_enabled = true
      advanced_networking = true
    }
    vm_tracking_enabled = true
    enable_control_plane_v2 = true
    upgrade_policy {
      control_plane_only = true
    }
    authorization {
      admin_users {
        username = "testuser@gmail.com"
      }
    }
    anti_affinity_groups {
      aag_config_disabled = true
    }
    auto_repair_config {
      enabled = true
    }
  }
`, context)
}

func testAccGkeonpremVmwareCluster_vmwareClusterUpdateManualLb(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_gkeonprem_vmware_cluster" "cluster" {
    name = "tf-test-cluster-%{random_suffix}"
    location = "us-west1"
    admin_cluster_membership = "projects/870316890899/locations/global/memberships/gkeonprem-terraform-test"
    description = "test cluster"
    on_prem_version = "1.13.1-gke.35"
    annotations = {}
    network_config {
      service_address_cidr_blocks = ["10.96.0.0/12"]
      pod_address_cidr_blocks = ["192.168.0.0/16"]
      host_config {
        dns_servers = ["10.254.41.1"]
        ntp_servers = ["216.239.35.8"]
        dns_search_domains = ["test-domain"]
      }
      static_ip_config {
        ip_blocks {
          netmask = "255.255.252.1"
          gateway = "10.251.31.255"
          ips {
            ip = "10.251.30.154"
            hostname = "test-hostname1-updated"
          }
          ips {
            ip = "10.251.31.206"
            hostname = "test-hostname2"
          }
          ips {
            ip = "10.251.31.193"
            hostname = "test-hostname3"
          }
          ips { 
            ip = "10.251.30.230"
            hostname = "test-hostname4"
          }
        }
      }
    }
    control_plane_node {
       cpus = 4
       memory = 8192
       replicas = 1
    }
    load_balancer {
      vip_config {
        control_plane_vip = "10.251.133.5"
        ingress_vip = "10.251.135.19"
      }
      manual_lb_config {
        ingress_http_node_port = 30006
        ingress_https_node_port = 30007
        control_plane_node_port = 30008
        konnectivity_server_node_port = 30009
      }
    }
    vcenter {
      resource_pool = "test-resource-pool-updated"
      datastore = "test-datastore-updated"
      datacenter = "test-datacenter-updated"
      cluster = "test-cluster-updated"
      folder = "test-folder-updated"
      ca_cert_data = "test-ca-cert-data-updated"
      storage_policy_name = "test-storage-policy-name-updated"
    }
    dataplane_v2 {
      dataplane_v2_enabled = true
      windows_dataplane_v2_enabled = true
      advanced_networking = true
    }
    vm_tracking_enabled = false
    enable_control_plane_v2 = false
    upgrade_policy {
      control_plane_only = true
    }
    authorization {
      admin_users {
        username = "testuser-updated@gmail.com"
      }
    }
    anti_affinity_groups {
      aag_config_disabled = true
    }
    auto_repair_config {
      enabled = true
    }
  }
`, context)
}
