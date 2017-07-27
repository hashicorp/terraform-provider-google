package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

const projectRegex = "(?:(?:[-a-z0-9]{1,63}\\.)*(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?):)?(?:[0-9]{1,19}|(?:[a-z0-9](?:[-a-z0-9]{0,61}[a-z0-9])?))"

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
