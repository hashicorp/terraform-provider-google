// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataprocmetastore

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const dataprocMetastoreProvidedOverride = "hive.metastore.warehouse.dir"

func dataprocMetastoreServiceOverrideSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the label provided by Google
	if strings.Contains(k, dataprocMetastoreProvidedOverride) && new == "" {
		return true
	}

	// Let diff be determined by labels (above)
	if strings.Contains(k, "hive_metastore_config.0.config_overrides.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}
