// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"testing"
)

func TestLabelDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		K, Old, New        string
		ExpectDiffSuppress bool
	}{
		"missing goog-dataplex-asset-id": {
			K:                  "labels.goog-dataplex-asset-id",
			Old:                "test-bucket",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-asset-id": {
			K:                  "labels.goog-dataplex-asset-id",
			Old:                "test-bucket",
			New:                "test-bucket-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-lake-id": {
			K:                  "labels.goog-dataplex-lake-id",
			Old:                "test-lake",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-lake-id": {
			K:                  "labels.goog-dataplex-lake-id",
			Old:                "test-lake",
			New:                "test-lake-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-project-id": {
			K:                  "labels.goog-dataplex-project-id",
			Old:                "test-project-12345",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-project-id": {
			K:                  "labels.goog-dataplex-project-id",
			Old:                "test-project-12345",
			New:                "test-project-12345-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-zone-id": {
			K:                  "labels.goog-dataplex-zone-id",
			Old:                "test-zone1",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-zone-id": {
			K:                  "labels.goog-dataplex-zone-id",
			Old:                "test-zone1",
			New:                "test-zone1-1",
			ExpectDiffSuppress: false,
		},
		"labels.%": {
			K:                  "labels.%",
			Old:                "5",
			New:                "1",
			ExpectDiffSuppress: true,
		},
		"deleted custom key": {
			K:                  "labels.my-label",
			Old:                "my-value",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"added custom key": {
			K:                  "labels.my-label",
			Old:                "",
			New:                "my-value",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if resourceDataplexLabelDiffSuppress(tc.K, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q: %q => %q expect DiffSuppress to return %t", tn, tc.K, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
