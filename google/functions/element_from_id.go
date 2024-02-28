// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

const noMatchesErrorSummary string = "No matches present in the input string"
const ambiguousMatchesWarningSummary string = "Ambiguous input string could contain more than one match"

// ValidateElementFromIdArguments is reusable validation logic used in provider-defined functions that use the getElementFromId function
func ValidateElementFromIdArguments(input string, regex *regexp.Regexp, pattern string, resp *function.RunResponse) {
	submatches := regex.FindAllStringSubmatchIndex(input, -1)

	// Zero matches means unusable input; error returned
	if len(submatches) == 0 {
		resp.Diagnostics.AddArgumentError(
			0,
			noMatchesErrorSummary,
			fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"%s\".", input, pattern),
		)
	}

	// >1 matches means input usable but not ideal; issue warning
	if len(submatches) > 1 {
		resp.Diagnostics.AddArgumentWarning(
			0,
			ambiguousMatchesWarningSummary,
			fmt.Sprintf("The input string \"%s\" contains more than one match for the pattern \"%s\". Terraform will use the first found match.", input, pattern),
		)
	}
}

// GetElementFromId is reusable logic that is used in multiple provider-defined functions for pulling elements out of self links and ids of resources and data sources
func GetElementFromId(input string, regex *regexp.Regexp, template string) string {
	submatches := regex.FindAllStringSubmatchIndex(input, -1)
	submatch := submatches[0] // Take the only / left-most submatch
	dst := []byte{}
	return string(regex.ExpandString(dst, template, input, submatch))
}
