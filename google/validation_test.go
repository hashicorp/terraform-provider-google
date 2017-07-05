package google

import (
	"fmt"
	"testing"
)

func TestValidateGCPName(t *testing.T) {
	x := []GCPNameTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "f"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a capital", Value: "Foobar", ExpectError: true},
		{TestName: "starts with a number", Value: "1foobar", ExpectError: true},
		{TestName: "has an underscore", Value: "foo_bar", ExpectError: true},
		{TestName: "too long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoob", ExpectError: true},
	}

	es := testGCPNames(x)
	if len(es) > 0 {
		t.Errorf("Failed to validate GCP names: %v", es)
	}
}

type GCPNameTestCase struct {
	TestName    string
	Value       string
	ExpectError bool
}

func testGCPNames(cases []GCPNameTestCase) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testGCPName(c)...)
	}

	return es
}

func testGCPName(testCase GCPNameTestCase) []error {
	_, es := validateGCPName(testCase.Value, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		} else {
			return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with string \"%s\"", testCase.TestName, testCase.Value)}
		}
	}

	return es
}
