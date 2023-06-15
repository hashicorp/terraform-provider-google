// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"regexp"
)

// ExtractFieldByPattern returns the value of a field extracted from a parent field according to the given regular expression pattern.
// An error is returned if the field already has a value different than the value extracted.
func ExtractFieldByPattern(fieldName, fieldValue, parentFieldValue, pattern string) (string, error) {
	var extractedValue string
	// Fetch value from container if the container exists.
	if parentFieldValue != "" {
		r := regexp.MustCompile(pattern)
		m := r.FindStringSubmatch(parentFieldValue)
		if m != nil && len(m) >= 2 {
			extractedValue = m[1]
		} else if fieldValue == "" {
			// The pattern didn't match and the value doesn't exist.
			return "", fmt.Errorf("parent of %q has no matching values from pattern %q in value %q", fieldName, pattern, parentFieldValue)
		}
	}

	// If both values exist and are different, error
	if fieldValue != "" && extractedValue != "" && fieldValue != extractedValue {
		return "", fmt.Errorf("%q has conflicting values of %q (from parent) and %q (from self)", fieldName, extractedValue, fieldValue)
	}

	// If value does not exist, use the value in container.
	if fieldValue == "" {
		return extractedValue, nil
	}

	return fieldValue, nil
}
