// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func SetTerraformLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
	}

	labels := d.Get("labels").(map[string]interface{})
	for k, v := range labels {
		terraformLabels[k] = v.(string)
	}

	if err := d.SetNew("terraform_labels", terraformLabels); err != nil {
		return fmt.Errorf("error setting new terraform_labels diff: %w", err)
	}

	return nil
}

func SetMetadataTerraformLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	v := d.Get("metadata")
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
	}

	labels := d.Get("metadata.0.labels").(map[string]interface{})
	for k, v := range labels {
		terraformLabels[k] = v.(string)
	}

	raw := l[0]
	original := raw.(map[string]interface{})
	original["terraform_labels"] = terraformLabels

	if err := d.SetNew("metadata", []interface{}{original}); err != nil {
		return fmt.Errorf("error setting new metadata diff: %w", err)
	}

	return nil
}
