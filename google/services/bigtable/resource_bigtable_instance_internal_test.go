// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"reflect"
	"testing"
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
