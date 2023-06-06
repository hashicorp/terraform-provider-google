// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func TestValidateRecordNameTrailingDot(t *testing.T) {
	cases := []verify.StringValidationTestCase{
		// No errors
		{TestName: "trailing dot", Value: "test-record.hashicorptest.com."},

		// With errors
		{TestName: "empty string", Value: "", ExpectError: true},
		{TestName: "no trailing dot", Value: "test-record.hashicorptest.com", ExpectError: true},
	}

	es := verify.TestStringValidationCases(cases, validateRecordNameTrailingDot)
	if len(es) > 0 {
		t.Errorf("Failed to validate DNS Record name with value: %v", es)
	}
}
