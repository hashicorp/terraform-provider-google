// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func FlattenContainerAwsNodePoolManagement(obj *containeraws.NodePoolManagement, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if obj == nil {
		return nil
	}
	transformed := make(map[string]interface{})

	if obj.AutoRepair == nil || obj.Empty() {
		transformed["auto_repair"] = false
	} else {
		transformed["auto_repair"] = obj.AutoRepair
	}

	return []interface{}{transformed}
}

func FlattenContainerAzureNodePoolManagement(obj *containerazure.NodePoolManagement, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if obj == nil {
		return nil
	}
	transformed := make(map[string]interface{})

	if obj.AutoRepair == nil || obj.Empty() {
		transformed["auto_repair"] = false
	} else {
		transformed["auto_repair"] = obj.AutoRepair
	}

	return []interface{}{transformed}
}
