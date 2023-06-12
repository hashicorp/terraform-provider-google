// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"cloud.google.com/go/bigtable"
	tpgbigtable "github.com/hashicorp/terraform-provider-google/google/services/bigtable"
)

// Recursively convert Bigtable GC policy to JSON format in a map.
func gcPolicyToGCRuleString(gc bigtable.GCPolicy, topLevel bool) (map[string]interface{}, error) {
	return tpgbigtable.GcPolicyToGCRuleString(gc, topLevel)
}
