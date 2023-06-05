// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgdclresource

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func ExpandStringArray(v interface{}) []string {
	arr, ok := v.([]string)

	if ok {
		return arr
	}

	if arr, ok := v.(*schema.Set); ok {
		return tpgresource.ConvertStringSet(arr)
	}

	arr = tpgresource.ConvertStringArr(v.([]interface{}))
	if arr == nil {
		// Send empty array specifically instead of nil
		return make([]string, 0)
	}
	return arr
}

func ExpandIntegerArray(v interface{}) []int64 {
	arr, ok := v.([]int64)

	if ok {
		return arr
	}

	if arr, ok := v.(*schema.Set); ok {
		return convertIntegerSet(arr)
	}

	return convertIntegerArr(v.([]interface{}))
}

func convertIntegerSet(v *schema.Set) []int64 {
	return convertIntegerArr(v.List())
}

func convertIntegerArr(v []interface{}) []int64 {
	var vi []int64
	for _, vs := range v {
		vi = append(vi, int64(vs.(int)))
	}
	return vi
}

// Returns the DCL representation of a three-state boolean value represented by a string in terraform.
func ExpandEnumBool(v interface{}) *bool {
	s, ok := v.(string)
	if !ok {
		return nil
	}

	switch {
	case strings.EqualFold(s, "true"):
		return boolPtr(true)
	case strings.EqualFold(s, "false"):
		return boolPtr(false)
	default:
		return nil
	}
}

// boolPtr returns a pointer to the given boolean.
func boolPtr(b bool) *bool {
	return &b
}
