package google

import (
	"fmt"
	"log"
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
		"four clusters (0-indexed, ordered)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"0": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"1": testGenBigtableCluster("cluster2", "us-central1-a", "", "SSD"),
				"2": testGenBigtableCluster("cluster3", "us-central1-a", "3", ""),
				"3": testGenBigtableCluster("cluster4", "us-central1-a", "", ""),
			},
		},
		"four clusters (0-indexed, unordered indexes)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"2": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"1": testGenBigtableCluster("cluster2", "us-central1-a", "", "SSD"),
				"3": testGenBigtableCluster("cluster3", "us-central1-a", "3", ""),
				"0": testGenBigtableCluster("cluster4", "us-central1-a", "", ""),
			},
		},
		"four clusters (0-indexed, unordered elements)": {
			StateVersion: 0,
			ClusterSpecs: map[string]map[string]string{
				"3": testGenBigtableCluster("cluster4", "us-central1-a", "", ""),
				"0": testGenBigtableCluster("cluster1", "us-central1-a", "3", "SSD"),
				"2": testGenBigtableCluster("cluster3", "us-central1-a", "3", ""),
				"1": testGenBigtableCluster("cluster2", "us-central1-a", "", "SSD"),
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
			log.Printf("[DEBUG] Running test '%s' (withRemovedFields=%#v)", tn, withRemovedFields)
			state := testGenBigtableInstanceState(tc.ClusterSpecs, withRemovedFields)
			newState, err := resourceBigtableInstanceMigrateState(tc.StateVersion, state, nil)
			if err != nil {
				t.Fatalf("bad: %s, err: %#v", tn, err)
			}

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

// Run tests to validate that the state has been correctly migrated, and that the
// resource data it represents is also generated correctly from the state
func testBigtableInstanceMigrateStateCheckClusters(
	tn string,
	clusterSpecs map[string]map[string]string,
	newState *terraform.InstanceState,
	t *testing.T,
) {
	attributes := newState.Attributes
	data := resourceBigtableInstance().Data(newState)
	clusters := data.Get("cluster").([]interface{})

	numClustersExpected := len(clusterSpecs)

	// First, validate that the expected number of clusters exist in both the state
	// and the resource data it represents
	validateNumClusters(tn, numClustersExpected, attributes, clusters, t)

	// Ensure that each expected cluster exists only once in both the state and the
	// resource data it represents
	expectedAttributes := validateClustersExist(tn, clusterSpecs, numClustersExpected, attributes, clusters, t)

	// Make sure clusters stayed in order if the input state was 0-indexed
	validateClusterOrder(tn, clusterSpecs, numClustersExpected, attributes, clusters, t)

	// Make sure no unexpected attributes are in the state
	validateNoUnexpectedAttributes(tn, attributes, expectedAttributes, t)
}

func validateNumClusters(
	tn string,
	numClustersExpected int,
	attributes map[string]string,
	clusters []interface{},
	t *testing.T,
) {
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
}

func validateClustersExist(
	tn string,
	clusterSpecs map[string]map[string]string,
	numClustersExpected int,
	attributes map[string]string,
	clusters []interface{},
	t *testing.T,
) map[string]bool {
	// Keep track of the attributes we expect to be in the state so we can later make
	// sure nothing unexpected is also there
	expectedAttributes := make(map[string]bool)
	coreAttributes := []string{"cluster.#", "display_name", "id", "instance_type", "name", "project"}
	for _, attr := range coreAttributes {
		expectedAttributes[attr] = true
	}

	for _, clusterSpec := range clusterSpecs {
		// Look for cluster in migrated state
		numMatches, expectedAttributesForCluster := findMatchingClustersInState(clusterSpec, numClustersExpected, attributes)
		if numMatches == 1 {
			// Add found attributes to the overall set of expected ones
			for attr, _ := range expectedAttributesForCluster {
				expectedAttributes[attr] = true
			}
		} else {
			assertionKey := "matching clusters in state attributes"
			in := fmt.Sprintf("for\n - cluster: %#v\n - state attributes: %#v", clusterSpec, attributes)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %d\n got: %s -> %d\n%s",
				tn, assertionKey, 1, assertionKey, numMatches, in)
		}

		// Look for cluster in resource data
		numMatches = findMatchingClustersInResourceData(clusterSpec, clusters)
		if numMatches != 1 {
			assertionKey := "matching clusters in resource data"
			in := fmt.Sprintf("for\n - cluster: %#v\n - resource data clusters: %#v", clusterSpec, clusters)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %d\n got: %s -> %d\n%s",
				tn, assertionKey, 1, assertionKey, numMatches, in)
		}
	}

	return expectedAttributes
}

func findMatchingClustersInState(
	clusterSpec map[string]string,
	numClustersExpected int,
	attributes map[string]string,
) (int, map[string]bool) {
	numMatches := 0
	expectedAttributes := make(map[string]bool)
	for i := 0; i < numClustersExpected; i++ {
		var hits, misses int = 0, 0
		clusterAttributes := make([]string, 0, len(clusterSpec))
		for key, expectedValue := range clusterSpec {
			attrKey := fmt.Sprintf("cluster.%d.%s", i, key)
			if value, exists := attributes[attrKey]; exists {
				if value == expectedValue {
					clusterAttributes = append(clusterAttributes, attrKey)
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
			for _, attr := range clusterAttributes {
				expectedAttributes[attr] = true
			}
			numMatches++
		}
	}

	return numMatches, expectedAttributes
}

func findMatchingClustersInResourceData(clusterSpec map[string]string, clusters []interface{}) int {
	numMatches := 0
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

	return numMatches
}

func validateClusterOrder(
	tn string,
	clusterSpecs map[string]map[string]string,
	numClustersExpected int,
	attributes map[string]string,
	clusters []interface{},
	t *testing.T,
) {
	shouldRetainOrder := true
	for idxOrHash, _ := range clusterSpecs {
		idxOrHashInt, err := strconv.Atoi(idxOrHash)
		if err == nil && idxOrHashInt < numClustersExpected {
			shouldRetainOrder = shouldRetainOrder && true
		} else {
			shouldRetainOrder = false
		}
	}
	if shouldRetainOrder {
		for key, clusterSpec := range clusterSpecs {
			idx, _ := strconv.Atoi(key)
			clusterIdSpec := clusterSpec["cluster_id"]
			// Check state
			attrKey := fmt.Sprintf("cluster.%d.cluster_id", idx)
			clusterIdState := attributes[attrKey]
			if clusterIdState != clusterIdSpec {
				assertionKey := fmt.Sprintf("attributes[%s]", attrKey)
				in := fmt.Sprintf("for\n - state attributes: %#v", attributes)
				t.Fatalf(
					"bad: %s\n\n expected: %s -> %s\n got: %s -> %s\n%s",
					tn, assertionKey, clusterIdSpec, assertionKey, clusterIdState, in)
			}
			// Also check resource data
			cluster := clusters[idx].(map[string]interface{})
			clusterIdResourceData := cluster["cluster_id"]
			if clusterIdResourceData != clusterIdSpec {
				assertionKey := fmt.Sprintf("data.Get(\"cluster\")[%d][\"cluster_id\"]", idx)
				in := fmt.Sprintf("for\n - resource data clusters: %#v", clusters)
				t.Fatalf(
					"bad: %s\n\n expected: %s -> %s\n got: %s -> %s\n%s",
					tn, assertionKey, clusterIdSpec, assertionKey, clusterIdResourceData, in)
			}
		}
	}
}

func validateNoUnexpectedAttributes(
	tn string,
	attributes map[string]string,
	expectedAttributes map[string]bool,
	t *testing.T,
) {
	for k, v := range attributes {
		if _, exists := expectedAttributes[k]; !exists {
			assertionKey := fmt.Sprintf("attributes[\"%s\"]", k)
			in := fmt.Sprintf("for\n - state attributes: %#v", attributes)
			t.Fatalf(
				"bad: %s\n\n expected: %s -> %s\n got: %s -> %s\n%s",
				tn, assertionKey, "(should not be set)", assertionKey, v, in)
		}
	}
}
