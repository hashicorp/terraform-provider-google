// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import "testing"

func TestParseLoggingSinkId(t *testing.T) {
	tests := []struct {
		val         string
		out         *LoggingSinkId
		errExpected bool
	}{
		{"projects/my-project/sinks/my-sink", &LoggingSinkId{"projects", "my-project", "my-sink"}, false},
		{"folders/foofolder/sinks/woo", &LoggingSinkId{"folders", "foofolder", "woo"}, false},
		{"kitchens/the-big-one/sinks/second-from-the-left", nil, true},
	}

	for _, test := range tests {
		out, err := ParseLoggingSinkId(test.val)
		if err != nil {
			if !test.errExpected {
				t.Errorf("Got error with val %#v: error = %#v", test.val, err)
			}
		} else {
			if *out != *test.out {
				t.Errorf("Mismatch on val %#v: expected %#v but got %#v", test.val, test.out, out)
			}
		}
	}
}

func TestLoggingSinkId(t *testing.T) {
	tests := []struct {
		val         LoggingSinkId
		canonicalId string
		parent      string
	}{
		{
			val:         LoggingSinkId{"projects", "my-project", "my-sink"},
			canonicalId: "projects/my-project/sinks/my-sink",
			parent:      "projects/my-project",
		}, {
			val:         LoggingSinkId{"folders", "foofolder", "woo"},
			canonicalId: "folders/foofolder/sinks/woo",
			parent:      "folders/foofolder",
		},
	}

	for _, test := range tests {
		canonicalId := test.val.canonicalId()

		if canonicalId != test.canonicalId {
			t.Errorf("canonicalId mismatch on val %#v: expected %#v but got %#v", test.val, test.canonicalId, canonicalId)
		}

		parent := test.val.parent()

		if parent != test.parent {
			t.Errorf("parent mismatch on val %#v: expected %#v but got %#v", test.val, test.parent, parent)
		}
	}
}
