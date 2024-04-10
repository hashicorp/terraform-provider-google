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

	if d.Get("effective_annotations") == nil {
		return fmt.Errorf("`effective_annotations` field is not present in the resource schema.")
	}

	// If "annotations" field is computed, set "effective_annotations" to computed.
	// https://github.com/hashicorp/terraform-provider-google/issues/16217
	if !d.GetRawPlan().GetAttr("annotations").IsWhollyKnown() {
		if err := d.SetNewComputed("effective_annotations"); err != nil {
			return fmt.Errorf("error setting effective_annotations to computed: %w", err)
		}
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

	// Fix the bug that the computed and nested "annotations" field disappears from the terraform plan.
	// https://github.com/hashicorp/terraform-provider-google/issues/17756
	// The bug is introduced by SetNew on "metadata" field with the object including "effective_annotations".
	// "effective_annotations" cannot be set directly due to a bug that SetNew doesn't work on nested fields
	// in terraform sdk.
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/459
	if !d.GetRawPlan().GetAttr("metadata").AsValueSlice()[0].GetAttr("annotations").IsWhollyKnown() {
		return nil
	}

	raw := d.Get("metadata.0.annotations")
	if raw == nil {
		return nil
	}

	if d.Get("metadata.0.effective_annotations") == nil {
		return fmt.Errorf("`metadata.0.effective_annotations` field is not present in the resource schema.")
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

// Sets the "annotations" field with the value of the field "effective_annotations" for data sources.
// When reading data source, as the annotations field is unavailable in the configuration of the data source,
// the "annotations" field will be empty. With this funciton, the labels "annotations" will have all of annotations in the resource.
func SetDataSourceAnnotations(d *schema.ResourceData) error {
	effectiveAnnotations := d.Get("effective_annotations")
	if effectiveAnnotations == nil {
		return nil
	}

	if d.Get("annotations") == nil {
		return fmt.Errorf("`annotations` field is not present in the resource schema.")
	}
	if err := d.Set("annotations", effectiveAnnotations); err != nil {
		return fmt.Errorf("Error setting annotations in data source: %s", err)
	}

	return nil
}
