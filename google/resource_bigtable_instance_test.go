package google

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Config: testAccBigtableInstance_cluster_reordered(instanceName, 5),
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

func TestUnitBigtableInstance_regression(t *testing.T) {
	t.Parallel()

	clusterIdxPairs := [][]string{{"0", "1"}, {"1234567890", "9876543210"}}

	for _, clusterIdxs := range clusterIdxPairs {
		clusterIdx1 := clusterIdxs[0]
		clusterIdx2 := clusterIdxs[1]

		state := testUnitGenInstanceState(clusterIdx1, clusterIdx2)

		data := resourceBigtableInstance().Data(state)

		log.Printf("%v", data.Get("cluster"))
		log.Printf("%T", data.Get("cluster"))
		for k, v := range data.Get("cluster").([]interface{}) {
			log.Printf("%v", k)
			log.Printf("%v", v)
		}

		if data.Id() != "foo" {
			t.Fatalf("ID incorrect: %s", data.Id())
		}

		// For v2.13.0
		//clusters := data.Get("cluster").(*schema.Set).List()
		// For v2.14.0+
		clusters := data.Get("cluster").([]interface{})

		testUnitBigtableInstance_checkClusters(clusters, t)
	}
}

func testUnitBigtableInstance_checkClusters(clusters []interface{}, t *testing.T) {
	numClusters := len(clusters)
	if numClusters != 2 {
		t.Fatalf("Num clusters incorrect: %d", numClusters)
	}
	cluster1 := clusters[0].(map[string]interface{})
	clusterId1 := cluster1["cluster_id"]
	if clusterId1 != "cluster1" {
		t.Fatalf("cluster_id (1) incorrect: %s", clusterId1)
	}
	cluster2 := clusters[1].(map[string]interface{})
	clusterId2 := cluster2["cluster_id"]
	if clusterId2 != "cluster2" {
		t.Fatalf("cluster_id (2) incorrect: %s", clusterId2)
	}
}

func TestUnitBigtableInstance_MigrateState(t *testing.T) {
	t.Parallel()

	clusterIdxPairs := [][]string{{"0", "1"}, {"1234567890", "9876543210"}}

	for _, clusterIdxs := range clusterIdxPairs {
		clusterIdx1 := clusterIdxs[0]
		clusterIdx2 := clusterIdxs[1]

		state := testUnitGenInstanceState(clusterIdx1, clusterIdx2)
		log.Printf("state: %v", state)
		newState, err := resourceBigtableInstance().MigrateState(0, state, nil)
		if err != nil {
			t.Fatalf("MigrateState returned error: %s", err)
		}
		log.Printf("newState: %v", newState)
		data := resourceBigtableInstance().Data(newState)
		log.Printf("data: %v", data)
		clusters := data.Get("cluster").([]interface{})

		testUnitBigtableInstance_checkClusters(clusters, t)
	}
}

func testUnitGenInstanceState(clusterIdx1 string, clusterIdx2 string) *terraform.InstanceState {
	clusterPrefix1 := fmt.Sprintf("cluster.%s", clusterIdx1)
	clusterPrefix2 := fmt.Sprintf("cluster.%s", clusterIdx2)

	state := &terraform.InstanceState{
		ID: "foo",
		Attributes: map[string]string{
			"cluster.#": "2",
			fmt.Sprintf("%s.cluster_id", clusterPrefix1):   "cluster1",
			fmt.Sprintf("%s.num_nodes", clusterPrefix1):    "3",
			fmt.Sprintf("%s.storage_type", clusterPrefix1): "SSD",
			fmt.Sprintf("%s.zone", clusterPrefix1):         "us-central1-a",
			fmt.Sprintf("%s.cluster_id", clusterPrefix2):   "cluster2",
			fmt.Sprintf("%s.num_nodes", clusterPrefix2):    "3",
			fmt.Sprintf("%s.storage_type", clusterPrefix2): "SSD",
			fmt.Sprintf("%s.zone", clusterPrefix2):         "us-central1-a",
			"cluster_id":                                   "",
			"display_name":                                 "foo",
			"instance_type":                                "PRODUCTION",
			"name":                                         "foo",
			"num_nodes":                                    "0",
			"project":                                      "some-project",
			"storage_type":                                 "SSD",
			"zone":                                         "",
		},
	}

	return state
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

func testAccBigtableInstance_cluster_reordered(instanceName string, numNodes int) string {
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

func testAccBigtableInstance_development_invalid_num_nodes(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id    = "%s"
		zone          = "us-central1-b"
        num_nodes     = 3
	}
	instance_type = "DEVELOPMENT"
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_development_invalid_no_cluster(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	instance_type = "DEVELOPMENT"
}
`, instanceName)
}
