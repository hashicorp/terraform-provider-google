package google

import (
	"fmt"
	"regexp"
)

func validateGCPName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	re := `^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$`
	if !regexp.MustCompile(re).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) doesn't match regexp %q", k, value, re))
	}
	return
}
