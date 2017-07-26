package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/container/v1"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	name := "tf-nodepool-test-" + acctest.RandString(10)
	zone := "us-central1-a"
	clusterConfig := SomeGoogleContainerCluster()
	nodepoolConfig := SomeGoogleContainerNodePool(clusterConfig).
		WithAttribute("name", name).
		WithAttribute("zone", zone).
		WithAttribute("initial_node_count", 2)

	var nodePool container.NodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig.String() + nodepoolConfig.String(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolExists(zone, clusterConfig.Name(), name, &nodePool),
					testAccCheckContainerNodePoolHasInitialNodeCount(&nodePool, 2),
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

func SomeGoogleContainerCluster() *ConfigBuilder {
	return NewResourceConfigBuilder("google_container_cluster", "cluster-"+acctest.RandString(10)).
		WithAttribute("name", "tf-cluster-nodepool-test-"+acctest.RandString(10)).
		WithAttribute("zone", "us-central1-a").
		WithAttribute("initial_node_count", 3).
		WithAttribute("master_auth", NewNestedConfig().
			WithAttribute("username", "mr.yoda").
			WithAttribute("password", "adoy.rm"))
}

func SomeGoogleContainerNodePool(cluster *ConfigBuilder) *ConfigBuilder {
	return NewResourceConfigBuilder("google_container_node_pool", "nodepool-"+acctest.RandString(10)).
		WithAttribute("name", "tf-nodepool-test-"+acctest.RandString(10)).
		WithAttribute("zone", "us-central1-a").
		WithAttribute("cluster", fmt.Sprintf("${google_container_cluster.%s.name}", cluster.ResourceName)).
		WithAttribute("initial_node_count", 2)
}

func testAccCheckContainerNodePoolExists(zone, clusterName, nodePoolName string, nodePool *container.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(config.Project, zone, clusterName, nodePoolName).Do()
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Unable to find resource")
		}

		*nodePool = *found
		return nil
	}
}

func testAccCheckContainerNodePoolHasInitialNodeCount(nodePool *container.NodePool, count int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nodePool.InitialNodeCount != count {
			return fmt.Errorf("Expected initial_node_count %d but found %d", count, nodePool.InitialNodeCount)
		}
		return nil
	}
}
