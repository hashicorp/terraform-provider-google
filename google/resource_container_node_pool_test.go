package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_maxPodsPerNode(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_maxPodsPerNode(cluster, np),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_namePrefix(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_namePrefix(cluster, "tf-np-"),
			},
			resource.TestStep{
				ResourceName:            "google_container_node_pool.np",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccContainerNodePool_noName(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_noName(cluster),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	nodePool := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfig(cluster, nodePool),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np_with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
				// autoscaling.# = 0 is equivalent to no autoscaling at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"autoscaling.#"},
			},
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool),
			},
			resource.TestStep{
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

func TestAccContainerNodePool_withNodeConfigTaints(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfigTaints(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np_with_node_config", "node_config.0.taint.#", "2"),
				),
			},
			// Don't include an import step because beta features can't yet be imported.
			// Once taints are in GA, consider merging this test with the _withNodeConfig test.
		},
	})
}

func TestAccContainerNodePool_withWorkloadMetadataConfig(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withWorkloadMetadataConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.node_metadata", "SECURE"),
				),
			},
			{
				ResourceName:      "google_container_node_pool.with_workload_metadata_config",
				ImportState:       true,
				ImportStateVerify: true,
				// Import always uses the v1 API, so beta features don't get imported.
				ImportStateVerifyIgnore: []string{
					"node_config.0.workload_metadata_config.#",
					"node_config.0.workload_metadata_config.0.node_metadata",
				},
			},
		},
	})
}

func TestAccContainerNodePool_withGPU(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withGPU(),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np_with_gpu",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withManagement(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	nodePool := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	management := `
	management {
		auto_repair = "true"
		auto_upgrade = "true"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.#", "1"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "false"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "false"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np_with_management",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccContainerNodePool_withManagement(cluster, nodePool, management),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.#", "1"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "true"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "true"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np_with_management",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfigScopeAlias(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfigScopeAlias(),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np_with_node_config_scope_alias",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//This test exists to validate a regional node pool *and* and update to it.
func TestAccContainerNodePool_regionalAutoscaling(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_regionalAutoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "3"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "5"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count"),
				),
			},
			resource.TestStep{
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_autoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "3"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count", "5"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_node_pool.np", "autoscaling.0.max_node_count"),
				),
			},
			resource.TestStep{
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

func testAccCheckContainerNodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_node_pool" {
			continue
		}

		attributes := rs.Primary.Attributes
		zone := attributes["zone"]

		var err error
		if zone != "" {
			_, err = config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
				config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		} else {
			name := fmt.Sprintf(
				"projects/%s/locations/%s/clusters/%s/nodePools/%s",
				config.Project,
				attributes["region"],
				attributes["cluster"],
				attributes["name"],
			)
			_, err = config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Get(name).Do()
		}

		if err == nil {
			return fmt.Errorf("NodePool still exists")
		}
	}

	return nil
}

func TestAccContainerNodePool_regionalClusters(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_regionalClusters(cluster, np),
			},
			resource.TestStep{
				ResourceName:      "google_container_node_pool.np",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerNodePool_basic(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
}`, cluster, np)
}

func testAccContainerNodePool_maxPodsPerNode(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name                     = "${google_compute_network.container_network.name}"
	network                  = "${google_compute_network.container_network.name}"
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
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"
	private_cluster = true
	master_ipv4_cidr_block = "10.42.0.0/28"
	ip_allocation_policy {
		cluster_secondary_range_name  = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.0.range_name}"
		services_secondary_range_name = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.1.range_name}"
	}
	master_authorized_networks_config {
		cidr_blocks = []
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	max_pods_per_node = 30
	initial_node_count = 2
}`, cluster, cluster, np)
}

func testAccContainerNodePool_regionalClusters(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	region = "us-central1"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name = "%s"
	cluster = "${google_container_cluster.cluster.name}"
	region = "us-central1"
	initial_node_count = 2
}`, cluster, np)
}

func testAccContainerNodePool_namePrefix(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name_prefix = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
}`, cluster, np)
}

func testAccContainerNodePool_noName(cluster string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
}`, cluster)
}

func testAccContainerNodePool_regionalAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	region = "us-central1"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name = "%s"
	region = "us-central1"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 1
		max_node_count = 3
	}
}`, cluster, np)
}

func testAccContainerNodePool_autoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 1
		max_node_count = 3
	}
}`, cluster, np)
}

func testAccContainerNodePool_updateAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 0
		max_node_count = 5
	}
}`, cluster, np)
}

func testAccContainerNodePool_additionalZones(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	node_count = 2
}`, cluster, nodePool)
}

func testAccContainerNodePool_resize(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	node_count = 3
}`, cluster, nodePool)
}

func testAccContainerNodePool_withManagement(cluster, nodePool, management string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name               = "%s"
	zone               = "us-central1-a"
	initial_node_count = 1
}

resource "google_container_node_pool" "np_with_management" {
	name               = "%s"
	zone               = "us-central1-a"
	cluster            = "${google_container_cluster.cluster.name}"
	initial_node_count = 1

	%s

	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
	}
}`, cluster, nodePool, management)
}

func testAccContainerNodePool_withNodeConfig(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
}
resource "google_container_node_pool" "np_with_node_config" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = [
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring"
		]
		preemptible = true
		min_cpu_platform = "Intel Broadwell"

		// Updatable fields
		image_type = "COS"
	}
}`, cluster, nodePool)
}

func testAccContainerNodePool_withNodeConfigUpdate(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
}
resource "google_container_node_pool" "np_with_node_config" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = [
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring"
		]
		preemptible = true
		min_cpu_platform = "Intel Broadwell"

		// Updatable fields
		image_type = "UBUNTU"
	}
}`, cluster, nodePool)
}

func testAccContainerNodePool_withNodeConfigTaints() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1
}
resource "google_container_node_pool" "np_with_node_config" {
	name = "tf-nodepool-test-%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		taint {
			key = "taint_key"
			value = "taint_value"
			effect = "PREFER_NO_SCHEDULE"
		}
		taint {
			key = "taint_key2"
			value = "taint_value2"
			effect = "NO_EXECUTE"
		}
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerNodePool_withWorkloadMetadataConfig() string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  zone = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
  name               = "tf-cluster-nodepool-test-%s"
  zone               = "us-central1-a"
  initial_node_count = 1
  min_master_version = "${data.google_container_engine_versions.central1a.latest_master_version}"
}

resource "google_container_node_pool" "with_workload_metadata_config" {
  name = "tf-nodepool-test-%s"
  zone = "us-central1-a"
  cluster = "${google_container_cluster.cluster.name}"
  initial_node_count = 1
  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ]

    workload_metadata_config {
      node_metadata = "SECURE"
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerNodePool_withGPU() string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1c" {
	zone = "us-central1-c"
}

resource "google_container_cluster" "cluster" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-c"
	initial_node_count = 1
	node_version = "${data.google_container_engine_versions.central1c.latest_node_version}"
	min_master_version = "${data.google_container_engine_versions.central1c.latest_master_version}"
}
resource "google_container_node_pool" "np_with_gpu" {
	name = "tf-nodepool-test-%s"
	zone = "us-central1-c"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "n1-standard-1"
		disk_size_gb = 10
		oauth_scopes = [
			"https://www.googleapis.com/auth/devstorage.read_only", 
			"https://www.googleapis.com/auth/logging.write", 
			"https://www.googleapis.com/auth/monitoring", 
			"https://www.googleapis.com/auth/service.management.readonly",
			"https://www.googleapis.com/auth/servicecontrol", 
			"https://www.googleapis.com/auth/trace.append"
		]
		preemptible = true
		service_account = "default"
		image_type = "COS"
		guest_accelerator = [
			{
				type = "nvidia-tesla-k80"
				count = 1
			}
		]
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerNodePool_withNodeConfigScopeAlias() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1
}
resource "google_container_node_pool" "np_with_node_config_scope_alias" {
	name = "tf-nodepool-test-%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerNodePool_version(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	zone = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	min_master_version = "${data.google_container_engine_versions.central1a.latest_master_version}"
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1

	version = "${data.google_container_engine_versions.central1a.valid_node_versions.1}"
}`, cluster, np)
}

func testAccContainerNodePool_updateVersion(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	zone = "us-central1-a"
}

resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	min_master_version = "${data.google_container_engine_versions.central1a.latest_master_version}"
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1

	version = "${data.google_container_engine_versions.central1a.valid_node_versions.0}"
}`, cluster, np)
}
