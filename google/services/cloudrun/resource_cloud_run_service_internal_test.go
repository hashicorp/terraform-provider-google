// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrun

import (
	"testing"
)

func TestCloudrunAnnotationDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		K, Old, New        string
		ExpectDiffSuppress bool
	}{
		"missing run.googleapis.com/operation-id": {
			K:                  "metadata.0.annotations.run.googleapis.com/operation-id",
			Old:                "12345abc",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"missing run.googleapis.com/ingress": {
			K:                  "metadata.0.annotations.run.googleapis.com/ingress",
			Old:                "all",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit run.googleapis.com/ingress": {
			K:                  "metadata.0.annotations.run.googleapis.com/ingress",
			Old:                "all",
			New:                "internal",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if got := cloudrunAnnotationDiffSuppress(tc.K, tc.Old, tc.New, nil); got != tc.ExpectDiffSuppress {
				t.Errorf("got %t; want %t", got, tc.ExpectDiffSuppress)
			}
		})
	}
}
