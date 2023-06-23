// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner

import (
	"testing"
)

// Unit Tests

func TestSpannerInstanceId_instanceUri(t *testing.T) {
	id := SpannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceUri()
	expected := "projects/project123/instances/instance456"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_instanceConfigUri(t *testing.T) {
	id := SpannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceConfigUri("conf987")
	expected := "projects/project123/instanceConfigs/conf987"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_parentProjectUri(t *testing.T) {
	id := SpannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.parentProjectUri()
	expected := "projects/project123"
	expectEquals(t, expected, actual)
}

func expectEquals(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}
