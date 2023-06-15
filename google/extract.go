// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/logging"
)

// ExtractFieldByPattern returns the value of a field extracted from a parent field according to the given regular expression pattern.
// An error is returned if the field already has a value different than the value extracted.
func ExtractFieldByPattern(fieldName, fieldValue, parentFieldValue, pattern string) (string, error) {
	return logging.ExtractFieldByPattern(fieldName, fieldValue, parentFieldValue, pattern)
}
