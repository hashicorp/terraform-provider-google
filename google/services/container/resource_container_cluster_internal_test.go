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
