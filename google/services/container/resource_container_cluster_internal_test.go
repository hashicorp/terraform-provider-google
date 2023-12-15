// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestContainerClusterEnableK8sBetaApisCustomizeDiff(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before           *schema.Set
		after            *schema.Set
		expectedForceNew bool
	}{
		"no need to force new from nil to empty apis": {
			before:           schema.NewSet(schema.HashString, nil),
			after:            schema.NewSet(schema.HashString, []interface{}{}),
			expectedForceNew: false,
		},
		"no need to force new from empty apis to nil": {
			before:           schema.NewSet(schema.HashString, []interface{}{}),
			after:            schema.NewSet(schema.HashString, nil),
			expectedForceNew: false,
		},
		"no need to force new from empty apis to empty apis": {
			before:           schema.NewSet(schema.HashString, []interface{}{}),
			after:            schema.NewSet(schema.HashString, []interface{}{}),
			expectedForceNew: false,
		},
		"no need to force new from nil to empty string apis": {
			before:           schema.NewSet(schema.HashString, nil),
			after:            schema.NewSet(schema.HashString, []interface{}{""}),
			expectedForceNew: false,
		},
		"no need to force new from empty string apis to empty string apis": {
			before:           schema.NewSet(schema.HashString, []interface{}{""}),
			after:            schema.NewSet(schema.HashString, []interface{}{""}),
			expectedForceNew: false,
		},
		"no need to force new for enabling new api from empty apis": {
			before:           schema.NewSet(schema.HashString, []interface{}{}),
			after:            schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			expectedForceNew: false,
		},
		"no need to force new for enabling new api from nil": {
			before:           schema.NewSet(schema.HashString, nil),
			after:            schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			expectedForceNew: false,
		},
		"no need to force new for passing same apis": {
			before:           schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			after:            schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			expectedForceNew: false,
		},
		"no need to force new for passing same apis with inconsistent order": {
			before:           schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo", "dummy.k8s.io/v1beta1/bar"}),
			after:            schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/bar", "dummy.k8s.io/v1beta1/foo"}),
			expectedForceNew: false,
		},
		"need to force new from empty string apis to nil": {
			before:           schema.NewSet(schema.HashString, []interface{}{""}),
			after:            schema.NewSet(schema.HashString, nil),
			expectedForceNew: true,
		},
		"need to force new for disabling existing api": {
			before:           schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			after:            schema.NewSet(schema.HashString, []interface{}{}),
			expectedForceNew: true,
		},
		"need to force new for disabling existing api with nil": {
			before:           schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			after:            schema.NewSet(schema.HashString, nil),
			expectedForceNew: true,
		},
		"need to force new for disabling existing apis": {
			before:           schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo", "dummy.k8s.io/v1beta1/bar", "dummy.k8s.io/v1beta1/baz"}),
			after:            schema.NewSet(schema.HashString, []interface{}{"dummy.k8s.io/v1beta1/foo"}),
			expectedForceNew: true,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"enable_k8s_beta_apis.0.enabled_apis": tc.before,
			},
			After: map[string]interface{}{
				"enable_k8s_beta_apis.0.enabled_apis": tc.after,
			},
		}
		err := containerClusterEnableK8sBetaApisCustomizeDiffFunc(d)
		if err != nil {
			t.Errorf("%s failed, found unexpected error: %s", tn, err)
		}
		if d.IsForceNew != tc.expectedForceNew {
			t.Errorf("%v: expected d.IsForceNew to be %v, but was %v", tn, tc.expectedForceNew, d.IsForceNew)
		}
	}
}

func TestContainerCluster_NodeVersionCustomizeDiff(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		BeforeName    string
		AfterName     string
		MasterVersion string
		NodeVersion   string
		ExpectError   bool
	}{
		"Master version and node version are exactly the same": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.9-gke.5",
			NodeVersion:   "1.10.9-gke.5",
			ExpectError:   false,
		},
		"Master version and node version have the same Kubernetes patch version but not the same gke-N suffix ": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.9-gke.5",
			NodeVersion:   "1.10.9-gke.9",
			ExpectError:   false,
		},
		"Master version and node version have different minor versions": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.9-gke.5",
			NodeVersion:   "1.11.6-gke.11",
			ExpectError:   true,
		},
		"Master version and node version have different Kubernetes Patch Versions": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.9-gke.5",
			NodeVersion:   "1.10.6-gke.11",
			ExpectError:   true,
		},
		"Master version is not set, but node version is": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "",
			NodeVersion:   "1.10.6-gke.11",
			ExpectError:   false,
		},
		"Node version is not set, but master version is": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.6-gke.11",
			NodeVersion:   "",
			ExpectError:   false,
		},
		"Node version and master version match, both do not have -gke.X suffix": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.6",
			NodeVersion:   "1.10.6",
			ExpectError:   false,
		},
		"Node version and master version do not match, both do not have -gke.X suffix": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.10.6",
			NodeVersion:   "1.11.6",
			ExpectError:   true,
		},
		"Node version and master version do not match, node version has -gke.X suffix but master version doesn't": {
			BeforeName:    "",
			AfterName:     "test",
			MasterVersion: "1.11.6",
			NodeVersion:   "1.10.6-gke.11",
			ExpectError:   true,
		},
		"Diff is executed in non-create scenario, master version and node version do not match": {
			BeforeName:    "test",
			AfterName:     "test-1",
			MasterVersion: "1.11.6-gke.11",
			NodeVersion:   "1.10.6-gke.11",
			ExpectError:   false,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"name":               tc.BeforeName,
				"min_master_version": "",
				"node_version":       "",
			},
			After: map[string]interface{}{
				"name":               tc.AfterName,
				"min_master_version": tc.MasterVersion,
				"node_version":       tc.NodeVersion,
			},
		}
		err := containerClusterNodeVersionCustomizeDiffFunc(d)

		if tc.ExpectError && err == nil {
			t.Errorf("%s failed, expected error but was none", tn)
		}
		if !tc.ExpectError && err != nil {
			t.Errorf("%s failed, found unexpected error: %s", tn, err)
		}
	}
}
