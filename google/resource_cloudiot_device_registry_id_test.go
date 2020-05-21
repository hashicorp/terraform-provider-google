package google

import (
	"strings"
	"testing"
)

func TestValidateCloudIoTDeviceRegistryId(t *testing.T) {
	x := []StringValidationTestCase{
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

	es := testStringValidationCases(x, validateCloudIotDeviceRegistryID)
	if len(es) > 0 {
		t.Errorf("Failed to validate CloudIoT ID names: %v", es)
	}
}
