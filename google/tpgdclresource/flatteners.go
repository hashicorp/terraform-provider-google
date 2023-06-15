// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgdclresource

// Returns the terraform representation of a three-state boolean value represented by a pointer to bool in DCL.
func FlattenEnumBool(v interface{}) string {
	b, ok := v.(*bool)
	if !ok || b == nil {
		return ""
	}
	if *b {
		return "TRUE"
	}
	return "FALSE"
}
