// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudiot_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/cloudiot"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func TestValidateCloudIoTDeviceRegistryId(t *testing.T) {
	x := []verify.StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foo"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a goog", Value: "googfoobar", ExpectError: true},
		{TestName: "starts with a number", Value: "1foobar", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an backslash", Value: "foo\bar", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 260), ExpectError: true},
	}

	es := verify.TestStringValidationCases(x, cloudiot.ValidateCloudIotDeviceRegistryID)
	if len(es) > 0 {
		t.Errorf("Failed to validate CloudIoT ID names: %v", es)
	}
}
