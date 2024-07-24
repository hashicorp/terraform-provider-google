// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkeonprem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGkeonpremVmwareNodePool_vmwareNodePoolUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGkeonpremVmwareNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGkeonpremVmwareNodePool_vmwareNodePoolUpdateStart(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_node_pool.nodepool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccGkeonpremVmwareNodePool_vmwareNodePoolUpdate(context),
			},
			{
				ResourceName:            "google_gkeonprem_vmware_node_pool.nodepool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func testAccGkeonpremVmwareNodePool_vmwareNodePoolUpdateStart(context map[string]interface{}) string {
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

  resource "google_gkeonprem_vmware_node_pool" "nodepool" {
    name = "tf-test-nodepool-%{random_suffix}"
    location = "us-west1"
    vmware_cluster = google_gkeonprem_vmware_cluster.cluster.name
    annotations = {
      env = "test"
    }
    config {
        cpus = 4
        memory_mb = 8196
        replicas = 3
        image_type = "ubuntu_containerd"
        image = "image"
        boot_disk_size_gb = 10
        taints {
            key = "key"
            value = "value"
        }
        labels = {}
        vsphere_config {
          datastore = "test-datastore"
          tags {
            category = "test-category-1"
            tag = "tag-1"
          }
          tags {
            category = "test-category-2"
            tag = "tag-2"
          }
          host_groups = ["host1", "host2"]
        }
        enable_load_balancer = true
    }
    node_pool_autoscaling {
        min_replicas = 1
        max_replicas = 5
    }
  }
`, context)
}

func testAccGkeonpremVmwareNodePool_vmwareNodePoolUpdate(context map[string]interface{}) string {
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

  resource "google_gkeonprem_vmware_node_pool" "nodepool" {
    name = "tf-test-nodepool-%{random_suffix}"
    location = "us-west1"
    vmware_cluster = google_gkeonprem_vmware_cluster.cluster.name
    annotations = {
      env = "test-update"
    }
    config {
        cpus = 5
        memory_mb = 4096
        replicas = 3
        image_type = "windows"
        image = "image-updated"
        boot_disk_size_gb = 12
        taints {
            key = "key-updated"
            value = "value-updated"
        }
        labels = {}
        vsphere_config {
          datastore = "test-datastore-update"
          tags {
            category = "test-category-3"
            tag = "tag-3"
          }
          tags {
            category = "test-category-4"
            tag = "tag-4"
          }
          host_groups = ["host3", "host4"]
        }
        enable_load_balancer = false
    }
    node_pool_autoscaling {
        min_replicas = 2
        max_replicas = 6
    }
  }
`, context)
}
