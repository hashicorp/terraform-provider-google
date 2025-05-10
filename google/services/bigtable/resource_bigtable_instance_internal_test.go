// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitBigtable_getUnavailableClusterZones(t *testing.T) {
	cases := map[string]struct {
		clusterZones     []string
		unavailableZones []string
		want             []string
	}{
		"not unavailalbe": {
			clusterZones:     []string{"us-central1", "eu-west1"},
			unavailableZones: []string{"us-central2", "eu-west2"},
			want:             nil,
		},
		"unavailable one to one": {
			clusterZones:     []string{"us-central2"},
			unavailableZones: []string{"us-central2"},
			want:             []string{"us-central2"},
		},
		"unavailable one to many": {
			clusterZones:     []string{"us-central2"},
			unavailableZones: []string{"us-central2", "us-central1"},
			want:             []string{"us-central2"},
		},
		"unavailable many to one": {
			clusterZones:     []string{"us-central2", "us-central1"},
			unavailableZones: []string{"us-central2"},
			want:             []string{"us-central2"},
		},
		"unavailable many to many": {
			clusterZones:     []string{"us-central2", "us-central1"},
			unavailableZones: []string{"us-central2", "us-central1"},
			want:             []string{"us-central2", "us-central1"},
		},
	}

	for tn, tc := range cases {
		var clusters []interface{}
		for _, zone := range tc.clusterZones {
			clusters = append(clusters, map[string]interface{}{"zone": zone})
		}
		if got := getUnavailableClusterZones(clusters, tc.unavailableZones); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("bad: %s, got %q, want %q", tn, got, tc.want)
		}
	}
}

func TestUnitBigtable_getInstanceFromResponse(t *testing.T) {
	instanceName := "test-instance"
	originalId := "original_value"
	cases := map[string]struct {
		instanceNames      []string
		listInstancesError error

		wantError        string
		wantInstanceName string
		wantStop         bool
		wantId           string
	}{
		"not found": {
			instanceNames:      []string{"wrong", "also_wrong"},
			listInstancesError: nil,

			wantError:        "",
			wantStop:         true,
			wantInstanceName: "",
			wantId:           "",
		},
		"found": {
			instanceNames:      []string{"wrong", "also_wrong", instanceName},
			listInstancesError: nil,

			wantError:        "",
			wantStop:         false,
			wantInstanceName: instanceName,
			wantId:           originalId,
		},
		"error": {
			instanceNames:      nil,
			listInstancesError: fmt.Errorf("some error"),

			wantError:        "Error retrieving instance.",
			wantStop:         true,
			wantInstanceName: "",
			wantId:           originalId,
		},
		"unavailable error": {
			instanceNames:      []string{"wrong", "also_wrong"},
			listInstancesError: bigtable.ErrPartiallyUnavailable{[]string{"some", "location"}},

			wantError:        "",
			wantStop:         false,
			wantInstanceName: "",
			wantId:           originalId,
		}}
	for tn, tc := range cases {
		instancesResponse := []*bigtable.InstanceInfo{}
		for _, existingInstance := range tc.instanceNames {
			instancesResponse = append(instancesResponse, &bigtable.InstanceInfo{Name: existingInstance})
		}
		d := &schema.ResourceData{}
		d.SetId(originalId)
		gotInstance, gotStop, gotErr := getInstanceFromResponse(instancesResponse, instanceName, tc.listInstancesError, d)

		if gotStop != tc.wantStop {
			t.Errorf("bad stop: %s, got %v, want %v", tn, gotStop, tc.wantStop)
		}
		if (gotErr != nil && tc.wantError == "") ||
			(gotErr == nil && tc.wantError != "") ||
			(gotErr != nil && !strings.Contains(gotErr.Error(), tc.wantError)) {
			t.Errorf("bad error: %s, got %q, want %q", tn, gotErr, tc.wantError)
		}
		if (gotInstance == nil && tc.wantInstanceName != "") ||
			(gotInstance != nil && tc.wantInstanceName == "") ||
			(gotInstance != nil && gotInstance.Name != tc.wantInstanceName) {
			t.Errorf("bad instance: %s, got %v, want %q", tn, gotInstance, tc.wantInstanceName)
		}
		gotId := d.Id()
		if gotId != tc.wantId {
			t.Errorf("bad ID: %s, got %v, want %q", tn, gotId, tc.wantId)
		}
	}
}

func TestUnitBigtable_flattenBigtableCluster(t *testing.T) {
	cases := map[string]struct {
		clusterInfo *bigtable.ClusterInfo
		want        map[string]interface{}
	}{
		"SSD auto scaling": {
			clusterInfo: &bigtable.ClusterInfo{
				StorageType: bigtable.SSD,
				Zone:        "zone1",
				ServeNodes:  5,
				Name:        "ssd-cluster",
				KMSKeyName:  "KMS",
				State:       "CREATING",
				AutoscalingConfig: &bigtable.AutoscalingConfig{
					MinNodes:                  3,
					MaxNodes:                  7,
					CPUTargetPercent:          50,
					StorageUtilizationPerNode: 60,
				},
			},
			want: map[string]interface{}{
				"zone":         "zone1",
				"num_nodes":    5,
				"cluster_id":   "ssd-cluster",
				"storage_type": "SSD",
				"kms_key_name": "KMS",
				"state":        "CREATING",
				"autoscaling_config": []map[string]interface{}{
					{
						"min_nodes":      3,
						"max_nodes":      7,
						"cpu_target":     50,
						"storage_target": 60,
					},
				},
				// unspecified node scaling factor in input will default to 1X
				"node_scaling_factor": "NodeScalingFactor1X",
			},
		},
		"HDD manual scaling": {
			clusterInfo: &bigtable.ClusterInfo{
				StorageType:       bigtable.HDD,
				Zone:              "zone2",
				ServeNodes:        7,
				Name:              "hdd-cluster",
				KMSKeyName:        "KMS",
				State:             "READY",
				NodeScalingFactor: bigtable.NodeScalingFactor2X,
			},
			want: map[string]interface{}{
				"zone":                "zone2",
				"num_nodes":           7,
				"cluster_id":          "hdd-cluster",
				"storage_type":        "HDD",
				"kms_key_name":        "KMS",
				"state":               "READY",
				"node_scaling_factor": "NodeScalingFactor2X",
			},
		},
	}

	for tn, tc := range cases {
		if got := flattenBigtableCluster(tc.clusterInfo); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("bad: %s, got %q, want %q", tn, got, tc.want)
		}
	}
}

func TestUnitBigtable_resourceBigtableInstanceClusterReorderTypeListFunc_error(t *testing.T) {
	d := &tpgresource.ResourceDiffMock{
		After: map[string]interface{}{
			"cluster.#": 0,
		},
	}
	if err := resourceBigtableInstanceClusterReorderTypeListFunc(d, nil); err == nil {
		t.Errorf("expected error, got success")
	}
}

func TestUnitBigtable_resourceBigtableInstanceClusterReorderTypeListFunc(t *testing.T) {
	cases := map[string]struct {
		before           map[string]interface{}
		after            map[string]interface{}
		wantClusterOrder []string
		wantForceNew     bool
	}{
		"create": {
			before: map[string]interface{}{
				"cluster.#":            1,
				"cluster.0.cluster_id": "some-id-a",
			},
			after: map[string]interface{}{
				"name":                 "some-name",
				"cluster.#":            1,
				"cluster.0.cluster_id": "some-id-a",
				"cluster.0": map[string]interface{}{
					"cluster_id": "some-id-a",
				},
			},
			wantClusterOrder: []string{},
			wantForceNew:     false,
		},
		"no force new change": {
			before: map[string]interface{}{
				"name":                 "some-name",
				"cluster.#":            4,
				"cluster.0.cluster_id": "some-id-a",
				"cluster.1.cluster_id": "some-id-b",
				"cluster.2.cluster_id": "some-id-c",
				"cluster.3.cluster_id": "some-id-e",
			},
			after: map[string]interface{}{
				"name":                 "some-name",
				"cluster.#":            3,
				"cluster.0.cluster_id": "some-id-c",
				"cluster.1.cluster_id": "some-id-a",
				"cluster.2.cluster_id": "some-id-d",
				"cluster.0": map[string]interface{}{
					"cluster_id": "some-id-c",
				},
				"cluster.1": map[string]interface{}{
					"cluster_id": "some-id-a",
				},
				"cluster.2": map[string]interface{}{
					"cluster_id": "some-id-d",
				},
			},
			wantClusterOrder: []string{"some-id-a", "some-id-d", "some-id-c"},
			wantForceNew:     false,
		},
		"force new - zone change": {
			before: map[string]interface{}{
				"name":                 "some-name",
				"cluster.#":            1,
				"cluster.0.cluster_id": "some-id-a",
				"cluster.0.zone":       "zone-a",
			},
			after: map[string]interface{}{
				"name":                 "some-name",
				"cluster.#":            1,
				"cluster.0.cluster_id": "some-id-a",
				"cluster.0.zone":       "zone-b",
				"cluster.0": map[string]interface{}{
					"cluster_id": "some-id-a",
					"zone":       "zone-b",
				},
			},
			wantClusterOrder: []string{"some-id-a"},
			wantForceNew:     true,
		},
		"force new - kms_key_name change": {
			before: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.kms_key_name": "key-a",
			},
			after: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.kms_key_name": "key-b",
				"cluster.0": map[string]interface{}{
					"cluster_id":   "some-id-a",
					"kms_key_name": "key-b",
				},
			},
			wantClusterOrder: []string{"some-id-a"},
			wantForceNew:     true,
		},
		"force new - storage_type change": {
			before: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.storage_type": "HDD",
				"cluster.0.state":        "READY",
			},
			after: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.storage_type": "SSD",
				"cluster.0": map[string]interface{}{
					"cluster_id":   "some-id-a",
					"storage_type": "SSD",
				},
			},
			wantClusterOrder: []string{"some-id-a"},
			wantForceNew:     true,
		},
		"skip force new - storage_type change for CREATING cluster": {
			before: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.storage_type": "SSD",
				"cluster.0.state":        "CREATING",
			},
			after: map[string]interface{}{
				"name":                   "some-name",
				"cluster.#":              1,
				"cluster.0.cluster_id":   "some-id-a",
				"cluster.0.storage_type": "HDD",
				"cluster.0": map[string]interface{}{
					"cluster_id":   "some-id-a",
					"storage_type": "HDD",
				},
			},
			wantClusterOrder: []string{"some-id-a"},
			wantForceNew:     false,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			d := &tpgresource.ResourceDiffMock{
				Before: tc.before,
				After:  tc.after,
			}
			var clusters []interface{}
			err := resourceBigtableInstanceClusterReorderTypeListFunc(d, func(gotClusters []interface{}) error {
				clusters = gotClusters
				return nil
			})
			if err != nil {
				t.Fatalf("bad: %s, error: %v", tn, err)
			}
			if d.IsForceNew != tc.wantForceNew {
				t.Errorf("bad: %s, got %v, want %v", tn, d.IsForceNew, tc.wantForceNew)
			}
			gotClusterOrder := []string{}
			for _, cluster := range clusters {
				clusterResource := cluster.(map[string]interface{})
				gotClusterOrder = append(gotClusterOrder, clusterResource["cluster_id"].(string))
			}
			if !reflect.DeepEqual(gotClusterOrder, tc.wantClusterOrder) {
				t.Errorf("bad: %s, got %q, want %q", tn, gotClusterOrder, tc.wantClusterOrder)
			}
		})
	}
}
