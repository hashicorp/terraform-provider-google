package google

import (
	"strings"
	"testing"
)

func TestValidateIAMWorkforcePoolWorkforcePoolId(t *testing.T) {
	x := []StringValidationTestCase{
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

	es := testStringValidationCases(x, validateWorkforcePoolId)
	if len(es) > 0 {
		t.Errorf("Failed to validate WorkforcePool names: %v", es)
	}
}
