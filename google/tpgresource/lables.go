// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FlattenLabels(labels map[string]string, d *schema.ResourceData) map[string]interface{} {
	transformed := make(map[string]interface{})

	if v, ok := d.GetOk("labels"); ok {
		if labels != nil {
			for k, _ := range v.(map[string]interface{}) {
				transformed[k] = labels[k]
			}
		}
	}

	return transformed
}
