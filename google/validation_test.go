package google

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/verify"
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

	es := testStringValidationCases(x, verify.ValidateGCEName)
	if len(es) > 0 {
		t.Errorf("Failed to validate GCP names: %v", es)
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

	es := testStringValidationCases(cases, verify.ValidateRFC3339Time)
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
		}, verify.ValidateRFC1035Name(c.Min, c.Max))

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

	es := testStringValidationCases(cases, verify.ValidateRegexp(verify.ServiceAccountLinkRegex))
	if len(es) > 0 {
		t.Errorf("Failed to validate Service Account Links: %v", es)
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

	es := testStringValidationCases(x, verify.ValidateProjectID())
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

	es := testStringValidationCases(x, verify.ValidateDSProjectID())
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

	es := testStringValidationCases(x, verify.ValidateProjectName())
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

	es := testStringValidationCases(x, verify.ValidateIAMCustomRoleID)
	if len(es) > 0 {
		t.Errorf("Failed to validate IAMCustomRole IDs: %v", es)
	}
}
