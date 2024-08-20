// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// SetLabels is called in the READ method of the resources to set
// the field "labels" and "terraform_labels" in the state based on the labels field in the configuration.
// So the field "labels" and "terraform_labels" in the state will only have the user defined labels.
// param "labels" is all of labels returned from API read reqeust.
// param "lineage" is the terraform lineage of the field and could be "labels" or "terraform_labels".
func SetLabels(labels map[string]string, d *schema.ResourceData, lineage string) error {
	transformed := make(map[string]interface{})

	if v, ok := d.GetOk(lineage); ok {
		if labels != nil {
			for k := range v.(map[string]interface{}) {
				transformed[k] = labels[k]
			}
		}
	}

	return d.Set(lineage, transformed)
}

// Sets the "labels" field and "terraform_labels" with the value of the field "effective_labels" for data sources.
// When reading data source, as the labels field is unavailable in the configuration of the data source,
// the "labels" field will be empty. With this funciton, the labels "field" will have all of labels in the resource.
func SetDataSourceLabels(d *schema.ResourceData) error {
	effectiveLabels := d.Get("effective_labels")
	if effectiveLabels == nil {
		return nil
	}

	if d.Get("labels") == nil {
		return fmt.Errorf("`labels` field is not present in the resource schema.")
	}
	if err := d.Set("labels", effectiveLabels); err != nil {
		return fmt.Errorf("Error setting labels in data source: %s", err)
	}

	if d.Get("terraform_labels") == nil {
		return fmt.Errorf("`terraform_labels` field is not present in the resource schema.")
	}
	if err := d.Set("terraform_labels", effectiveLabels); err != nil {
		return fmt.Errorf("Error setting terraform_labels in data source: %s", err)
	}

	return nil
}

// Sets the values of terraform_labels and effective_labels fields when labels field is in root level
func setLabelsFields(labelsField string, d *schema.ResourceDiff, meta interface{}, skipAttribution bool) error {
	raw := d.Get(labelsField)
	if raw == nil {
		return nil
	}

	if d.Get("terraform_labels") == nil {
		return fmt.Errorf("`terraform_labels` field is not present in the resource schema.")
	}

	if d.Get("effective_labels") == nil {
		return fmt.Errorf("`effective_labels` field is not present in the resource schema.")
	}

	// If "labels" field is computed, set "terraform_labels" and "effective_labels" to computed.
	// https://github.com/hashicorp/terraform-provider-google/issues/16217
	if !d.GetRawPlan().GetAttr(labelsField).IsWhollyKnown() {
		if err := d.SetNewComputed("terraform_labels"); err != nil {
			return fmt.Errorf("error setting terraform_labels to computed: %w", err)
		}

		if err := d.SetNewComputed("effective_labels"); err != nil {
			return fmt.Errorf("error setting effective_labels to computed: %w", err)
		}
		return nil
	}

	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
	}

	// Append optional label indicating the resource was provisioned using Terraform
	if !skipAttribution && config.AddTerraformAttributionLabel {
		if el, ok := d.Get("effective_labels").(map[string]any); ok {
			_, hasExistingLabel := el[transport_tpg.AttributionKey]
			if hasExistingLabel ||
				config.TerraformAttributionLabelAdditionStrategy == transport_tpg.ProactiveAttributionStrategy ||
				(config.TerraformAttributionLabelAdditionStrategy == transport_tpg.CreateOnlyAttributionStrategy && d.Id() == "") {
				terraformLabels[transport_tpg.AttributionKey] = transport_tpg.AttributionValue
			}
		}
	}

	labels := raw.(map[string]interface{})
	for k, v := range labels {
		terraformLabels[k] = v.(string)
	}

	if err := d.SetNew("terraform_labels", terraformLabels); err != nil {
		return fmt.Errorf("error setting new terraform_labels diff: %w", err)
	}

	o, n := d.GetChange("terraform_labels")
	effectiveLabels := d.Get("effective_labels").(map[string]interface{})

	for k, v := range n.(map[string]interface{}) {
		effectiveLabels[k] = v.(string)
	}

	for k := range o.(map[string]interface{}) {
		if _, ok := n.(map[string]interface{})[k]; !ok {
			delete(effectiveLabels, k)
		}
	}

	if err := d.SetNew("effective_labels", effectiveLabels); err != nil {
		return fmt.Errorf("error setting new effective_labels diff: %w", err)
	}

	return nil
}

func SetLabelsDiffWithoutAttributionLabel(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return setLabelsFields("labels", d, meta, true)
}

// The CustomizeDiff func to set the values of terraform_labels and effective_labels fields
// when labels field is at the root level and named "labels".
func SetLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return setLabelsFields("labels", d, meta, false)
}

// The CustomizeDiff func to set the values of terraform_labels and effective_labels fields
// when labels field is at the root level and has a diffent name (e.g. resource_labels) than "labels"
func SetDiffForLabelsWithCustomizedName(labelsField string) func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
		return setLabelsFields(labelsField, d, meta, false)
	}
}

func SetMetadataLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	l := d.Get("metadata").([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	// Fix the bug that the computed and nested "labels" field disappears from the terraform plan.
	// https://github.com/hashicorp/terraform-provider-google/issues/17756
	// The bug is introduced by SetNew on "metadata" field with the object including terraform_labels and effective_labels.
	// "terraform_labels" and "effective_labels" cannot be set directly due to a bug that SetNew doesn't work on nested fields
	// in terraform sdk.
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/459
	values := d.GetRawPlan().GetAttr("metadata").AsValueSlice()
	if len(values) > 0 && !values[0].GetAttr("labels").IsWhollyKnown() {
		return nil
	}

	raw := d.Get("metadata.0.labels")
	if raw == nil {
		return nil
	}

	if d.Get("metadata.0.terraform_labels") == nil {
		return fmt.Errorf("`metadata.0.terraform_labels` field is not present in the resource schema.")
	}

	if d.Get("metadata.0.effective_labels") == nil {
		return fmt.Errorf("`metadata.0.effective_labels` field is not present in the resource schema.")
	}

	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
	}

	// Append optional label indicating the resource was provisioned using Terraform
	if config.AddTerraformAttributionLabel {
		if el, ok := d.Get("metadata.0.effective_labels").(map[string]any); ok {
			_, hasExistingLabel := el[transport_tpg.AttributionKey]
			if hasExistingLabel ||
				config.TerraformAttributionLabelAdditionStrategy == transport_tpg.ProactiveAttributionStrategy ||
				(config.TerraformAttributionLabelAdditionStrategy == transport_tpg.CreateOnlyAttributionStrategy && d.Id() == "") {
				terraformLabels[transport_tpg.AttributionKey] = transport_tpg.AttributionValue
			}
		}
	}

	labels := raw.(map[string]interface{})
	for k, v := range labels {
		terraformLabels[k] = v.(string)
	}

	original := l[0].(map[string]interface{})

	original["terraform_labels"] = terraformLabels
	if err := d.SetNew("metadata", []interface{}{original}); err != nil {
		return fmt.Errorf("error setting new metadata diff: %w", err)
	}

	o, n := d.GetChange("metadata.0.terraform_labels")
	effectiveLabels := d.Get("metadata.0.effective_labels").(map[string]interface{})

	for k, v := range n.(map[string]interface{}) {
		effectiveLabels[k] = v.(string)
	}

	for k := range o.(map[string]interface{}) {
		if _, ok := n.(map[string]interface{})[k]; !ok {
			delete(effectiveLabels, k)
		}
	}

	original["effective_labels"] = effectiveLabels
	if err := d.SetNew("metadata", []interface{}{original}); err != nil {
		return fmt.Errorf("error setting new metadata diff: %w", err)
	}

	return nil
}

// Upgrade the field "labels" in the state to exclude the labels with the labels prefix
// and the field "effective_labels" to have all of labels, including the labels with the labels prefix
func LabelsStateUpgrade(rawState map[string]interface{}, labesPrefix string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)
	log.Printf("[DEBUG] Attributes before migration labels: %#v", rawState["labels"])
	log.Printf("[DEBUG] Attributes before migration effective_labels: %#v", rawState["effective_labels"])

	if rawState["labels"] != nil {
		rawLabels := rawState["labels"].(map[string]interface{})
		labels := make(map[string]interface{})
		effectiveLabels := make(map[string]interface{})

		for k, v := range rawLabels {
			effectiveLabels[k] = v

			if !strings.HasPrefix(k, labesPrefix) {
				labels[k] = v
			}
		}

		rawState["labels"] = labels
		rawState["effective_labels"] = effectiveLabels
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	log.Printf("[DEBUG] Attributes after migration labels: %#v", rawState["labels"])
	log.Printf("[DEBUG] Attributes after migration effective_labels: %#v", rawState["effective_labels"])

	return rawState, nil
}

// Upgrade the field "terraform_labels" in the state to have the value of filed "labels"
// when it is not set but "labels" field is set in the state
func TerraformLabelsStateUpgrade(rawState map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)
	log.Printf("[DEBUG] Attributes before migration terraform_labels: %#v", rawState["terraform_labels"])

	if rawState["terraform_labels"] == nil && rawState["labels"] != nil {
		rawState["terraform_labels"] = rawState["labels"]
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	log.Printf("[DEBUG] Attributes after migration terraform_labels: %#v", rawState["terraform_labels"])

	return rawState, nil
}
