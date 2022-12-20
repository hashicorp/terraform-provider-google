package google

import (
	"strings"
	"testing"
)

func TestValidateIAMWorkforcePoolWorkforcePoolProviderId(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foo-"},
		{TestName: "long", Value: strings.Repeat("f", 32)},
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

	es := testStringValidationCases(x, validateWorkforcePoolProviderId)
	if len(es) > 0 {
		t.Errorf("Failed to validate WorkforcePoolProvider names: %v", es)
	}
}
