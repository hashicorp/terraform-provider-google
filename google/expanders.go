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

	return convertStringArr(v.([]interface{}))
}
