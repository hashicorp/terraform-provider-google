package google

import (
	"fmt"
	//"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestBigtableInstanceMigrateState(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		StateVersion int
		ClusterSpecs map[string]map[string]string
	}{
		"one cluster (0-indexed)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"0": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
			},
		},
		"two clusters (0-indexed, ordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"0": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"1": testGenBigtableCluster("cluster2", "us-central1-a", "3", "SSD"),
			},
		},
		"two clusters (0-indexed, unordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"1": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"0": testGenBigtableCluster("cluster2", "us-central1-a", "3", "SSD"),
			},
		},
		"one cluster (hash-indexed)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"1234567890": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
			},
		},
		"two clusters (hash-indexed, ordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"1234567890": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"9876543210": testGenBigtableCluster("cluster2", "us-central1-a", "3", "SSD"),
			},
		},
		"two clusters (hash-indexed, unordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"9876543210": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"1234567890": testGenBigtableCluster("cluster2", "us-central1-a", "3", "SSD"),
			},
		},
		"four clusters (hash-indexed, unordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"1234567890": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"9876543210": testGenBigtableCluster("cluster2", "us-central1-a", "", "SSD"),
				"6789054321": testGenBigtableCluster("cluster3", "us-central1-a", "3", ""),
				"5432167890": testGenBigtableCluster("cluster4", "us-central1-a", "", ""),
			},
		},
	}

	// Test with and without top-level fields marked as "Removed"
	withRemovedFieldsOptions := []bool{true, false}

	for tn, tc := range cases {
		for _, withRemovedFields := range withRemovedFieldsOptions {
			//log.Printf("clusterSpecs: %v", clusterSpecs)
			state := testGenBigtableInstanceState(tc.ClusterSpecs, withRemovedFields)
			//log.Printf("state: %v", state)
			newState, err := resourceBigtableInstanceMigrateState(tc.StateVersion, state, nil)
			if err != nil {
				t.Fatalf("bad: %s, err: %#v", tn, err)
			}
			//log.Printf("newState: %v", newState)

			testBigtableInstanceMigrateStateCheckClusters(tn, tc.ClusterSpecs, newState, t)
		}
	}
}

func testGenBigtableCluster(clusterId string, zone string, numNodes string, storageType string) map[string]string {
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

func testGenBigtableInstanceState(clusterSpecs map[string]map[string]string, withRemovedFields bool) *terraform.InstanceState {
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
func testBigtableInstanceMigrateStateCheckClusters(
	tn string,
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
		assertionKey := "num clusters in migrated state"
		t.Fatalf(
			"bad: %s\n\n expected: %s -> %d\n got: %s -> %d",
			tn, assertionKey, numClustersExpected, assertionKey, numClustersActualState)
	}
	if numClustersActualResource != numClustersExpected {
		assertionKey := "num clusters in resource data"
		t.Fatalf(
			"bad: %s\n\n expected: %s -> %d\n got: %s -> %d",
			tn, assertionKey, numClustersExpected, assertionKey, numClustersActualResource)
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
		if numMatches != 1 {
			assertionKey := "matching clusters in state attributes"
			in := fmt.Sprintf("for\n - cluster: %#v\n - state attributes: %#v", clusterSpec, attributes)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %d\n got: %s -> %d\n%s",
				tn, assertionKey, 1, assertionKey, numMatches, in)
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
		if numMatches != 1 {
			assertionKey := "matching clusters in resource data"
			in := fmt.Sprintf("for\n - cluster: %#v\n - resouce data clusters: %#v", clusterSpec, clusters)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %d\n got: %s -> %d\n%s",
				tn, assertionKey, 1, assertionKey, numMatches, in)
		}
	}

	// Make sure nothing exists that shouldn't
	fields := []string{"cluster_id", "num_nodes", "storage_type", "zone"}
	for idx, clusterSpec := range clusterSpecByIndex {
		for _, field := range fields {
			_, existsSpec := clusterSpec[field]
			attrKey := fmt.Sprintf("cluster.%d.%s", idx, field)
			stateValue, existsState := attributes[attrKey]
			if existsState && !existsSpec {
				assertionKey := fmt.Sprintf("cluster[%s]", field)
				in := fmt.Sprintf("for\n - cluster: %#v\n - state attributes: %#v", clusterSpec, attributes)
				t.Fatalf(
					"bad: %s\n\n expected: %s -> %s\n got: %s -> %s\n%s",
					tn, assertionKey, "(should not be set)", assertionKey, stateValue, in)
			}
		}
	}

	// Make sure fields marked as Removed do not exist in the state attributes
	for _, field := range fields {
		if v, exists := attributes[field]; exists {
			assertionKey := fmt.Sprintf("attributes[%s]", field)
			in := fmt.Sprintf("for\n - state attributes: %#v", attributes)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %s\n got: %s -> %s\n%s",
				tn, assertionKey, "(should not be set; removed from schema)", assertionKey, v, in)
		}
	}
}
