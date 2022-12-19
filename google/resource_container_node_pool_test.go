package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basic(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_basicWithClusterId(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_basicWithClusterId(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
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
	skipIfVcr(t)
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_namePrefix(cluster, "tf-np-"),
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
	skipIfVcr(t)
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_noName(cluster),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "DEFAULT"),
			},
			{
				ResourceName:      "google_container_node_pool.with_logging_variant",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "MAX_THROUGHPUT"),
			},
			{
				ResourceName:      "google_container_node_pool.with_logging_variant",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withLoggingVariant(cluster, nodePool, "DEFAULT"),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNodeConfig(cluster, nodePool),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
			{
				Config: testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool),
			},
			{
				ResourceName:      "google_container_node_pool.np_with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
		},
	})
}

func TestAccContainerNodePool_withReservationAffinity(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withReservationAffinity(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	reservation := fmt.Sprintf("tf-test-reservation-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withReservationAffinitySpecific(cluster, reservation, np),
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

	pid := getTestProjectFromEnv()
	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withWorkloadMetadataConfig(cluster, np),
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
				Config: testAccContainerNodePool_withWorkloadMetadataConfig_gkeMetadata(pid, cluster, np),
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

func TestAccContainerNodePool_withNetworkConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNetworkConfig(cluster, np, network),
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
		},
	})
}

func TestAccContainerNodePool_withEnablePrivateNodesToggle(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
  min_master_version = "1.23"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, 2, 3, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, 2, 1, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, 1, 1, "SURGE", "", 0, 0.0, ""),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, 0, 0, "BLUE_GREEN", "100s", 1, 0.0, "0s"),
			},
			{
				ResourceName:      "google_container_node_pool.with_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_withUpgradeSettings(cluster, np, 0, 0, "BLUE_GREEN", "100s", 0, 0.5, "1s"),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withGPU(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	nodePool := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))
	management := `
	management {
		auto_repair = "false"
		auto_upgrade = "false"
	}`

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, ""),
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
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, management),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNodeConfigScopeAlias(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_regionalAutoscaling(cluster, np),
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
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
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
				Config: testAccContainerNodePool_basic(cluster, np),
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

//This test exists to validate a node pool with total size *and* and update to it.
func TestAccContainerNodePool_totalSize(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_totalSize(cluster, np),
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
				Config: testAccContainerNodePool_updateTotalSize(cluster, np),
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
				Config: testAccContainerNodePool_basicTotalSize(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_autoscaling(cluster, np),
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
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
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
				Config: testAccContainerNodePool_basic(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_additionalZones(cluster, np),
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
				Config: testAccContainerNodePool_resize(cluster, np),
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
	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_version(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_updateVersion(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_version(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_regionalClusters(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_012_ConfigModeAttr1(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerNodePool_012_ConfigModeAttr2(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_EmptyGuestAccelerator(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test alternative way to specify an empty node pool
				Config: testAccContainerNodePool_EmptyGuestAccelerator(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test alternative way to specify an empty node pool
				Config: testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np, 1),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Assert that changes in count from 1 result in a diff
				Config:             testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np, 2),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Assert that adding another accelerator block will also result in a diff
				Config:             testAccContainerNodePool_PartialEmptyGuestAccelerator2(cluster, np),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccContainerNodePool_shieldedInstanceConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_shieldedInstanceConfig(cluster, np),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np1 := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))
	np2 := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_concurrentCreate(cluster, np1, np2),
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
				Config: testAccContainerNodePool_concurrentUpdate(cluster, np1, np2),
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

func TestAccContainerNodePool_gcfsConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_gcfsConfig(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_gcfsConfig(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, np)
}

func TestAccContainerNodePool_gvnic(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_gvnic(cluster, np),
			},
			{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_gvnic(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, np)
}

func testAccCheckContainerNodePoolDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_container_node_pool" {
				continue
			}

			attributes := rs.Primary.Attributes
			location := attributes["location"]

			var err error
			if location != "" {
				_, err = config.NewContainerClient(config.userAgent).Projects.Zones.Clusters.NodePools.Get(
					config.Project, attributes["location"], attributes["cluster"], attributes["name"]).Do()
			} else {
				name := fmt.Sprintf(
					"projects/%s/locations/%s/clusters/%s/nodePools/%s",
					config.Project,
					attributes["location"],
					attributes["cluster"],
					attributes["name"],
				)
				_, err = config.NewContainerClient(config.userAgent).Projects.Locations.Clusters.NodePools.Get(name).Do()
			}

			if err == nil {
				return fmt.Errorf("NodePool still exists")
			}
		}

		return nil
	}
}

func testAccContainerNodePool_basic(cluster, np string) string {
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
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, np)
}

func testAccContainerNodePool_withLoggingVariant(cluster, np, loggingVariant string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, np, loggingVariant)
}

func testAccContainerNodePool_basicWithClusterId(cluster, np string) string {
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
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  cluster            = google_container_cluster.cluster.id
  initial_node_count = 2
}
`, cluster, np)
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

func testAccContainerNodePool_regionalClusters(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  cluster            = google_container_cluster.cluster.name
  location           = "us-central1"
  initial_node_count = 2
}
`, cluster, np)
}

func testAccContainerNodePool_namePrefix(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
}

resource "google_container_node_pool" "np" {
  name_prefix        = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, np)
}

func testAccContainerNodePool_noName(cluster string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
}

resource "google_container_node_pool" "np" {
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster)
}

func testAccContainerNodePool_regionalAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
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
`, cluster, np)
}

func testAccContainerNodePool_totalSize(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  min_master_version = "1.24"
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
`, cluster, np)
}

func testAccContainerNodePool_updateTotalSize(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 3
  min_master_version = "1.24"
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
`, cluster, np)
}

func testAccContainerNodePool_basicTotalSize(cluster, np string) string {
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
  min_master_version = "1.24"
}

resource "google_container_node_pool" "np" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
}
`, cluster, np)
}

func testAccContainerNodePool_autoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
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
`, cluster, np)
}

func testAccContainerNodePool_updateAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
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
`, cluster, np)
}

func testAccContainerNodePool_additionalZones(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]
}

resource "google_container_node_pool" "np" {
  name       = "%s"
  location   = "us-central1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 2
}
`, cluster, nodePool)
}

func testAccContainerNodePool_resize(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]
}

resource "google_container_node_pool" "np" {
  name       = "%s"
  location   = "us-central1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 3
}
`, cluster, nodePool)
}

func testAccContainerNodePool_withManagement(cluster, nodePool, management string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  release_channel {
	  channel = "UNSPECIFIED"
  }
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
`, cluster, nodePool, management)
}

func testAccContainerNodePool_withNodeConfig(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, nodePool)
}

func testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, nodePool)
}

func testAccContainerNodePool_withReservationAffinity(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
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
`, cluster, np)
}

func testAccContainerNodePool_withReservationAffinitySpecific(cluster, reservation, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
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
`, cluster, reservation, np)
}

func testAccContainerNodePool_withWorkloadMetadataConfig(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
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
`, cluster, np)
}

func testAccContainerNodePool_withWorkloadMetadataConfig_gkeMetadata(projectID, cluster, np string) string {
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
`, projectID, cluster, np)
}

func testAccContainerNodePool_withNetworkConfig(cluster, np, network string) string {
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

`, network, cluster, np, np)
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

func testAccContainerNodePool_withUpgradeSettings(clusterName string, nodePoolName string, maxSurge int, maxUnavailable int, strategy string, nodePoolSoakDuration string, batchNodeCount int, batchPercentage float64, batchSoakDuration string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1" {
  location = "us-central1"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1
  min_master_version = "${data.google_container_engine_versions.central1.latest_master_version}"
}

resource "google_container_node_pool" "with_upgrade_settings" {
  name = "%s"
  location = "us-central1"
  cluster = "${google_container_cluster.cluster.name}"
  initial_node_count = 1
  %s
}
`, clusterName, nodePoolName, makeUpgradeSettings(maxSurge, maxUnavailable, strategy, nodePoolSoakDuration, batchNodeCount, batchPercentage, batchSoakDuration))
}

func testAccContainerNodePool_withGPU(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1c" {
  location = "us-central1-c"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-c"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1c.latest_master_version
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
      gpu_sharing_config {
        gpu_sharing_strategy = "TIME_SHARING"
        max_shared_clients_per_gpu = 2
      }
    }
  }
}
`, cluster, np)
}

func testAccContainerNodePool_withNodeConfigScopeAlias(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, np)
}

func testAccContainerNodePool_version(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  version = data.google_container_engine_versions.central1a.valid_node_versions[1]
}
`, cluster, np)
}

func testAccContainerNodePool_updateVersion(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  version = data.google_container_engine_versions.central1a.valid_node_versions[0]
}
`, cluster, np)
}

func testAccContainerNodePool_012_ConfigModeAttr1(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator {
      count = 1
      type  = "nvidia-tesla-p100"
    }
	machine_type = "n1-highmem-4"
  }
}
`, cluster, np)
}

func testAccContainerNodePool_012_ConfigModeAttr2(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
}

resource "google_container_node_pool" "np" {
  name               = "%s"
  location           = "us-central1-f"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  node_config {
    guest_accelerator = []
  }
}
`, cluster, np)
}

func testAccContainerNodePool_EmptyGuestAccelerator(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
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
`, cluster, np)
}

func testAccContainerNodePool_PartialEmptyGuestAccelerator(cluster, np string, count int) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
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
`, cluster, np, count)
}

func testAccContainerNodePool_PartialEmptyGuestAccelerator2(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 3
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
`, cluster, np)
}

func testAccContainerNodePool_shieldedInstanceConfig(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
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
`, cluster, np)
}

func testAccContainerNodePool_concurrentCreate(cluster, np1, np2 string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
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
`, cluster, np1, np2)
}

func testAccContainerNodePool_concurrentUpdate(cluster, np1, np2 string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
}

resource "google_container_node_pool" "np1" {
  name               = "%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 2
  version 		     = "1.23.13-gke.900"
}

resource "google_container_node_pool" "np2" {
	name               = "%s"
	location           = "us-central1-a"
	cluster            = google_container_cluster.cluster.name
	initial_node_count = 2
	version 		   = "1.23.13-gke.900"
  }
`, cluster, np1, np2)
}
