// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iamworkforcepool_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/iamworkforcepool"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func TestValidateIAMWorkforcePoolWorkforcePoolId(t *testing.T) {
	x := []verify.StringValidationTestCase{
		// No errors
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foobar"},
		{TestName: "long", Value: strings.Repeat("f", 63)},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a gcp-", Value: "gcp-foobar", ExpectError: true},
		{TestName: "with uppercase", Value: "fooBar", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an backslash", Value: "foo\bar", ExpectError: true},
		{TestName: "too short", Value: "foooo", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 64), ExpectError: true},
		{TestName: "doesn't start with a lowercase letter", Value: "123foo", ExpectError: true},
		{TestName: "ends with a hyphen", Value: "foobar-", ExpectError: true},
	}

	es := verify.TestStringValidationCases(x, iamworkforcepool.ValidateWorkforcePoolId)
	if len(es) > 0 {
		t.Errorf("Failed to validate WorkforcePool names: %v", es)
	}
}
