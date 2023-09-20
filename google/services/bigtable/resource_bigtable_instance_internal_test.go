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
)

func TestGetUnavailableClusterZones(t *testing.T) {
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

func TestGetInstanceFromResponse(t *testing.T) {
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
		"unavailble error": {
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
