// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgdclresource

import (
	"reflect"
	"testing"
)

func TestExpandEnumBool(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		exp   *bool
	}{
		{
			name:  "true",
			input: "true",
			exp:   boolPtr(true),
		},
		{
			name:  "TRUE",
			input: "TRUE",
			exp:   boolPtr(true),
		},
		{
			name:  "True",
			input: "True",
			exp:   boolPtr(true),
		},
		{
			name:  "false",
			input: "false",
			exp:   boolPtr(false),
		},
		{
			name:  "FALSE",
			input: "FALSE",
			exp:   boolPtr(false),
		},
		{
			name:  "False",
			input: "False",
			exp:   boolPtr(false),
		},
		{
			name:  "empty_string",
			input: "",
			exp:   nil,
		},
		{
			name:  "apple",
			input: "apple",
			exp:   nil,
		},
		{
			name:  "unicode",
			input: "🚀",
			exp:   nil,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got, want := ExpandEnumBool(tc.input), tc.exp; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}
