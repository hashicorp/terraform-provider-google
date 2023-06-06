// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestDatastreamStreamCustomDiff(t *testing.T) {
	t.Parallel()

	cases := []struct {
		isNew     bool
		old       string
		new       string
		wantError bool
	}{
		{
			isNew:     true,
			new:       "NOT_STARTED",
			wantError: false,
		},
		{
			isNew:     true,
			new:       "RUNNING",
			wantError: false,
		},
		{
			isNew:     true,
			new:       "PAUSED",
			wantError: true,
		},
		{
			isNew:     true,
			new:       "MAINTENANCE",
			wantError: true,
		},
		{
			// Normally this transition is okay, but if the resource is "new"
			// (for example being recreated) it's not.
			isNew:     true,
			old:       "RUNNING",
			new:       "PAUSED",
			wantError: true,
		},
		{
			old:       "NOT_STARTED",
			new:       "RUNNING",
			wantError: false,
		},
		{
			old:       "NOT_STARTED",
			new:       "MAINTENANCE",
			wantError: true,
		},
		{
			old:       "NOT_STARTED",
			new:       "PAUSED",
			wantError: true,
		},
		{
			old:       "NOT_STARTED",
			new:       "NOT_STARTED",
			wantError: false,
		},
		{
			old:       "RUNNING",
			new:       "PAUSED",
			wantError: false,
		},
		{
			old:       "RUNNING",
			new:       "NOT_STARTED",
			wantError: true,
		},
		{
			old:       "RUNNING",
			new:       "RUNNING",
			wantError: false,
		},
		{
			old:       "RUNNING",
			new:       "MAINTENANCE",
			wantError: true,
		},
		{
			old:       "PAUSED",
			new:       "PAUSED",
			wantError: false,
		},
		{
			old:       "PAUSED",
			new:       "NOT_STARTED",
			wantError: true,
		},
		{
			old:       "PAUSED",
			new:       "RUNNING",
			wantError: false,
		},
		{
			old:       "PAUSED",
			new:       "MAINTENANCE",
			wantError: true,
		},
	}
	for _, tc := range cases {
		name := "whatever"
		tn := fmt.Sprintf("%s => %s", tc.old, tc.new)
		if tc.isNew {
			name = ""
			tn = fmt.Sprintf("(new) %s => %s", tc.old, tc.new)
		}
		t.Run(tn, func(t *testing.T) {
			diff := &tpgresource.ResourceDiffMock{
				Before: map[string]interface{}{
					"desired_state": tc.old,
				},
				After: map[string]interface{}{
					"name":          name,
					"desired_state": tc.new,
				},
			}
			err := resourceDatastreamStreamCustomDiffFunc(diff)
			if tc.wantError && err == nil {
				t.Fatalf("want error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
		})
	}
}
