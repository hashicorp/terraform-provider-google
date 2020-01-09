package google

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	// Copied from the official Google Cloud auto-generated client.
	ProjectRegex         = "(?:(?:[-a-z0-9]{1,63}\\.)*(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?):)?(?:[0-9]{1,19}|(?:[a-z0-9](?:[-a-z0-9]{0,61}[a-z0-9])?))"
	ProjectRegexWildCard = "(?:(?:[-a-z0-9]{1,63}\\.)*(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?):)?(?:[0-9]{1,19}|(?:[a-z0-9](?:[-a-z0-9]{0,61}[a-z0-9])?)|-)"
	RegionRegex          = "[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?"
	SubnetworkRegex      = "[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?"

	SubnetworkLinkRegex = "projects/(" + ProjectRegex + ")/regions/(" + RegionRegex + ")/subnetworks/(" + SubnetworkRegex + ")$"

	RFC1035NameTemplate = "[a-z](?:[-a-z0-9]{%d,%d}[a-z0-9])"
	CloudIoTIdRegex     = "^[a-zA-Z][-a-zA-Z0-9._+~%]{2,254}$"

	// Format of default Compute service accounts created by Google
	// ${PROJECT_ID}-compute@developer.gserviceaccount.com where PROJECT_ID is an int64 (max 20 digits)
	ComputeServiceAccountNameRegex = "[0-9]{1,20}-compute@developer.gserviceaccount.com"

	// https://cloud.google.com/iam/docs/understanding-custom-roles#naming_the_role
	IAMCustomRoleIDRegex = "^[a-zA-Z0-9_\\.]{3,64}$"
)

var (
	// Service account name must have a length between 6 and 30.
	// The first and last characters have different restrictions, than
	// the middle characters. The middle characters length must be between
	// 4 and 28 since the first and last character are excluded.
	ServiceAccountNameRegex = fmt.Sprintf(RFC1035NameTemplate, 4, 28)

	ServiceAccountLinkRegexPrefix = "projects/" + ProjectRegexWildCard + "/serviceAccounts/"
	PossibleServiceAccountNames   = []string{
		ServiceDefaultAccountNameRegex,
		ComputeServiceAccountNameRegex,
		CreatedServiceAccountNameRegex,
	}
	ServiceAccountLinkRegex = ServiceAccountLinkRegexPrefix + "(" + strings.Join(PossibleServiceAccountNames, "|") + ")"

	ServiceAccountKeyNameRegex = ServiceAccountLinkRegexPrefix + "(.+)/keys/(.+)"

	// Format of service accounts created through the API
	CreatedServiceAccountNameRegex = fmt.Sprintf(RFC1035NameTemplate, 4, 28) + "@" + ProjectNameInDNSFormRegex + "\\.iam\\.gserviceaccount\\.com$"

	// Format of service-created service account
	// examples are:
	// 		$PROJECTID@cloudbuild.gserviceaccount.com
	// 		$PROJECTID@cloudservices.gserviceaccount.com
	// 		$PROJECTID@appspot.gserviceaccount.com
	ServiceDefaultAccountNameRegex = ProjectRegex + "@[a-z]+.gserviceaccount.com$"

	ProjectNameInDNSFormRegex = "[-a-z0-9\\.]{1,63}"
	ProjectNameRegex          = "^[A-Za-z0-9-'\"\\s!]{4,30}$"
)

var rfc1918Networks = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

func validateGCPName(v interface{}, k string) (ws []string, errors []error) {
	re := `^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$`
	return validateRegexp(re)(v, k)
}

func validateRegexp(re string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !regexp.MustCompile(re).MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q (%q) doesn't match regexp %q", k, value, re))
		}

		return
	}
}

func validateRFC1918Network(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {

		s, es = validation.CIDRNetwork(min, max)(i, k)
		if len(es) > 0 {
			return
		}

		v, _ := i.(string)
		ip, _, _ := net.ParseCIDR(v)
		for _, c := range rfc1918Networks {
			if _, ipnet, _ := net.ParseCIDR(c); ipnet.Contains(ip) {
				return
			}
		}

		es = append(es, fmt.Errorf("expected %q to be an RFC1918-compliant CIDR, got: %s", k, v))

		return
	}
}

func validateRFC3339Time(v interface{}, k string) (warnings []string, errors []error) {
	time := v.(string)
	if len(time) != 5 || time[2] != ':' {
		errors = append(errors, fmt.Errorf("%q (%q) must be in the format HH:mm (RFC3399)", k, time))
		return
	}
	if hour, err := strconv.ParseUint(time[:2], 10, 0); err != nil || hour > 23 {
		errors = append(errors, fmt.Errorf("%q (%q) does not contain a valid hour (00-23)", k, time))
		return
	}
	if min, err := strconv.ParseUint(time[3:], 10, 0); err != nil || min > 59 {
		errors = append(errors, fmt.Errorf("%q (%q) does not contain a valid minute (00-59)", k, time))
		return
	}
	return
}

func validateRFC1035Name(min, max int) schema.SchemaValidateFunc {
	if min < 2 || max < min {
		return func(i interface{}, k string) (s []string, errors []error) {
			if min < 2 {
				errors = append(errors, fmt.Errorf("min must be at least 2. Got: %d", min))
			}
			if max < min {
				errors = append(errors, fmt.Errorf("max must greater than min. Got [%d, %d]", min, max))
			}
			return
		}
	}

	return validateRegexp(fmt.Sprintf("^"+RFC1035NameTemplate+"$", min-2, max-2))
}

func validateIpCidrRange(v interface{}, k string) (warnings []string, errors []error) {
	_, _, err := net.ParseCIDR(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IP CIDR range: %s", k, err))
	}
	return
}

func validateIAMCustomRoleID(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(IAMCustomRoleIDRegex).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) doesn't match regexp %q", k, value, IAMCustomRoleIDRegex))
	}
	return
}

func orEmpty(f schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if ok && v == "" {
			return nil, nil
		}
		return f(i, k)
	}
}

func validateProjectID() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		if !regexp.MustCompile("^" + ProjectRegex + "$").MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q project_id must be 6 to 30 with lowercase letters, digits, hyphens and start with a letter. Trailing hyphens are prohibited.", value))
		}
		return
	}
}

func validateProjectName() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		if !regexp.MustCompile(ProjectNameRegex).MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q name must be 4 to 30 characters with lowercase and uppercase letters, numbers, hyphen, single-quote, double-quote, space, and exclamation point.", value))
		}
		return
	}
}

func validateDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if _, err := time.ParseDuration(v); err != nil {
			es = append(es, fmt.Errorf("expected %s to be a duration, but parsing gave an error: %s", k, err.Error()))
			return
		}

		return
	}
}

func validateNonNegativeDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		dur, err := time.ParseDuration(v)
		if err != nil {
			es = append(es, fmt.Errorf("expected %s to be a duration, but parsing gave an error: %s", k, err.Error()))
			return
		}

		if dur < 0 {
			es = append(es, fmt.Errorf("duration %v must be a non-negative duration", dur))
			return
		}

		return
	}
}

func validateIpAddress(i interface{}, val string) ([]string, []error) {
	ip := net.ParseIP(i.(string))
	if ip == nil {
		return nil, []error{fmt.Errorf("could not parse %q to IP address", val)}
	}
	return nil, nil
}

// StringNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and that it matches none of the element in the invalid slice.
// if ignorecase is true, case is ignored.
func StringNotInSlice(invalid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		for _, str := range invalid {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				es = append(es, fmt.Errorf("expected %s to not match any of %v, got %s", k, invalid, v))
				return
			}
		}

		return
	}
}

// Ensure that hourly timestamp strings "HH:MM" have the minutes zeroed out for hourly only inputs
func validateHourlyOnly(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		errs = append(errs, fmt.Errorf("%q must be in the format HH:00, got: %s", key, v))
		return
	}
	if parts[1] != "00" {
		errs = append(errs, fmt.Errorf("%q does not allow minutes, it must be in the format HH:00, got: %s", key, v))
	}
	i, err := strconv.Atoi(parts[0])
	if err != nil {
		errs = append(errs, fmt.Errorf("%q cannot be parsed, it must be in the format HH:00, got: %s", key, v))
	} else if i < 0 || i > 23 {
		errs = append(errs, fmt.Errorf("%q does not specify a valid hour, it must be in the format HH:00 where HH : [00-23], got: %s", key, v))
	}
	return
}
