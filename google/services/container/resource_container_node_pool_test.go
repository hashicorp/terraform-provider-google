// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basic(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_resourceManagerTags(t *testing.T) {
	t.Parallel()
	pid := envvar.GetTestProjectFromEnv()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)

	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_resourceManagerTags(pid, clusterName, networkName, subnetworkName, randomSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_node_pool.primary_nodes", "node_config.0.resource_manager_tags.%"),
				),
			},
			{
				ResourceName:            "google_container_node_pool.primary_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "cluster"},
			},
			{
				Config: testAccContainerNodePool_resourceManagerTagsUpdate1(pid, clusterName, networkName, subnetworkName, randomSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_node_pool.primary_nodes", "node_config.0.resource_manager_tags.%"),
				),
			},
			{
				ResourceName:            "google_container_node_pool.primary_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "cluster"},
			},
			{
				Config: testAccContainerNodePool_resourceManagerTagsUpdate2(pid, clusterName, networkName, subnetworkName, randomSuffix),
			},
			{
				ResourceName:            "google_container_node_pool.primary_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "cluster"},
			},
		},
	})
}

func TestAccContainerNodePool_basicWithClusterId(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basicWithClusterId(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_node_pool.np",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster"},
			},
		},
	})
}

func TestAccContainerNodePool_nodeLocations(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_nodeLocations(cluster, np, network),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_maxPodsPerNode(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_maxPodsPerNode(cluster, np, network),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_namePrefix(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_namePrefix(cluster, "tf-np-", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_node_pool.np",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccContainerNodePool_noName(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_noName(cluster, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withLoggingVariantUpdates(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "DEFAULT", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_logging_variant",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "MAX_THROUGHPUT", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_logging_variant",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "DEFAULT", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_logging_variant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNodeConfig(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#", "node_config.0.taint"},
			},
			{
				Config: testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#", "node_config.0.taint"},
			},
		},
	})
}

func TestAccContainerNodePool_withTaintsUpdate(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basic(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withTaintsUpdate(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#", "node_config.0.taint"},
			},
		},
	})
}

func TestAccContainerNodePool_withMachineAndDiskUpdate(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basic(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withMachineAndDiskUpdate(cluster, nodePool, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#", "node_config.0.taint"},
			},
		},
	})
}

func TestAccContainerNodePool_withReservationAffinity(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withReservationAffinity(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.#", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "ANY_RESERVATION"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_reservation_affinity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withReservationAffinitySpecific(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	reservation := fmt.Sprintf("tf-test-reservation-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withReservationAffinitySpecific(cluster, reservation, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.#", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "SPECIFIC_RESERVATION"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.0.key", "compute.googleapis.com/reservation-name"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.0.values.#", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_reservation_affinity",
						"node_config.0.reservation_affinity.0.values.0", reservation),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_reservation_affinity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withWorkloadIdentityConfig(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withWorkloadMetadataConfig(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.mode", "GCE_METADATA"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_workload_metadata_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withWorkloadMetadataConfig_gkeMetadata(pid, cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.mode", "GKE_METADATA"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_workload_metadata_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withKubeletConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withKubeletConfig(cluster, np, "static", "100ms", networkName, subnetworkName, true, 2048),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_kubelet_config",
						"node_config.0.kubelet_config.0.cpu_cfs_quota", "true"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_kubelet_config",
						"node_config.0.kubelet_config.0.pod_pids_limit", "2048"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_kubelet_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withKubeletConfig(cluster, np, "", "", networkName, subnetworkName, false, 1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_kubelet_config",
						"node_config.0.kubelet_config.0.cpu_cfs_quota", "false"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_kubelet_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withInvalidKubeletCpuManagerPolicy(t *testing.T) {
	t.Parallel()
	// Unit test, no interactions
	acctest.SkipIfVcr(t)

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerNodePool_withKubeletConfig(cluster, np, "dontexist", "100us", networkName, subnetworkName, true, 1024),
				ExpectError: regexp.MustCompile(`.*to be one of \["?static"? "?none"? "?"?\].*`),
			},
		},
	})
}

func TestAccContainerNodePool_withLinuxNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			// Create a node pool with empty `linux_node_config.sysctls`.
			{
				Config: testAccContainerNodePool_withLinuxNodeConfig(cluster, np, "", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_linux_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withLinuxNodeConfig(cluster, np, "1000 20000 100000", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_linux_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Perform an update.
			{
				Config: testAccContainerNodePool_withLinuxNodeConfig(cluster, np, "1000 20000 200000", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_linux_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withCgroupMode(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withCgroupMode(cluster, np, "CGROUP_MODE_V2", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Perform an update.
			{
				Config: testAccContainerNodePool_withCgroupMode(cluster, np, "CGROUP_MODE_UNSPECIFIED", networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withNetworkConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNetworkConfig(cluster, np, network, "TIER_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.with_pco_disabled", "network_config.0.pod_cidr_overprovision_config.0.disabled", "true"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tier1_net", "network_config.0.network_performance_config.#", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tier1_net", "network_config.0.network_performance_config.0.total_egress_bandwidth_tier", "TIER_1"),
				),
			},
			{
				ResourceName:            "google_container_node_pool.with_manual_pod_cidr",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_config.0.create_pod_range"},
			},
			{
				ResourceName:            "google_container_node_pool.with_auto_pod_cidr",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_config.0.create_pod_range"},
			},
			// edit the updateable network config
			{
				Config: testAccContainerNodePool_withNetworkConfig(cluster, np, network, "TIER_UNSPECIFIED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_tier1_net", "network_config.0.network_performance_config.#", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tier1_net", "network_config.0.network_performance_config.0.total_egress_bandwidth_tier", "TIER_UNSPECIFIED"),
				),
			},
		},
	})
}

func TestAccContainerNodePool_withMultiNicNetworkConfig(t *testing.T) {
	t.Parallel()

	randstr := acctest.RandString(t, 10)
	cluster := fmt.Sprintf("tf-test-cluster-%s", randstr)
	np := fmt.Sprintf("tf-test-np-%s", randstr)
	network := fmt.Sprintf("tf-test-net-%s", randstr)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withMultiNicNetworkConfig(cluster, np, network),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_config.0.create_pod_range", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerNodePool_withEnablePrivateNodesToggle(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withEnablePrivateNodesToggle(cluster, np, network, "true"),
			},
			{
				ResourceName:            "google_container_node_pool.with_enable_private_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerNodePool_withEnablePrivateNodesToggle(cluster, np, network, "false"),
			},
			{
				ResourceName:            "google_container_node_pool.with_enable_private_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func testAccContainerNodePool_withEnablePrivateNodesToggle(cluster, np, network, flag string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = "1.27"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
  deletion_protection = false
}

resource "google_container_node_pool" "with_enable_private_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  network_config {
    create_pod_range = false
    enable_private_nodes = %s
    pod_range = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
  }
  node_config {
	oauth_scopes = [
	  "https://www.googleapis.com/auth/cloud-platform",
	]
  }
}
`, network, cluster, np, flag)
}

func TestAccContainerNodePool_withUpgradeSettings(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, networkName, subnetworkName, 2, 3, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, networkName, subnetworkName, 2, 1, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, networkName, subnetworkName, 1, 1, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, networkName, subnetworkName, 0, 0, "BLUE_GREEN", "100s", 1, 0.0, "0s"),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, networkName, subnetworkName, 0, 0, "BLUE_GREEN", "100s", 0, 0.5, "1s"),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withGPU(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withGPU(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_gpu",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withManagement(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	management := `
	management {
		auto_repair = "false"
		auto_upgrade = "false"
	}`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, "", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.#", "1"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "true"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_upgrade", "true"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_management",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, management, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.#", "1"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "false"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_upgrade", "false"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_management",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfigScopeAlias(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNodeConfigScopeAlias(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_node_config_scope_alias",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This test exists to validate a regional node pool *and* and update to it.
func TestAccContainerNodePool_regionalAutoscaling(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_regionalAutoscaling(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "3"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "5"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_basic(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
		},
	})
}

// This test exists to validate a node pool with total size *and* and update to it.
func TestAccContainerNodePool_totalSize(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_totalSize(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.total_min_node_count", "4"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.total_max_node_count", "12"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.location_policy", "BALANCED"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_updateTotalSize(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.total_min_node_count", "2"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.total_max_node_count", "22"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.location_policy", "ANY"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_basicTotalSize(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
		},
	})
}

func TestAccContainerNodePool_autoscaling(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_autoscaling(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "3"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "5"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_basic(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
		},
	})
}

func TestAccContainerNodePool_resize(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_additionalZones(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "node_count", "2"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_resize(cluster, np, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "node_count", "3"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_version(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_version(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_updateVersion(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_version(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_regionalClusters(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_regionalClusters(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_012_ConfigModeAttr(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_012_ConfigModeAttr1(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_012_ConfigModeAttr2(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test guest_accelerator.count = 0 is the same as guest_accelerator = []
				Config:             testAccContainerNodePool_EmptyGuestAccelerator(cluster, np, networkName, subnetworkName),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccContainerNodePool_EmptyGuestAccelerator(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test alternative way to specify an empty node pool
				Config: testAccContainerNodePool_EmptyGuestAccelerator(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test alternative way to specify an empty node pool
				Config: testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np, networkName, subnetworkName, 1),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Assert that changes in count from 1 result in a diff
				Config:             testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np, networkName, subnetworkName, 2),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Assert that adding another accelerator block will also result in a diff
				Config:             testAccContainerNodePool_PartialEmptyGuestAccelerator2(cluster, np, networkName, subnetworkName),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccContainerNodePool_shieldedInstanceConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_shieldedInstanceConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_node_pool.np",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_pods_per_node"},
			},
		},
	})
}

func TestAccContainerNodePool_concurrent(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np1 := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	np2 := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_concurrentCreate(cluster, np1, np2, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_node_pool.np2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_concurrentUpdate(cluster, np1, np2, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_node_pool.np2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withSoleTenantConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withSoleTenantConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_sole_tenant_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_ephemeralStorageLocalSsdConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_ephemeralStorageLocalSsdConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_ephemeralStorageLocalSsdConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	location       = "us-central1-a"
	// this feature became available in 1.25.3-gke.1800, not sure if theres a better way to do
	version_prefix = "1.25"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-1"
    ephemeral_storage_local_ssd_config {
      local_ssd_count = 1
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_localNvmeSsdBlockConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_localNvmeSsdBlockConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_localNvmeSsdBlockConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	location       = "us-central1-a"
	// this feature became available in 1.25.3-gke.1800, not sure if theres a better way to do
	version_prefix = "1.25"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-1"
    local_nvme_ssd_block_config {
      local_ssd_count = 1
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_secondaryBootDisks(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_secondaryBootDisks(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_node_pool.np-no-mode",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_secondaryBootDisks(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  min_master_version = "1.28"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
	gcfs_config {
  		enabled = true
	}
    secondary_boot_disks {
      disk_image = ""
      mode = "CONTAINER_IMAGE_CACHE"
    }
  }
}

resource "google_container_node_pool" "np-no-mode" {
  name               = "%s-no-mode"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
	gcfs_config {
  		enabled = true
	}
    secondary_boot_disks {
      disk_image = ""
    }
  }
}
`, cluster, networkName, subnetworkName, np, np)
}

func TestAccContainerNodePool_gcfsConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_gcfsConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_gcfsConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    gcfs_config {
      enabled = true
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_gvnic(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_gvnic(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_gvnic(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    gvnic {
      enabled = true
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_fastSocket(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_fastSocket(cluster, np, networkName, subnetworkName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np",
						"node_config.0.fast_socket.0.enabled", "true"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_fastSocket(cluster, np, networkName, subnetworkName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np",
						"node_config.0.fast_socket.0.enabled", "false"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_fastSocket(cluster, np, networkName, subnetworkName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1
  min_master_version = "1.28"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    guest_accelerator {
      type  = "nvidia-tesla-t4"
      count = 1
      }
    gvnic {
      enabled = true
    }
    fast_socket {
      enabled = %t
    }
  }
}
`, cluster, networkName, subnetworkName, np, enabled)
}

func TestAccContainerNodePool_compactPlacement(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_compactPlacement(cluster, np, "COMPACT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerNodePool_compactPlacement(cluster, np, placementType, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
    machine_type = "c2-standard-4"
  }
  placement_policy {
    type = "%s"
  }
}
`, cluster, networkName, subnetworkName, np, placementType)
}

func TestAccContainerNodePool_customPlacementPolicy(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	policy := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_customPlacementPolicy(cluster, np, policy, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "node_config.0.machine_type", "c2-standard-4"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "placement_policy.0.policy_name", policy),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "placement_policy.0.type", "COMPACT"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_customPlacementPolicy(cluster, np, policyName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_compute_resource_policy" "policy" {
  name = "%s"
  region = "us-central1"
  group_placement_policy {
    collocation = "COLLOCATED"
  }
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
	autoscaling {}

  node_config {
    machine_type = "c2-standard-4"
  }
  placement_policy {
	type = "COMPACT"
    policy_name = google_compute_resource_policy.policy.name
  }
}
`, cluster, networkName, subnetworkName, policyName, np)
}

func TestAccContainerNodePool_enableQueuedProvisioning(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_enableQueuedProvisioning(cluster, np, networkName, subnetworkName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "node_config.0.machine_type", "n1-standard-2"),
					resource.TestCheckResourceAttr("google_container_node_pool.np",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "NO_RESERVATION"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "queued_provisioning.0.enabled", "true"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_enableQueuedProvisioning(cluster, np, networkName, subnetworkName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name                = "%s"
  location            = "us-central1-a"
  initial_node_count  = 1
  min_master_version  = "1.28"
  deletion_protection = false
  network             = "%s"
  subnetwork          = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  autoscaling {
    total_min_node_count = 0
    total_max_node_count = 1
  }

  node_config {
    machine_type = "n1-standard-2"
    guest_accelerator {
      type  = "nvidia-tesla-t4"
      count = 1
      gpu_driver_installation_config {
        gpu_driver_version = "LATEST"
      }
    }
    reservation_affinity {
      consume_reservation_type = "NO_RESERVATION"
    }
  }
  queued_provisioning {
    enabled = %t
  }
}
`, cluster, networkName, subnetworkName, np, enabled)
}

func TestAccContainerNodePool_threadsPerCore(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_threadsPerCore(cluster, np, networkName, subnetworkName, 1),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerNodePool_threadsPerCore(cluster, np, networkName, subnetworkName string, threadsPerCore int) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"

  node_config {
    machine_type = "c2-standard-4"
	advanced_machine_features {
		threads_per_core = "%v"
	}
  }
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
    machine_type = "c2-standard-4"
	advanced_machine_features {
		threads_per_core = "%v"
	}
  }
}
`, cluster, networkName, subnetworkName, threadsPerCore, np, threadsPerCore)
}

func TestAccContainerNodePool_nestedVirtualization(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_nestedVirtualization(cluster, np, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerNodePool_nestedVirtualization(cluster, np, networkName, subnetworkName string, enableNV bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"

  node_config {
    machine_type = "c2-standard-4"
    advanced_machine_features {
      threads_per_core = 1
      enable_nested_virtualization = "%t"
    }
  }
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
    machine_type = "c2-standard-4"
    advanced_machine_features {
      threads_per_core = 1
      enable_nested_virtualization = "%t"
    }
  }
}
`, cluster, networkName, subnetworkName, enableNV, np, enableNV)
}

func testAccCheckContainerNodePoolDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_container_node_pool" {
				continue
			}

			attributes := rs.Primary.Attributes
			location := attributes["location"]

			var err error
			if location != "" {
				_, err = config.NewContainerClient(config.UserAgent).Projects.Zones.Clusters.NodePools.Get(
					config.Project, attributes["location"], attributes["cluster"], attributes["name"]).Do()
			} else {
				name := fmt.Sprintf(
					"projects/%s/locations/%s/clusters/%s/nodePools/%s",
					config.Project,
					attributes["location"],
					attributes["cluster"],
					attributes["name"],
				)
				_, err = config.NewContainerClient(config.UserAgent).Projects.Locations.Clusters.NodePools.Get(name).Do()
			}

			if err == nil {
				return fmt.Errorf("NodePool still exists")
			}
		}

		return nil
	}
}

func testAccContainerNodePool_basic(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_container_cluster" "cluster" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withLoggingVariant(cluster, np, loggingVariant, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_logging_variant" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.with_logging_variant.name
  initial_node_count = 1
  node_config {
    logging_variant = "%s"
  }
}
`, cluster, networkName, subnetworkName, np, loggingVariant)
}

func testAccContainerNodePool_basicWithClusterId(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_container_cluster" "cluster" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  cluster            = google_container_cluster.cluster.id
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_nodeLocations(cluster, np, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }

  master_authorized_networks_config {
  }
  deletion_protection = false
}

resource "google_container_node_pool" "np" {
  name     = "%s"
  location = "us-central1"
  cluster  = google_container_cluster.cluster.name

  initial_node_count = 1
  node_locations     = ["us-central1-a", "us-central1-c"]
}
`, network, cluster, np)
}

func testAccContainerNodePool_maxPodsPerNode(cluster, np, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }

  master_authorized_networks_config {
  }
  deletion_protection = false
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  max_pods_per_node  = 30
  initial_node_count = 2
}
`, network, cluster, np)
}

func testAccContainerNodePool_regionalClusters(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  cluster            = google_container_cluster.cluster.name
  location           = "us-central1"
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_namePrefix(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name_prefix        = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_noName(cluster, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName)
}

func testAccContainerNodePool_regionalAutoscaling(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  autoscaling {
    min_node_count = 1
    max_node_count = 3
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_totalSize(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  min_master_version = "1.27"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  autoscaling {
    total_min_node_count = 4
    total_max_node_count = 12
    location_policy      = "BALANCED"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_updateTotalSize(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  min_master_version = "1.27"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  autoscaling {
    total_min_node_count = 2
    total_max_node_count = 22
    location_policy      = "ANY"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_basicTotalSize(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_container_cluster" "cluster" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  min_master_version = "1.27"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_autoscaling(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  autoscaling {
    min_node_count = 1
    max_node_count = 3
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_updateAutoscaling(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  autoscaling {
    min_node_count = 0
    max_node_count = 5
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_additionalZones(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name       = "%s"
  location   = "us-central1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 2
}
`, cluster, networkName, subnetworkName, nodePool)
}

func testAccContainerNodePool_resize(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name       = "%s"
  location   = "us-central1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 3
}
`, cluster, networkName, subnetworkName, nodePool)
}

func testAccContainerNodePool_withManagement(cluster, nodePool, management, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  release_channel {
	  channel = "UNSPECIFIED"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np_with_management" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  %s

  node_config {
    machine_type = "g1-small"
    disk_size_gb = 10
    oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
  }
}
`, cluster, networkName, subnetworkName, nodePool, management)
}

func testAccContainerNodePool_withNodeConfig(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np_with_node_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "g1-small"
    disk_size_gb = 10
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    preemptible      = true
    min_cpu_platform = "Intel Broadwell"

    taint {
      key    = "taint_key"
      value  = "taint_value"
      effect = "PREFER_NO_SCHEDULE"
    }

    taint {
      key    = "taint_key2"
      value  = "taint_value2"
      effect = "NO_EXECUTE"
    }

    // Updatable fields
    image_type = "COS_CONTAINERD"

    tags = ["foo"]

    labels = {
      "test.terraform.io/key1" = "foo"
    }

    resource_labels = {
      "key1" = "foo"
    }
  }
}
`, cluster, networkName, subnetworkName, nodePool)
}

func testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np_with_node_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "g1-small"
    disk_size_gb = 10
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    preemptible      = true
    min_cpu_platform = "Intel Broadwell"

    taint {
      key    = "taint_key"
      value  = "taint_value"
      effect = "PREFER_NO_SCHEDULE"
    }

    taint {
      key    = "taint_key2"
      value  = "taint_value2"
      effect = "NO_EXECUTE"
    }

    // Updatable fields
    image_type = "UBUNTU_CONTAINERD"

    tags = ["bar", "foobar"]

    labels = {
      "test.terraform.io/key1" = "bar"
      "test.terraform.io/key2" = "foo"
    }

    resource_labels = {
      "key1" = "bar"
      "key2" = "foo"
    }
  }
}
`, cluster, networkName, subnetworkName, nodePool)
}

func testAccContainerNodePool_withTaintsUpdate(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_container_cluster" "cluster" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
	taint {
      key    = "taint_key"
      value  = "taint_value"
      effect = "PREFER_NO_SCHEDULE"
    }
  }


}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withMachineAndDiskUpdate(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_container_cluster" "cluster" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
	machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withReservationAffinity(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_reservation_affinity" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type    = "n1-standard-1"
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    reservation_affinity {
      consume_reservation_type = "ANY_RESERVATION"
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withReservationAffinitySpecific(cluster, reservation, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_compute_reservation" "gce_reservation" {
  name = "%s"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type     = "n1-standard-1"
    }
  }

  specific_reservation_required = true
}

resource "google_container_node_pool" "with_reservation_affinity" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type    = "n1-standard-1"
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    reservation_affinity {
      consume_reservation_type = "SPECIFIC_RESERVATION"
      key = "compute.googleapis.com/reservation-name"
      values = [
        google_compute_reservation.gce_reservation.name
      ]
    }
  }
}
`, cluster, networkName, subnetworkName, reservation, np)
}

func testAccContainerNodePool_withWorkloadMetadataConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_workload_metadata_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    spot         = true
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    workload_metadata_config {
      mode = "GCE_METADATA"
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withWorkloadMetadataConfig_gkeMetadata(projectID, cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_workload_metadata_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}
`, projectID, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withKubeletConfig(cluster, np, policy, period, networkName, subnetworkName string, quota bool, podPidsLimit int) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

# cpu_manager_policy & cpu_cfs_quota_period cannot be blank if cpu_cfs_quota is set to true
# cpu_manager_policy & cpu_cfs_quota_period must not set if cpu_cfs_quota is set to false
resource "google_container_node_pool" "with_kubelet_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    image_type = "COS_CONTAINERD"
    kubelet_config {
      cpu_manager_policy   = %q
      cpu_cfs_quota        = %v
      cpu_cfs_quota_period = %q
      pod_pids_limit			 = %d
    }
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    logging_variant = "DEFAULT"
  }
}
`, cluster, networkName, subnetworkName, np, policy, quota, period, podPidsLimit)
}

func testAccContainerNodePool_withLinuxNodeConfig(cluster, np, tcpMem, networkName, subnetworkName string) string {
	linuxNodeConfig := `
    linux_node_config {
      sysctls = {}
    }
`
	if len(tcpMem) != 0 {
		linuxNodeConfig = fmt.Sprintf(`
    linux_node_config {
      sysctls = {
        "net.core.netdev_max_backlog" = "10000"
        "net.core.rmem_max"           = 10000
        "net.core.wmem_default"       = 10000
        "net.core.wmem_max"           = 20000
        "net.core.optmem_max"         = 10000
        "net.core.somaxconn"          = 12800
        "net.ipv4.tcp_rmem"           = "%s"
        "net.ipv4.tcp_wmem"           = "%s"
        "net.ipv4.tcp_tw_reuse"       = 1
      }
    }
`, tcpMem, tcpMem)
	}

	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_linux_node_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    image_type = "COS_CONTAINERD"
    %s
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
`, cluster, networkName, subnetworkName, np, linuxNodeConfig)
}

func testAccContainerNodePool_withCgroupMode(cluster, np, mode, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name                = "%s"
  location            = "us-central1-a"
  initial_node_count  = 1
  min_master_version  = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    image_type = "COS_CONTAINERD"
    linux_node_config {
      cgroup_mode = "%s"
    }
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
`, cluster, networkName, subnetworkName, np, mode)
}

func testAccContainerNodePool_withNetworkConfig(cluster, np, network, netTier string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }

  secondary_ip_range {
    range_name    = "another-pod"
    ip_cidr_range = "10.1.32.0/22"
  }

  lifecycle {
    ignore_changes = [
      # The auto nodepool creates a secondary range which diffs this resource.
      secondary_ip_range,
    ]
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
  release_channel {
	channel = "RAPID"
  }
  deletion_protection = false
}

resource "google_container_node_pool" "with_manual_pod_cidr" {
  name               = "%s-manual"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  network_config {
    create_pod_range = false
    pod_range = google_compute_subnetwork.container_subnetwork.secondary_ip_range[2].range_name
  }
  node_config {
	oauth_scopes = [
	  "https://www.googleapis.com/auth/cloud-platform",
	]
  }
}

resource "google_container_node_pool" "with_auto_pod_cidr" {
  name               = "%s-auto"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  network_config {
	create_pod_range    = true
	pod_range           = "auto-pod-range"
	pod_ipv4_cidr_block = "10.2.0.0/20"
  }
  node_config {
	oauth_scopes = [
	  "https://www.googleapis.com/auth/cloud-platform",
	]
  }
}

resource "google_container_node_pool" "with_pco_disabled" {
  name               = "%s-pco"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  network_config {
	pod_cidr_overprovision_config {
		disabled = true
	}
  }
  node_config {
	oauth_scopes = [
	  "https://www.googleapis.com/auth/cloud-platform",
	]
  }
}

resource "google_container_node_pool" "with_tier1_net" {
  name               = "%s-tier1"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  node_locations = [
	"us-central1-a",
  ]
  network_config {
	network_performance_config {
		total_egress_bandwidth_tier = "%s"
	}
  }
  node_config {
	machine_type = "n2-standard-32"
	gvnic {
		enabled = true
	}
	oauth_scopes = [
		"https://www.googleapis.com/auth/cloud-platform",
	]
  }
}

`, network, cluster, np, np, np, np, netTier)
}

func testAccContainerNodePool_withMultiNicNetworkConfig(cluster, np, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s-1"
  auto_create_subnetworks = false
}

resource "google_compute_network" "addn_net_1" {
  name                    = "%s-2"
  auto_create_subnetworks = false
}

resource "google_compute_network" "addn_net_2" {
  name                    = "%s-3"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%s-subnet-1"
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }

  lifecycle {
    ignore_changes = [
      # The auto nodepool creates a secondary range which diffs this resource.
      secondary_ip_range,
    ]
  }
}

resource "google_compute_subnetwork" "subnet1" {
  name                     = "%s-subnet-2"
  network                  = google_compute_network.addn_net_1.name
  ip_cidr_range            = "10.0.37.0/24"
  region                   = "us-central1"
}

resource "google_compute_subnetwork" "subnet2" {
  name                     = "%s-subnet-3"
  network                  = google_compute_network.addn_net_2.name
  ip_cidr_range            = "10.0.38.0/24"
  region                   = "us-central1"

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.64.0/19"
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
  private_cluster_config {
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }
  release_channel {
	channel = "RAPID"
  }
  enable_multi_networking = true
  datapath_provider = "ADVANCED_DATAPATH"
  deletion_protection = false
}

resource "google_container_node_pool" "with_multi_nic" {
  name               = "%s-mutli-nic"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  node_count = 1
  network_config {
    create_pod_range = false
    enable_private_nodes = true
    pod_range = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    additional_node_network_configs {
      network    = google_compute_network.addn_net_1.name
      subnetwork = google_compute_subnetwork.subnet1.name
    }
    additional_node_network_configs {
      network    = google_compute_network.addn_net_2.name
      subnetwork = google_compute_subnetwork.subnet2.name
    }
    additional_pod_network_configs {
      subnetwork          = google_compute_subnetwork.subnet2.name
      secondary_pod_range = "pod"
      max_pods_per_node   = 32
    }
  }
  node_config {
    machine_type = "n2-standard-8"
	oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
    image_type = "COS_CONTAINERD"
  }
}

`, network, network, network, network, network, network, cluster, np)
}

func makeUpgradeSettings(maxSurge int, maxUnavailable int, strategy string, nodePoolSoakDuration string, batchNodeCount int, batchPercentage float64, batchSoakDuration string) string {
	if strategy == "BLUE_GREEN" {
		return fmt.Sprintf(`
upgrade_settings {
	strategy = "%s"
	blue_green_settings {
		node_pool_soak_duration = "%s"
		standard_rollout_policy {
			batch_node_count = %d
			batch_percentage = %f
			batch_soak_duration = "%s"
		}
	}
}
`, strategy, nodePoolSoakDuration, batchNodeCount, batchPercentage, batchSoakDuration)
	}
	return fmt.Sprintf(`
upgrade_settings {
	max_surge = %d
	max_unavailable = %d
	strategy = "%s"
}
`, maxSurge, maxUnavailable, strategy)
}

func testAccContainerNodePool_withUpgradeSettings(clusterName, nodePoolName, networkName, subnetworkName string, maxSurge int, maxUnavailable int, strategy string, nodePoolSoakDuration string, batchNodeCount int, batchPercentage float64, batchSoakDuration string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1" {
  location = "us-central1"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1
  min_master_version = "${data.google_container_engine_versions.central1.latest_master_version}"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_upgrade_settings" {
  name = "%s"
  location = "us-central1"
  cluster = "${google_container_cluster.cluster.name}"
  initial_node_count = 1
  %s
}
`, clusterName, networkName, subnetworkName, nodePoolName, makeUpgradeSettings(maxSurge, maxUnavailable, strategy, nodePoolSoakDuration, batchNodeCount, batchPercentage, batchSoakDuration))
}

func testAccContainerNodePool_withGPU(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1c" {
  location = "us-central1-c"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-c"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1c.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np_with_gpu" {
  name     = "%s"
  location = "us-central1-c"
  cluster  = google_container_cluster.cluster.name

  initial_node_count = 1

  node_config {
    machine_type = "a2-highgpu-1g"  // can't be e2 because of accelerator
    disk_size_gb = 32

    oauth_scopes = [
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
    ]

    preemptible     = true
    service_account = "default"
    image_type      = "COS_CONTAINERD"

    guest_accelerator {
      type  = "nvidia-tesla-a100"
      gpu_partition_size = "1g.5gb"
      count = 1
	  gpu_driver_installation_config {
		gpu_driver_version = "LATEST"
	  }
      gpu_sharing_config {
        gpu_sharing_strategy = "TIME_SHARING"
        max_shared_clients_per_gpu = 2
      }
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_withNodeConfigScopeAlias(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np_with_node_config_scope_alias" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "g1-small"
    disk_size_gb = 10
    oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_version(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  version = data.google_container_engine_versions.central1a.valid_node_versions[1]
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_updateVersion(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  version = data.google_container_engine_versions.central1a.valid_node_versions[0]
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_012_ConfigModeAttr1(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator {
      count = 1
      type  = "nvidia-tesla-t4"
    }
	machine_type = "n1-highmem-4"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_012_ConfigModeAttr2(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator = []
	machine_type = "n1-highmem-4"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_EmptyGuestAccelerator(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator {
      count = 0
      type  = "nvidia-tesla-p100"
    }
	machine_type = "n1-highmem-4"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np, networkName, subnetworkName string, count int) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator {
      count = 0
      type  = "nvidia-tesla-p100"
    }

    guest_accelerator {
      count = %d
      type  = "nvidia-tesla-p100"
    }
	machine_type = "n1-highmem-4"
  }
}
`, cluster, networkName, subnetworkName, np, count)
}

func testAccContainerNodePool_PartialEmptyGuestAccelerator2(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator {
      count = 0
      type  = "nvidia-tesla-p100"
    }

    guest_accelerator {
      count = 1
      type  = "nvidia-tesla-p100"
    }

    guest_accelerator {
      count = 1
      type  = "nvidia-tesla-p9000"
    }
	machine_type = "n1-highmem-4"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_shieldedInstanceConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  node_config {
    shielded_instance_config {
      enable_integrity_monitoring = true
      enable_secure_boot          = true
    }
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_concurrentCreate(cluster, np1, np2, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np1" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}

resource "google_container_node_pool" "np2" {
	name               = "%s"
	location           = "us-central1-a"
	cluster            = google_container_cluster.cluster.name
	initial_node_count = 2
  }
`, cluster, networkName, subnetworkName, np1, np2)
}

func testAccContainerNodePool_concurrentUpdate(cluster, np1, np2, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np1" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  version            = "1.29.4-gke.1043002"
}

resource "google_container_node_pool" "np2" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  version            = "1.29.4-gke.1043002"
}
`, cluster, networkName, subnetworkName, np1, np2)
}

func testAccContainerNodePool_withSoleTenantConfig(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_compute_node_template" "soletenant-tmpl" {
  name      = "tf-test-soletenant-tmpl"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "nodes" {
  name        = "tf-test-soletenant-group"
  zone        = "us-central1-a"
  initial_size	= 1
  node_template = google_compute_node_template.soletenant-tmpl.id
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_sole_tenant_config" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
       machine_type = "n1-standard-2"
       sole_tenant_config {
			node_affinity {
				key = "compute.googleapis.com/node-group-name"
				operator = "IN"
				values = [google_compute_node_group.nodes.name]
			}
       }
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
`, cluster, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_withConfidentialNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withConfidentialNodes(clusterName, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_disableConfidentialNodes(clusterName, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withConfidentialNodes(clusterName, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_withConfidentialNodes(clusterName, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  node_config {
    confidential_nodes {
      enabled = false
    }
    machine_type = "n2-standard-2"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "n2d-standard-2" // can't be e2 because Confidential Nodes require AMD CPUs
    confidential_nodes {
      enabled = true
    }
  }
}
`, clusterName, networkName, subnetworkName, np)
}

func testAccContainerNodePool_disableConfidentialNodes(clusterName, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  node_config {
    confidential_nodes {
      enabled = false
    }
    machine_type = "n2-standard-2"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "n2d-standard-2" // can't be e2 because Confidential Nodes require AMD CPUs
    confidential_nodes {
      enabled = false
    }
  }
}
`, clusterName, networkName, subnetworkName, np)
}

func TestAccContainerNodePool_tpuTopology(t *testing.T) {
	t.Parallel()
	t.Skip("https://github.com/hashicorp/terraform-provider-google/issues/15254#issuecomment-1646277473")

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np1 := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	np2 := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_tpuTopology(cluster, np1, np2, "2x2x2", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.regular_pool", "node_config.0.machine_type", "n1-standard-4"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tpu_topology", "node_config.0.machine_type", "ct4p-hightpu-4t"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tpu_topology", "placement_policy.0.tpu_topology", "2x2x2"),
					resource.TestCheckResourceAttr("google_container_node_pool.with_tpu_topology", "placement_policy.0.type", "COMPACT"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_tpu_topology",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_tpuTopology(cluster, np1, np2, tpuTopology, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "regular_pool" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-4"
  }
}

resource "google_container_node_pool" "with_tpu_topology" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2

  node_config {
    machine_type = "ct4p-hightpu-4t"

  }
  placement_policy {
  type = "COMPACT"
  tpu_topology = "%s"
  }
}
`, cluster, networkName, subnetworkName, np1, np2, tpuTopology)
}

func TestAccContainerNodePool_withConfidentialBootDisk(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withConfidentialBootDisk(cluster, np, kms.CryptoKey.Name, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.with_confidential_boot_disk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_withConfidentialBootDisk(cluster, np string, kmsKeyName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "with_confidential_boot_disk" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name

  node_config {
    confidential_nodes {
      enabled = true
    }
    image_type = "COS_CONTAINERD"
    boot_disk_kms_key = "%s"
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    enable_confidential_storage = true
    machine_type = "n2d-standard-2"
    disk_type = "hyperdisk-balanced"
  }
}
`, cluster, networkName, subnetworkName, np, kmsKeyName)
}

func TestAccContainerNodePool_withoutConfidentialBootDisk(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withoutConfidentialBootDisk(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_node_pool.without_confidential_boot_disk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_withoutConfidentialBootDisk(cluster, np, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}

resource "google_container_node_pool" "without_confidential_boot_disk" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name

  node_config {
    image_type = "COS_CONTAINERD"
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
    enable_confidential_storage = false
    machine_type = "n2-standard-2"
    disk_type = "pd-balanced"
  }
}
`, cluster, networkName, subnetworkName, np)
}

func testAccContainerNodePool_resourceManagerTags(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_tags_tag_key.key1]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%[3]s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1

  deletion_protection = false
  network             = "%[4]s"
  subnetwork          = "%[5]s"

  timeouts {
    create = "30m"
    update = "40m"
  }

  depends_on = [time_sleep.wait_120_seconds]
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = google_container_cluster.primary.name
  location   = "us-central1-a"
  cluster    = google_container_cluster.primary.name

  version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15

    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key1.name}" = "tagValues/${google_tags_tag_value.value1.name}"
    }
  }
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func testAccContainerNodePool_resourceManagerTagsUpdate1(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_tags_tag_key.key1]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%[3]s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1

  deletion_protection = false
  network             = "%[4]s"
  subnetwork          = "%[5]s"

  timeouts {
    create = "30m"
    update = "40m"
  }

  depends_on = [time_sleep.wait_120_seconds]
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = google_container_cluster.primary.name
  location   = "us-central1-a"
  cluster    = google_container_cluster.primary.name

  version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15

    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key1.name}" = "tagValues/${google_tags_tag_value.value1.name}"
	  "tagKeys/${google_tags_tag_key.key2.name}" = "tagValues/${google_tags_tag_value.value2.name}"
    }
  }
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func testAccContainerNodePool_resourceManagerTagsUpdate2(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_tags_tag_key.key1]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%[3]s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1

  deletion_protection = false
  network             = "%[4]s"
  subnetwork          = "%[5]s"

  timeouts {
    create = "30m"
    update = "40m"
  }

  depends_on = [time_sleep.wait_120_seconds]
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = google_container_cluster.primary.name
  location   = "us-central1-a"
  cluster    = google_container_cluster.primary.name

  version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15
  }
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func TestAccContainerNodePool_privateRegistry(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodepool := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	secretID := fmt.Sprintf("tf-test-secret-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_privateRegistryEnabled(secretID, cluster, nodepool, networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np",
						"node_config.0.containerd_config.0.private_registry_access_config.0.enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np",
						"node_config.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.#",
						"2",
					),
					// First CA config
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np",
						"node_config.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.0.fqdns.0",
						"my.custom.domain",
					),
					// Second CA config
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np",
						"node_config.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.1.fqdns.0",
						"10.1.2.32",
					),
				),
			},
		},
	})
}

func testAccContainerNodePool_privateRegistryEnabled(secretID, cluster, nodepool, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_project" "test_project" { 
	}

resource "google_secret_manager_secret" "secret-basic" { 
	secret_id     = "%s" 
	replication { 
		user_managed { 
		replicas { 
			location = "us-central1" 
		} 
		} 
	} 
}

resource "google_secret_manager_secret_version" "secret-version-basic" { 
	secret = google_secret_manager_secret.secret-basic.id 
	secret_data = "dummypassword" 
  } 
   
resource "google_secret_manager_secret_iam_member" "secret_iam" { 
	secret_id  = google_secret_manager_secret.secret-basic.id 
	role       = "roles/secretmanager.admin" 
	member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com" 
	depends_on = [google_secret_manager_secret_version.secret-version-basic] 
  }

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
	
resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
	
  node_config {
	oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    containerd_config {
      private_registry_access_config {
        enabled = true
        certificate_authority_domain_config {
          fqdns = [ "my.custom.domain", "10.0.0.127:8888" ]
          gcp_secret_manager_certificate_config {
            secret_uri = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
        certificate_authority_domain_config {
          fqdns = [ "10.1.2.32" ]
          gcp_secret_manager_certificate_config {
            secret_uri = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
      }
    }
  }
}
`, secretID, cluster, network, subnetwork, nodepool)
}

func TestAccContainerNodePool_defaultDriverInstallation(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_defaultDriverInstallation(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_defaultDriverInstallation(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  deletion_protection = false

  min_master_version = data.google_container_engine_versions.central1a.release_channel_latest_version["RAPID"]
  release_channel {
    channel = "RAPID"
  }
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  version            = "1.30.1-gke.1329003"

  node_config {
    service_account = "default"
    machine_type = "n1-standard-8"

    guest_accelerator {
      type  = "nvidia-tesla-t4"
      count = 1
      gpu_sharing_config {
	    gpu_sharing_strategy = "TIME_SHARING"
	    max_shared_clients_per_gpu = 3
      }
    }
  }
}
`, cluster, np)
}
