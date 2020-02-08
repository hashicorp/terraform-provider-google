package google

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				Config:      testAccBigtableInstance_invalid(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block"),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableInstance(instanceName, 4),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
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
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableInstance_clusterModified(instanceName, 5),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
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
				Config:      testAccBigtableInstance_development_invalid_no_cluster(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: instance with instance_type=\"DEVELOPMENT\" should have exactly one \"cluster\" block"),
			},
			{
				Config:      testAccBigtableInstance_development_invalid_num_nodes(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: num_nodes cannot be set for instance_type=\"DEVELOPMENT\""),
			},
			{
				Config: testAccBigtableInstance_development(instanceName),
			},
			{
				ResourceName:      "google_bigtable_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccBigtableInstance_invalid(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
}
`, instanceName)
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
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
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
    num_nodes    = 3
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-d"
    zone         = "us-central1-f"
    num_nodes    = 3
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-e"
    zone         = "us-east1-a"
    num_nodes    = 3
    storage_type = "HDD"
  }
}
`, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableInstance_clusterReordered(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
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
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_clusterModified(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s-c"
    zone         = "us-central1-c"
    num_nodes    = %d
    storage_type = "HDD"
  }
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
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_development(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }
  instance_type = "DEVELOPMENT"
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_development_invalid_num_nodes(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
    num_nodes  = 3
  }
  instance_type = "DEVELOPMENT"
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_development_invalid_no_cluster(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
}
`, instanceName)
}
