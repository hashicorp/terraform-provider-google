package google

import (
	"fmt"
	//"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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

// Check both state and resource
func testUnitBigtableInstance_checkClusters(
	clusterSpecs map[string]map[string]string,
	newState *terraform.InstanceState,
	t *testing.T,
) {
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
