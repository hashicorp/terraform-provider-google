package google

// Returns the terraform representation of a three-state boolean value represented by a pointer to bool in DCL.
func flattenEnumBool(v interface{}) string {
	b, ok := v.(*bool)
	if !ok || b == nil {
		return ""
	}
	if *b {
		return "TRUE"
	}
	return "FALSE"
}
