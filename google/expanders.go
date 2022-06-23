package google

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func expandStringArray(v interface{}) []string {
	arr, ok := v.([]string)

	if ok {
		return arr
	}

	if arr, ok := v.(*schema.Set); ok {
		return convertStringSet(arr)
	}

	arr = convertStringArr(v.([]interface{}))
	if arr == nil {
		// Send empty array specifically instead of nil
		return make([]string, 0)
	}
	return arr
}

func expandIntegerArray(v interface{}) []int64 {
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
func expandEnumBool(v interface{}) *bool {
	s, ok := v.(string)
	if !ok {
		return nil
	}
	switch s {
	case "TRUE":
		b := true
		return &b
	case "FALSE":
		b := false
		return &b
	}
	return nil
}
