package verify

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func TestValidateRFC1918Network(t *testing.T) {
	x := []RFC1918NetworkTestCase{
		// No errors
		{TestName: "valid 10.x", CIDR: "10.0.0.0/8", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid 172.x", CIDR: "172.16.0.0/16", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid 192.x", CIDR: "192.168.0.0/32", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid, bounded 10.x CIDR", CIDR: "10.0.0.0/8", MinPrefix: 8, MaxPrefix: 32},
		{TestName: "valid, bounded 172.x CIDR", CIDR: "172.16.0.0/16", MinPrefix: 12, MaxPrefix: 32},
		{TestName: "valid, bounded 192.x CIDR", CIDR: "192.168.0.0/32", MinPrefix: 16, MaxPrefix: 32},

		// With errors
		{TestName: "empty CIDR", CIDR: "", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "missing mask", CIDR: "10.0.0.0", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "invalid CIDR", CIDR: "10.1.0.0/8", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "valid 10.x CIDR with lower bound violation", CIDR: "10.0.0.0/8", MinPrefix: 16, MaxPrefix: 32, ExpectError: true},
		{TestName: "valid 10.x CIDR with upper bound violation", CIDR: "10.0.0.0/24", MinPrefix: 8, MaxPrefix: 16, ExpectError: true},
		{TestName: "valid public CIDR", CIDR: "8.8.8.8/32", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
	}

	es := testRFC1918Networks(x)
	if len(es) > 0 {
		t.Errorf("Failed to validate RFC1918 Networks: %v", es)
	}
}

type RFC1918NetworkTestCase struct {
	TestName    string
	CIDR        string
	MinPrefix   int
	MaxPrefix   int
	ExpectError bool
}

func testRFC1918Networks(cases []RFC1918NetworkTestCase) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testRFC1918Network(c)...)
	}

	return es
}

func testRFC1918Network(testCase RFC1918NetworkTestCase) []error {
	f := ValidateRFC1918Network(testCase.MinPrefix, testCase.MaxPrefix)
	_, es := f(testCase.CIDR, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		}
		return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with CIDR=\"%s\" MinPrefix=%v MaxPrefix=%v",
			testCase.TestName, testCase.CIDR, testCase.MinPrefix, testCase.MaxPrefix)}
	}

	return es
}

func TestProjectRegex(t *testing.T) {
	tests := []struct {
		project string
		want    bool
	}{
		{"example", true},
		{"example.com:example", true},
		{"12345", true},
		{"", false},
		{"example_", false},
	}
	r := regexp.MustCompile("^" + ProjectRegex + "$")
	for _, test := range tests {
		if got := r.MatchString(test.project); got != test.want {
			t.Errorf("got %t, want %t for project %v", got, test.want, test.project)
		}
	}
}

func TestOrEmpty(t *testing.T) {
	cases := map[string]struct {
		Value                  string
		ValidateFunc           schema.SchemaValidateFunc
		ExpectValidationErrors bool
	}{
		"accept empty value": {
			Value:                  "",
			ExpectValidationErrors: false,
		},
		"non empty value is accepted when valid": {
			Value:                  "valid",
			ExpectValidationErrors: false,
		},
		"non empty value is rejected if invalid": {
			Value:                  "invalid",
			ExpectValidationErrors: true,
		},
	}

	for tn, tc := range cases {
		validateFunc := OrEmpty(validation.StringInSlice([]string{"valid"}, false))
		_, errors := validateFunc(tc.Value, tn)
		if len(errors) > 0 && !tc.ExpectValidationErrors {
			t.Errorf("%s: unexpected errors %s", tn, errors)
		} else if len(errors) == 0 && tc.ExpectValidationErrors {
			t.Errorf("%s: expected errors but got none", tn)
		}
	}
}
