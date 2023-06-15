// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/bigtable"
)

func flattenColumnFamily(families []string) []map[string]interface{} {
	return bigtable.FlattenColumnFamily(families)
}
