package google

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigtableInstance_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance(instanceName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance", 3),
				),
			},
			{
				Config: testAccBigtableInstance(instanceName, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance", 4),
				),
			},
		},
	})
}

func TestAccBigtableInstance_cluster(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableInstance_clusterMax(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: Too many cluster blocks: No more than 4 \"cluster\" blocks are allowed"),
			},
			{
				Config: testAccBigtableInstance_cluster(instanceName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance", 3),
				),
			},
			{
				Config: testAccBigtableInstance_cluster_reordered(instanceName, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance", 5),
				),
			},
		},
	})
}

func TestAccBigtableInstance_development(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_development(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance", 0),
				),
			},
		},
	})
}

func testAccCheckBigtableInstanceDestroy(s *terraform.State) error {
	var ctx = context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigtable_instance" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewInstanceAdminClient(config.Project)
		if err != nil {
			return fmt.Errorf("Error starting instance admin client. %s", err)
		}

		defer c.Close()

		_, err = c.InstanceInfo(ctx, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Instance %s still exists.", rs.Primary.Attributes["name"])
		}
	}

	return nil
}

func testAccBigtableInstanceExists(n string, numNodes int) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewInstanceAdminClient(config.Project)
		if err != nil {
			return fmt.Errorf("Error starting instance admin client. %s", err)
		}

		defer c.Close()

		_, err = c.InstanceInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving instance %s.", rs.Primary.Attributes["name"])
		}

		clusters, err := c.Clusters(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving cluster list for instance %s.", rs.Primary.Attributes["name"])
		}

		for _, c := range clusters {
			if c.ServeNodes != numNodes {
				return fmt.Errorf("Expected cluster %s to have %d nodes but got %d nodes for instance %s.",
					c.Name,
					numNodes,
					c.ServeNodes,
					rs.Primary.Attributes["name"])
			}
		}

		return nil
	}
}

func testAccBigtableInstance(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s"
		zone         = "us-central1-b"
		num_nodes    = %d
		storage_type = "HDD"
	}
}
`, instanceName, instanceName, numNodes)
}

func testAccBigtableInstance_cluster(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s-a"
		zone         = "us-central1-a"
		num_nodes    = %d
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-b"
		zone         = "us-central1-b"
		num_nodes    = %d
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-c"
		zone         = "us-central1-c"
		num_nodes    = %d
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-d"
		zone         = "us-central1-f"
		num_nodes    = %d
		storage_type = "HDD"
	}
}
`, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_clusterMax(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s-a"
		zone         = "us-central1-a"
		num_nodes    = 3
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-b"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-c"
		zone         = "us-central1-c"
		num_nodes    = %d
		storage_type = "HDD"
	}
}
`, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableInstance_cluster_reordered(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		num_nodes    = %d
		cluster_id   = "%s-b"
		zone         = "us-central1-c"
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-d"
		zone         = "us-central1-f"
		num_nodes    = %d
		storage_type = "HDD"
	}
	cluster {
		zone         = "us-central1-b"
		cluster_id   = "%s-a"
		num_nodes    = %d
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-e"
		zone         = "us-east1-a"
		num_nodes    = %d
		storage_type = "HDD"
	}
}
`, instanceName, numNodes, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName)
}

func testAccBigtableInstance_development(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id    = "%s"
		zone          = "us-central1-b"
	}
	instance_type = "DEVELOPMENT"
}
`, instanceName, instanceName)
}
