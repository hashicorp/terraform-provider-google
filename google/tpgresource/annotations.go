// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SetAnnotationsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	raw := d.Get("annotations")
	if raw == nil {
		return nil
	}

	o, n := d.GetChange("annotations")
	effectiveAnnotations := d.Get("effective_annotations").(map[string]interface{})

	for k, v := range n.(map[string]interface{}) {
		effectiveAnnotations[k] = v.(string)
	}

	for k := range o.(map[string]interface{}) {
		if _, ok := n.(map[string]interface{})[k]; !ok {
			delete(effectiveAnnotations, k)
		}
	}

	if err := d.SetNew("effective_annotations", effectiveAnnotations); err != nil {
		return fmt.Errorf("error setting new effective_annotations diff: %w", err)
	}

	return nil
}

func SetMetadataAnnotationsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	l := d.Get("metadata").([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	raw := d.Get("metadata.0.annotations")
	if raw == nil {
		return nil
	}

	o, n := d.GetChange("metadata.0.annotations")
	effectiveAnnotations := d.Get("metadata.0.effective_annotations").(map[string]interface{})

	for k, v := range n.(map[string]interface{}) {
		effectiveAnnotations[k] = v.(string)
	}

	for k := range o.(map[string]interface{}) {
		if _, ok := n.(map[string]interface{})[k]; !ok {
			delete(effectiveAnnotations, k)
		}
	}

	original := l[0].(map[string]interface{})
	original["effective_annotations"] = effectiveAnnotations

	if err := d.SetNew("metadata", []interface{}{original}); err != nil {
		return fmt.Errorf("error setting new metadata diff: %w", err)
	}

	return nil
}
