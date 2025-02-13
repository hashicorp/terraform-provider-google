// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iambeta_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/iambeta"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func TestValidateIAMBetaWorkloadIdentityPoolId(t *testing.T) {
	x := []verify.StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foos"},
		{TestName: "long", Value: "12345678901234567890123456789012"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a gcp-", Value: "gcp-foobar", ExpectError: true},
		{TestName: "with uppercase", Value: "fooBar", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an backslash", Value: "foo\bar", ExpectError: true},
		{TestName: "too short", Value: "foo", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 33), ExpectError: true},
	}

	es := verify.TestStringValidationCases(x, iambeta.ValidateWorkloadIdentityPoolId)
	if len(es) > 0 {
		t.Errorf("Failed to validate WorkloadIdentityPool names: %v", es)
	}
}
