package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccContainerNodePool_maxPodsPerNode(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_maxPodsPerNode(cluster, np),
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
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

func TestAccContainerNodePool_withNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	nodePool := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

func TestAccContainerNodePool_withGPU(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withGPU(),
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
			{
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
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "true"),
					resource.TestCheckResourceAttr(
						"google_container_node_pool.np_with_management", "management.0.auto_repair", "true"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePool_withNodeConfigScopeAlias(),
			},
			{
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

func TestAccContainerNodePool_autoscaling(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

func TestAccContainerNodePool_regionalClusters(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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

func testAccCheckContainerNodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_node_pool" {
			continue
		}

		attributes := rs.Primary.Attributes
		location := attributes["location"]

		var err error
		if location != "" {
			_, err = config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
				config.Project, attributes["location"], attributes["cluster"], attributes["name"]).Do()
		} else {
			name := fmt.Sprintf(
				"projects/%s/locations/%s/clusters/%s/nodePools/%s",
				config.Project,
				attributes["location"],
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

func testAccContainerNodePool_basic(cluster, np string) string {
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
}
`, cluster, np)
}

func testAccContainerNodePool_maxPodsPerNode(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
`, cluster, cluster, np)
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
    image_type = "COS"
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
    image_type = "UBUNTU"
  }
}
`, cluster, nodePool)
}

func testAccContainerNodePool_withGPU() string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1c" {
  location = "us-central1-c"
}

resource "google_container_cluster" "cluster" {
  name               = "tf-cluster-nodepool-test-%s"
  location           = "us-central1-c"
  initial_node_count = 1
  node_version       = data.google_container_engine_versions.central1c.latest_node_version
  min_master_version = data.google_container_engine_versions.central1c.latest_master_version
}

resource "google_container_node_pool" "np_with_gpu" {
  name     = "tf-nodepool-test-%s"
  location = "us-central1-c"
  cluster  = google_container_cluster.cluster.name

  initial_node_count = 1

  node_config {
    machine_type = "n1-standard-1"
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
    image_type      = "COS"

    guest_accelerator {
      type  = "nvidia-tesla-k80"
      count = 1
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerNodePool_withNodeConfigScopeAlias() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "tf-cluster-nodepool-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_container_node_pool" "np_with_node_config_scope_alias" {
  name               = "tf-nodepool-test-%s"
  location           = "us-central1-a"
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1
  node_config {
    machine_type = "g1-small"
    disk_size_gb = 10
    oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
  }
}
`, acctest.RandString(10), acctest.RandString(10))
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
