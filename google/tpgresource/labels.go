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
			for k, _ := range v.(map[string]interface{}) {
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

func SetLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
	}

	raw := d.Get("labels")
	if raw == nil {
		return nil
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

func SetMetadataLabelsDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	l := d.Get("metadata").([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	raw := d.Get("metadata.0.labels")
	if raw == nil {
		return nil
	}

	config := meta.(*transport_tpg.Config)

	// Merge provider default labels with the user defined labels in the resource to get terraform managed labels
	terraformLabels := make(map[string]string)
	for k, v := range config.DefaultLabels {
		terraformLabels[k] = v
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
