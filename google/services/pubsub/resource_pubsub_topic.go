// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package pubsub

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourcePubsubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubTopicCreate,
		Read:   resourcePubsubTopicRead,
		Update: resourcePubsubTopicUpdate,
		Delete: resourcePubsubTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePubsubTopicImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `Name of the topic.`,
			},
			"kms_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The resource name of the Cloud KMS CryptoKey to be used to protect access
to messages published on this topic. Your project's PubSub service account
('service-{{PROJECT_NUMBER}}@gcp-sa-pubsub.iam.gserviceaccount.com') must have
'roles/cloudkms.cryptoKeyEncrypterDecrypter' to use this feature.
The expected format is 'projects/*/locations/*/keyRings/*/cryptoKeys/*'`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `A set of key/value label pairs to assign to this Topic.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"message_retention_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Indicates the minimum duration to retain a message after it is published
to the topic. If this field is set, messages published to the topic in
the last messageRetentionDuration are always available to subscribers.
For instance, it allows any attached subscription to seek to a timestamp
that is up to messageRetentionDuration in the past. If this field is not
set, message retention is controlled by settings on individual subscriptions.
Cannot be more than 31 days or less than 10 minutes.`,
			},
			"message_storage_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Description: `Policy constraining the set of Google Cloud Platform regions where
messages published to the topic may be stored. If not present, then no
constraints are in effect.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_persistence_regions": {
							Type:     schema.TypeList,
							Required: true,
							Description: `A list of IDs of GCP regions where messages that are published to
the topic may be persisted in storage. Messages published by
publishers running in non-allowed GCP regions (or running outside
of GCP altogether) will be routed for storage in one of the
allowed regions. An empty list means that no regions are allowed,
and is not a valid configuration.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"schema_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `Settings for validating messages published against a schema.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The name of the schema that messages published should be
validated against. Format is projects/{project}/schemas/{schema}.
The value of this field will be _deleted-schema_
if the schema has been deleted.`,
						},
						"encoding": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidateEnum([]string{"ENCODING_UNSPECIFIED", "JSON", "BINARY", ""}),
							Description:  `The encoding of messages validated against schema. Default value: "ENCODING_UNSPECIFIED" Possible values: ["ENCODING_UNSPECIFIED", "JSON", "BINARY"]`,
							Default:      "ENCODING_UNSPECIFIED",
						},
					},
				},
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourcePubsubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandPubsubTopicName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	kmsKeyNameProp, err := expandPubsubTopicKmsKeyName(d.Get("kms_key_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(kmsKeyNameProp)) && (ok || !reflect.DeepEqual(v, kmsKeyNameProp)) {
		obj["kmsKeyName"] = kmsKeyNameProp
	}
	messageStoragePolicyProp, err := expandPubsubTopicMessageStoragePolicy(d.Get("message_storage_policy"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("message_storage_policy"); !tpgresource.IsEmptyValue(reflect.ValueOf(messageStoragePolicyProp)) && (ok || !reflect.DeepEqual(v, messageStoragePolicyProp)) {
		obj["messageStoragePolicy"] = messageStoragePolicyProp
	}
	schemaSettingsProp, err := expandPubsubTopicSchemaSettings(d.Get("schema_settings"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("schema_settings"); !tpgresource.IsEmptyValue(reflect.ValueOf(schemaSettingsProp)) && (ok || !reflect.DeepEqual(v, schemaSettingsProp)) {
		obj["schemaSettings"] = schemaSettingsProp
	}
	messageRetentionDurationProp, err := expandPubsubTopicMessageRetentionDuration(d.Get("message_retention_duration"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("message_retention_duration"); !tpgresource.IsEmptyValue(reflect.ValueOf(messageRetentionDurationProp)) && (ok || !reflect.DeepEqual(v, messageRetentionDurationProp)) {
		obj["messageRetentionDuration"] = messageRetentionDurationProp
	}
	labelsProp, err := expandPubsubTopicEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	obj, err = resourcePubsubTopicEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Topic: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Topic: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "PUT",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutCreate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.PubsubTopicProjectNotReady},
	})
	if err != nil {
		return fmt.Errorf("Error creating Topic: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/topics/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = transport_tpg.PollingWaitTime(resourcePubsubTopicPollRead(d, meta), transport_tpg.PollCheckForExistence, "Creating Topic", d.Timeout(schema.TimeoutCreate), 1)
	if err != nil {
		log.Printf("[ERROR] Unable to confirm eventually consistent Topic %q finished updating: %q", d.Id(), err)
	}

	log.Printf("[DEBUG] Finished creating Topic %q: %#v", d.Id(), res)

	return resourcePubsubTopicRead(d, meta)
}

func resourcePubsubTopicPollRead(d *schema.ResourceData, meta interface{}) transport_tpg.PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*transport_tpg.Config)

		url, err := tpgresource.ReplaceVars(d, config, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
		if err != nil {
			return nil, err
		}

		billingProject := ""

		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return nil, fmt.Errorf("Error fetching project for Topic: %s", err)
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return nil, err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              billingProject,
			RawURL:               url,
			UserAgent:            userAgent,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.PubsubTopicProjectNotReady},
		})
		if err != nil {
			return res, err
		}
		return res, nil
	}
}

func resourcePubsubTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Topic: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.PubsubTopicProjectNotReady},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("PubsubTopic %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	if err := d.Set("name", flattenPubsubTopicName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("kms_key_name", flattenPubsubTopicKmsKeyName(res["kmsKeyName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("labels", flattenPubsubTopicLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("message_storage_policy", flattenPubsubTopicMessageStoragePolicy(res["messageStoragePolicy"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("schema_settings", flattenPubsubTopicSchemaSettings(res["schemaSettings"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("message_retention_duration", flattenPubsubTopicMessageRetentionDuration(res["messageRetentionDuration"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("terraform_labels", flattenPubsubTopicTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("effective_labels", flattenPubsubTopicEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	return nil
}

func resourcePubsubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Topic: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	kmsKeyNameProp, err := expandPubsubTopicKmsKeyName(d.Get("kms_key_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, kmsKeyNameProp)) {
		obj["kmsKeyName"] = kmsKeyNameProp
	}
	messageStoragePolicyProp, err := expandPubsubTopicMessageStoragePolicy(d.Get("message_storage_policy"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("message_storage_policy"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, messageStoragePolicyProp)) {
		obj["messageStoragePolicy"] = messageStoragePolicyProp
	}
	schemaSettingsProp, err := expandPubsubTopicSchemaSettings(d.Get("schema_settings"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("schema_settings"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, schemaSettingsProp)) {
		obj["schemaSettings"] = schemaSettingsProp
	}
	messageRetentionDurationProp, err := expandPubsubTopicMessageRetentionDuration(d.Get("message_retention_duration"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("message_retention_duration"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, messageRetentionDurationProp)) {
		obj["messageRetentionDuration"] = messageRetentionDurationProp
	}
	labelsProp, err := expandPubsubTopicEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	obj, err = resourcePubsubTopicUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Topic %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("kms_key_name") {
		updateMask = append(updateMask, "kmsKeyName")
	}

	if d.HasChange("message_storage_policy") {
		updateMask = append(updateMask, "messageStoragePolicy")
	}

	if d.HasChange("schema_settings") {
		updateMask = append(updateMask, "schemaSettings")
	}

	if d.HasChange("message_retention_duration") {
		updateMask = append(updateMask, "messageRetentionDuration")
	}

	if d.HasChange("effective_labels") {
		updateMask = append(updateMask, "labels")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "PATCH",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutUpdate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.PubsubTopicProjectNotReady},
	})

	if err != nil {
		return fmt.Errorf("Error updating Topic %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Topic %q: %#v", d.Id(), res)
	}

	return resourcePubsubTopicRead(d, meta)
}

func resourcePubsubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Topic: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubBasePath}}projects/{{project}}/topics/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Topic %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "DELETE",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutDelete),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.PubsubTopicProjectNotReady},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Topic")
	}

	log.Printf("[DEBUG] Finished deleting Topic %q: %#v", d.Id(), res)
	return nil
}

func resourcePubsubTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/topics/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/topics/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenPubsubTopicName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenPubsubTopicKmsKeyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubTopicLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenPubsubTopicMessageStoragePolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["allowed_persistence_regions"] =
		flattenPubsubTopicMessageStoragePolicyAllowedPersistenceRegions(original["allowedPersistenceRegions"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubTopicMessageStoragePolicyAllowedPersistenceRegions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubTopicSchemaSettings(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["schema"] =
		flattenPubsubTopicSchemaSettingsSchema(original["schema"], d, config)
	transformed["encoding"] =
		flattenPubsubTopicSchemaSettingsEncoding(original["encoding"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubTopicSchemaSettingsSchema(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubTopicSchemaSettingsEncoding(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubTopicMessageRetentionDuration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubTopicTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenPubsubTopicEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandPubsubTopicName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return tpgresource.GetResourceNameFromSelfLink(v.(string)), nil
}

func expandPubsubTopicKmsKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubTopicMessageStoragePolicy(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowedPersistenceRegions, err := expandPubsubTopicMessageStoragePolicyAllowedPersistenceRegions(original["allowed_persistence_regions"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowedPersistenceRegions); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowedPersistenceRegions"] = transformedAllowedPersistenceRegions
	}

	return transformed, nil
}

func expandPubsubTopicMessageStoragePolicyAllowedPersistenceRegions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubTopicSchemaSettings(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSchema, err := expandPubsubTopicSchemaSettingsSchema(original["schema"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSchema); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["schema"] = transformedSchema
	}

	transformedEncoding, err := expandPubsubTopicSchemaSettingsEncoding(original["encoding"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEncoding); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["encoding"] = transformedEncoding
	}

	return transformed, nil
}

func expandPubsubTopicSchemaSettingsSchema(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubTopicSchemaSettingsEncoding(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubTopicMessageRetentionDuration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubTopicEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func resourcePubsubTopicEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	delete(obj, "name")
	return obj, nil
}

func resourcePubsubTopicUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	newObj := make(map[string]interface{})
	newObj["topic"] = obj
	return newObj, nil
}
