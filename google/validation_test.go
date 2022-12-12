package google

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func TestValidateGCEName(t *testing.T) {
	x := []StringValidationTestCase{
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

	es := testStringValidationCases(x, validateGCEName)
	if len(es) > 0 {
		t.Errorf("Failed to validate GCP names: %v", es)
	}
}

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

func TestValidateRFC3339Time(t *testing.T) {
	cases := []StringValidationTestCase{
		// No errors
		{TestName: "midnight", Value: "00:00"},
		{TestName: "one minute before midnight", Value: "23:59"},

		// With errors
		{TestName: "single-digit hour", Value: "3:00", ExpectError: true},
		{TestName: "hour out of range", Value: "24:00", ExpectError: true},
		{TestName: "minute out of range", Value: "03:60", ExpectError: true},
		{TestName: "missing colon", Value: "0100", ExpectError: true},
		{TestName: "not numbers", Value: "ab:cd", ExpectError: true},
	}

	es := testStringValidationCases(cases, validateRFC3339Time)
	if len(es) > 0 {
		t.Errorf("Failed to validate RFC3339 times: %v", es)
	}
}

func TestValidateRFC1035Name(t *testing.T) {
	cases := []struct {
		TestName    string
		Value       string
		Min, Max    int
		ExpectError bool
	}{
		{TestName: "valid", Min: 6, Max: 30, Value: "a-valid-name0"},
		{TestName: "valid lower bound", Min: 12, Max: 30, Value: "a-valid-name"},
		{TestName: "valid upper bound", Min: 6, Max: 12, Value: "a-valid-name"},
		{TestName: "valid with numbers", Min: 6, Max: 30, Value: "valid000-name"},
		{TestName: "must start with a letter", Min: 6, Max: 10, Value: "0invalid", ExpectError: true},
		{TestName: "cannot end with a dash", Min: 6, Max: 10, Value: "invalid-", ExpectError: true},
		{TestName: "too short", Min: 6, Max: 10, Value: "short", ExpectError: true},
		{TestName: "too long", Min: 6, Max: 10, Value: "toolooooong", ExpectError: true},
		{TestName: "min too small", Min: 1, Max: 10, Value: "", ExpectError: true},
		{TestName: "min < max", Min: 6, Max: 5, Value: "", ExpectError: true},
	}

	for _, c := range cases {
		errors := testStringValidation(StringValidationTestCase{
			TestName:    c.TestName,
			Value:       c.Value,
			ExpectError: c.ExpectError,
		}, validateRFC1035Name(c.Min, c.Max))

		if len(errors) > 0 {
			t.Errorf("%s failed; %v", c.TestName, errors)
		}
	}
}

func TestValidateServiceAccountLink(t *testing.T) {
	cases := []StringValidationTestCase{
		// These test cases focus on the project name part of the regex
		// The service account name is covered by the RFC1035Name tests above

		// No errors
		{TestName: "valid with dash", Value: "projects/my-project/serviceAccounts/svcacct@my-project.iam.gserviceaccount.com"},
		{TestName: "valid with colon", Value: "projects/my:project/serviceAccounts/svcacct@project.my.iam.gserviceaccount.com"},
		{TestName: "valid with dot and colon", Value: "projects/my.thing:project/serviceAccounts/svcacct@project.my.thing.iam.gserviceaccount.com"},
		{TestName: "valid with compute default service account", Value: "projects/my-project/serviceAccounts/123456-compute@developer.gserviceaccount.com"},
		{TestName: "valid with app engine default service account", Value: "projects/my-project/serviceAccounts/my-project@appspot.gserviceaccount.com"},

		// Errors
		{TestName: "multiple colons", Value: "projects/my:project:thing/serviceAccounts/svcacct@thing.project.my.iam.gserviceaccount.com", ExpectError: true},
		{TestName: "project name empty", Value: "projects//serviceAccounts/svcacct@.iam.gserviceaccount.com", ExpectError: true},
		{TestName: "dot only with no colon", Value: "projects/my.project/serviceAccounts/svcacct@my.project.iam.gserviceaccount.com", ExpectError: true},
		{
			TestName: "too long",
			Value: "projects/foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoob/serviceAccounts/svcacct@" +
				"foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoob.iam.gserviceaccount.com",
			ExpectError: true,
		},
	}

	es := testStringValidationCases(cases, validateRegexp(ServiceAccountLinkRegex))
	if len(es) > 0 {
		t.Errorf("Failed to validate Service Account Links: %v", es)
	}
}

type StringValidationTestCase struct {
	TestName    string
	Value       string
	ExpectError bool
}

type RFC1918NetworkTestCase struct {
	TestName    string
	CIDR        string
	MinPrefix   int
	MaxPrefix   int
	ExpectError bool
}

func testStringValidationCases(cases []StringValidationTestCase, validationFunc schema.SchemaValidateFunc) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testStringValidation(c, validationFunc)...)
	}

	return es
}

func testStringValidation(testCase StringValidationTestCase, validationFunc schema.SchemaValidateFunc) []error {
	_, es := validationFunc(testCase.Value, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		} else {
			return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with string \"%s\"", testCase.TestName, testCase.Value)}
		}
	}

	return es
}

func testRFC1918Networks(cases []RFC1918NetworkTestCase) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testRFC1918Network(c)...)
	}

	return es
}

func testRFC1918Network(testCase RFC1918NetworkTestCase) []error {
	f := validateRFC1918Network(testCase.MinPrefix, testCase.MaxPrefix)
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
		validateFunc := orEmpty(validation.StringInSlice([]string{"valid"}, false))
		_, errors := validateFunc(tc.Value, tn)
		if len(errors) > 0 && !tc.ExpectValidationErrors {
			t.Errorf("%s: unexpected errors %s", tn, errors)
		} else if len(errors) == 0 && tc.ExpectValidationErrors {
			t.Errorf("%s: expected errors but got none", tn)
		}
	}
}

func TestValidateProjectID(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foofoo"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobar"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an uppercase letter", Value: "foo-Bar", ExpectError: true},
		{TestName: "has a final hyphen", Value: "foo-bar-", ExpectError: true},
	}

	es := testStringValidationCases(x, validateProjectID())
	if len(es) > 0 {
		t.Errorf("Failed to validate project ID's: %v", es)
	}
}

func TestValidateDSProjectID(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foofoo"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobar"},
		{TestName: "has projects", Value: "projects/foo-bar"},
		{TestName: "has multiple projects", Value: "projects/projects/foobar"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "has an uppercase letter", Value: "foo-Bar", ExpectError: true},
		{TestName: "has a final hyphen", Value: "foo-bar-", ExpectError: true},
	}

	es := testStringValidationCases(x, validateDSProjectID())
	if len(es) > 0 {
		t.Errorf("Failed to validate project ID's: %v", es)
	}
}

func TestValidateProjectName(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "fooBar"},
		{TestName: "complex", Value: "project! 'A-1234'"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foof"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobar"},
		{TestName: "has a hyphen", Value: "foo-bar"},
		{TestName: "starts with a number", Value: "1foobar"},
		{TestName: "has a final hyphen", Value: "foo-bar-"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("a", 31), ExpectError: true},
	}

	es := testStringValidationCases(x, validateProjectName())
	if len(es) > 0 {
		t.Errorf("Failed to validate project ID's: %v", es)
	}
}

func TestValidateIAMCustomRoleIDRegex(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "with capitals", Value: "FooBar"},
		{TestName: "short", Value: "foo"},
		{TestName: "long", Value: strings.Repeat("f", 64)},
		{TestName: "has a dot", Value: "foo.bar"},
		{TestName: "has an underscore", Value: "foo_bar"},
		{TestName: "all of the above", Value: "foo.BarBaz_123"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has a hyphen", Value: "foo-bar", ExpectError: true},
		{TestName: "has a dollar", Value: "foo$", ExpectError: true},
		{TestName: "has a space", Value: "foo bar", ExpectError: true},
		{TestName: "too short", Value: "fo", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 65), ExpectError: true},
	}

	es := testStringValidationCases(x, validateIAMCustomRoleID)
	if len(es) > 0 {
		t.Errorf("Failed to validate IAMCustomRole IDs: %v", es)
	}
}
