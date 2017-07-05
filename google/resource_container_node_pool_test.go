package google

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
		},
	})
}

func TestAccContainerNodePool_autoscaling(t *testing.T) {
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
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
			resource.TestStep{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
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
		_, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
			config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("NodePool still exists")
		}
	}

	return nil
}

func testAccCheckContainerNodePoolMatches(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attributes := rs.Primary.Attributes
		found, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
			config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != attributes["name"] {
			return fmt.Errorf("NodePool not found")
		}

		inc, err := strconv.Atoi(attributes["initial_node_count"])
		if err != nil {
			return err
		}
		if found.InitialNodeCount != int64(inc) {
			return fmt.Errorf("Mismatched initialNodeCount. TF State: %s. GCP State: %d",
				attributes["initial_node_count"], found.InitialNodeCount)
		}

		tfAS := attributes["autoscaling.#"] == "1"
		if gcpAS := found.Autoscaling != nil && found.Autoscaling.Enabled == true; tfAS != gcpAS {
			return fmt.Errorf("Mismatched autoscaling status. TF State: %t. GCP State: %t", tfAS, gcpAS)
		}
		if tfAS {
			if tf := attributes["autoscaling.0.min_node_count"]; strconv.FormatInt(found.Autoscaling.MinNodeCount, 10) != tf {
				return fmt.Errorf("Mismatched Autoscaling.MinNodeCount. TF State: %s. GCP State: %d",
					tf, found.Autoscaling.MinNodeCount)
			}

			if tf := attributes["autoscaling.0.max_node_count"]; strconv.FormatInt(found.Autoscaling.MaxNodeCount, 10) != tf {
				return fmt.Errorf("Mismatched Autoscaling.MaxNodeCount. TF State: %s. GCP State: %d",
					tf, found.Autoscaling.MaxNodeCount)
			}

		}

		return nil
	}
}

func testAccContainerNodePool_basic(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
}`, cluster, np)
}

func testAccContainerNodePool_autoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
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

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 1
		max_node_count = 5
	}
}`, cluster, np)
}
