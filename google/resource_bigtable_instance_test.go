package google

import (
	"context"
	"fmt"
	//"log"
	"regexp"
	"strconv"
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

// Check both state and resource
func testUnitBigtableInstance_checkClusters(clusterSpecs map[string]map[string]string, newState *terraform.InstanceState, t *testing.T) {
	attributes := newState.Attributes
	data := resourceBigtableInstance().Data(newState)
	//log.Printf("data: %v", data)
	clusters := data.Get("cluster").([]interface{})
	//log.Printf("clusters: %v", clusters)

	numClustersExpected := len(clusterSpecs)
	numClustersActualState, _ := strconv.Atoi(attributes["cluster.#"])
	numClustersActualResource := len(clusters)
	if numClustersActualState != numClustersExpected {
		t.Fatalf("Num clusters in migrated state (%d) incorrect; expected %d", numClustersActualState, numClustersExpected)
	}
	if numClustersActualResource != numClustersExpected {
		t.Fatalf("Num clusters in resource data (%d) incorrect; expected %d", numClustersActualResource, numClustersExpected)
	}

	clusterSpecByIndex := make(map[int]map[string]string)

	for _, clusterSpec := range clusterSpecs {
		// Look for cluster in migrated state
		numMatches := 0
		// TODO: Look for more clusters?
		for i := 0; i < numClustersExpected; i++ {
			var hits, misses int = 0, 0
			for key, expectedValue := range clusterSpec {
				attrKey := fmt.Sprintf("cluster.%d.%s", i, key)
				if value, exists := attributes[attrKey]; exists {
					if value == expectedValue {
						hits++
					} else {
						misses++
					}
				}
				if hits+misses == len(clusterSpec) {
					continue
				}
			}
			if hits == len(clusterSpec) {
				clusterSpecByIndex[i] = clusterSpec
				numMatches++
			}
		}
		if numMatches == 0 {
			t.Fatalf("Did not find cluster %#v in state attributes %#v", clusterSpec, attributes)
		} else if numMatches > 1 {
			t.Fatalf("Found multiple matches for cluster %#v in state attributes %#v", clusterSpec, attributes)
		}

		// Look for cluster in resource data
		numMatches = 0
		for _, cl := range clusters {
			cluster := cl.(map[string]interface{})
			hits := 0
			for key, expectedValue := range clusterSpec {
				if value, exists := cluster[key]; exists {
					if key == "num_nodes" {
						expected, _ := strconv.Atoi(expectedValue)
						if value == expected {
							hits++
						}
					} else {
						if value == expectedValue {
							hits++
						}
					}
				}
			}
			if hits == len(clusterSpec) {
				numMatches++
			}
		}
		if numMatches == 0 {
			t.Fatalf("Did not find cluster %#v in resource data %#v", clusterSpec, clusters)
		} else if numMatches > 1 {
			t.Fatalf("Found multiple matches for cluster %#v in resource data %#v", clusterSpec, clusters)
		}
	}

	// Make sure nothing exists that shouldn't
	// TODO: Also do for resource data?
	fields := []string{"cluster_id", "num_nodes", "storage_type", "zone"}
	for idx, clusterSpec := range clusterSpecByIndex {
		for _, field := range fields {
			_, existsSpec := clusterSpec[field]
			attrKey := fmt.Sprintf("cluster.%d.%s", idx, field)
			_, existsState := attributes[attrKey]
			if existsState && !existsSpec {
				t.Fatalf("Found %s unexpectedly for cluster_id=%s in state attributes %#v", field, clusterSpec["cluster_id"], attributes)
			}
		}
	}

	// Make sure fields marked as Removed do not exist in the state attributes
	for _, field := range fields {
		if v, exists := attributes[field]; exists {
			t.Fatalf("Found removed field %s (value: '%s') unexpectedly in state attributes %#v", field, v, attributes)
		}
	}
}

func testUnitGenCluster(clusterId string, zone string, numNodes string, storageType string) map[string]string {
	m := map[string]string{
		"cluster_id": clusterId,
		"zone":       zone,
	}
	if numNodes != "" {
		m["num_nodes"] = numNodes
	}
	if storageType != "" {
		m["storage_type"] = storageType
	}

	return m
}

func TestUnitBigtableInstance_MigrateState(t *testing.T) {
	t.Parallel()

	testCases := []map[string]map[string]string{
		{
			"0": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
		},
		{
			"0": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
			"1": testUnitGenCluster("cluster2", "us-central1-a", "3", "SSD"),
		},
		{
			"1": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
			"0": testUnitGenCluster("cluster2", "us-central1-a", "3", "SSD"),
		},
		{
			"1234567890": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
		},
		{
			"1234567890": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
			"9876543210": testUnitGenCluster("cluster2", "us-central1-a", "3", "SSD"),
		},
		{
			"9876543210": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
			"1234567890": testUnitGenCluster("cluster2", "us-central1-a", "3", "SSD"),
		},
		{
			"1234567890": testUnitGenCluster("cluster1", "us-central1-a", "3", "SSD"),
			"9876543210": testUnitGenCluster("cluster2", "us-central1-a", "", "SSD"),
			"6789054321": testUnitGenCluster("cluster3", "us-central1-a", "3", ""),
			"5432167890": testUnitGenCluster("cluster4", "us-central1-a", "", ""),
		},
	}

	// Test with and without top-level fields marked as "Removed"
	withRemovedFieldsOptions := []bool{true, false}

	for _, clusterSpecs := range testCases {
		for _, withRemovedFields := range withRemovedFieldsOptions {
			//log.Printf("clusterSpecs: %v", clusterSpecs)
			state := testUnitGenInstanceState(clusterSpecs, withRemovedFields)
			//log.Printf("state: %v", state)
			newState, err := resourceBigtableInstance().MigrateState(0, state, nil)
			if err != nil {
				t.Fatalf("MigrateState returned error: %s", err)
			}
			//log.Printf("newState: %v", newState)

			testUnitBigtableInstance_checkClusters(clusterSpecs, newState, t)
		}
	}
}

func testUnitGenInstanceState(clusterSpecs map[string]map[string]string, withRemovedFields bool) *terraform.InstanceState {
	attributes := map[string]string{
		"display_name":  "foo",
		"instance_type": "PRODUCTION",
		"name":          "foo",
		"project":       "some-project",
	}

	attributes["cluster.#"] = strconv.Itoa(len(clusterSpecs))

	for idxOrHash, clusterSpec := range clusterSpecs {
		for k, v := range clusterSpec {
			attribute := fmt.Sprintf("cluster.%s.%s", idxOrHash, k)
			attributes[attribute] = v
		}
	}

	if withRemovedFields {
		attributes["cluster_id"] = ""
		attributes["num_nodes"] = "0"
		attributes["storage_type"] = "SSD"
		attributes["zone"] = ""
	}

	state := &terraform.InstanceState{
		ID:         "foo",
		Attributes: attributes,
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
