package google

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBigtableInstance_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableInstance_invalid(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block"),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance(instanceName, 4),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_cluster(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableInstance_clusterMax(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: Too many cluster blocks: No more than 4 \"cluster\" blocks are allowed"),
			},
			{
				Config: testAccBigtableInstance_cluster(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterModified(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_development(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_development(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_allowDestroy(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_noAllowDestroy(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config:      testAccBigtableInstance_noAllowDestroy(instanceName, 3),
				Destroy:     true,
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
		},
	})
}

func testAccCheckBigtableInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_instance" {
				continue
			}

			config := googleProviderConfig(t)
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

  deletion_protection = false
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

  deletion_protection = false
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

  deletion_protection = false
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

  deletion_protection = false
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

  deletion_protection = false
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

  deletion_protection = false
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_noAllowDestroy(instanceName string, numNodes int) string {
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
