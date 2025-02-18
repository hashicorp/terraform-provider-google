// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/pubsublite/Topic.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package pubsublite

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourcePubsubLiteTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubLiteTopicCreate,
		Read:   resourcePubsubLiteTopicRead,
		Update: resourcePubsubLiteTopicUpdate,
		Delete: resourcePubsubLiteTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePubsubLiteTopicImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
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
			"partition_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The settings for this topic's partitions.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The number of partitions in the topic. Must be at least 1.`,
						},
						"capacity": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The capacity configuration.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publish_mib_per_sec": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `Subscribe throughput capacity per partition in MiB/s. Must be >= 4 and <= 16.`,
									},
									"subscribe_mib_per_sec": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `Publish throughput capacity per partition in MiB/s. Must be >= 4 and <= 16.`,
									},
								},
							},
						},
					},
				},
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region of the pubsub lite topic.`,
			},
			"reservation_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The settings for this topic's Reservation usage.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"throughput_reservation": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The Reservation to use for this topic's throughput capacity.`,
						},
					},
				},
			},
			"retention_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The settings for a topic's message retention.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"per_partition_bytes": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The provisioned storage, in bytes, per partition. If the number of bytes stored
in any of the topic's partitions grows beyond this value, older messages will be
dropped to make room for newer ones, regardless of the value of period.`,
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `How long a published message is retained. If unset, messages will be retained as
long as the bytes retained for each partition is below perPartitionBytes. A
duration in seconds with up to nine fractional digits, terminated by 's'.
Example: "3.5s".`,
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The zone of the pubsub lite topic.`,
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

func resourcePubsubLiteTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	partitionConfigProp, err := expandPubsubLiteTopicPartitionConfig(d.Get("partition_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("partition_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(partitionConfigProp)) && (ok || !reflect.DeepEqual(v, partitionConfigProp)) {
		obj["partitionConfig"] = partitionConfigProp
	}
	retentionConfigProp, err := expandPubsubLiteTopicRetentionConfig(d.Get("retention_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("retention_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(retentionConfigProp)) && (ok || !reflect.DeepEqual(v, retentionConfigProp)) {
		obj["retentionConfig"] = retentionConfigProp
	}
	reservationConfigProp, err := expandPubsubLiteTopicReservationConfig(d.Get("reservation_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("reservation_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(reservationConfigProp)) && (ok || !reflect.DeepEqual(v, reservationConfigProp)) {
		obj["reservationConfig"] = reservationConfigProp
	}

	obj, err = resourcePubsubLiteTopicEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubLiteBasePath}}projects/{{project}}/locations/{{zone}}/topics?topicId={{name}}")
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

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Topic: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{zone}}/topics/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Topic %q: %#v", d.Id(), res)

	return resourcePubsubLiteTopicRead(d, meta)
}

func resourcePubsubLiteTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubLiteBasePath}}projects/{{project}}/locations/{{zone}}/topics/{{name}}")
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

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("PubsubLiteTopic %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	if err := d.Set("partition_config", flattenPubsubLiteTopicPartitionConfig(res["partitionConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("retention_config", flattenPubsubLiteTopicRetentionConfig(res["retentionConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("reservation_config", flattenPubsubLiteTopicReservationConfig(res["reservationConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	return nil
}

func resourcePubsubLiteTopicUpdate(d *schema.ResourceData, meta interface{}) error {
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
	partitionConfigProp, err := expandPubsubLiteTopicPartitionConfig(d.Get("partition_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("partition_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, partitionConfigProp)) {
		obj["partitionConfig"] = partitionConfigProp
	}
	retentionConfigProp, err := expandPubsubLiteTopicRetentionConfig(d.Get("retention_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("retention_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, retentionConfigProp)) {
		obj["retentionConfig"] = retentionConfigProp
	}
	reservationConfigProp, err := expandPubsubLiteTopicReservationConfig(d.Get("reservation_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("reservation_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, reservationConfigProp)) {
		obj["reservationConfig"] = reservationConfigProp
	}

	obj, err = resourcePubsubLiteTopicEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubLiteBasePath}}projects/{{project}}/locations/{{zone}}/topics/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Topic %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("partition_config") {
		updateMask = append(updateMask, "partitionConfig")
	}

	if d.HasChange("retention_config") {
		updateMask = append(updateMask, "retentionConfig")
	}

	if d.HasChange("reservation_config") {
		updateMask = append(updateMask, "reservationConfig")
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

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error updating Topic %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Topic %q: %#v", d.Id(), res)
		}

	}

	return resourcePubsubLiteTopicRead(d, meta)
}

func resourcePubsubLiteTopicDelete(d *schema.ResourceData, meta interface{}) error {
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

	url, err := tpgresource.ReplaceVars(d, config, "{{PubsubLiteBasePath}}projects/{{project}}/locations/{{zone}}/topics/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Topic %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Topic")
	}

	log.Printf("[DEBUG] Finished deleting Topic %q: %#v", d.Id(), res)
	return nil
}

func resourcePubsubLiteTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<zone>[^/]+)/topics/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<zone>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{zone}}/topics/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenPubsubLiteTopicPartitionConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["count"] =
		flattenPubsubLiteTopicPartitionConfigCount(original["count"], d, config)
	transformed["capacity"] =
		flattenPubsubLiteTopicPartitionConfigCapacity(original["capacity"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubLiteTopicPartitionConfigCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenPubsubLiteTopicPartitionConfigCapacity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["publish_mib_per_sec"] =
		flattenPubsubLiteTopicPartitionConfigCapacityPublishMibPerSec(original["publishMibPerSec"], d, config)
	transformed["subscribe_mib_per_sec"] =
		flattenPubsubLiteTopicPartitionConfigCapacitySubscribeMibPerSec(original["subscribeMibPerSec"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubLiteTopicPartitionConfigCapacityPublishMibPerSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenPubsubLiteTopicPartitionConfigCapacitySubscribeMibPerSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenPubsubLiteTopicRetentionConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["per_partition_bytes"] =
		flattenPubsubLiteTopicRetentionConfigPerPartitionBytes(original["perPartitionBytes"], d, config)
	transformed["period"] =
		flattenPubsubLiteTopicRetentionConfigPeriod(original["period"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubLiteTopicRetentionConfigPerPartitionBytes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubLiteTopicRetentionConfigPeriod(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenPubsubLiteTopicReservationConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["throughput_reservation"] =
		flattenPubsubLiteTopicReservationConfigThroughputReservation(original["throughputReservation"], d, config)
	return []interface{}{transformed}
}
func flattenPubsubLiteTopicReservationConfigThroughputReservation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func expandPubsubLiteTopicPartitionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCount, err := expandPubsubLiteTopicPartitionConfigCount(original["count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["count"] = transformedCount
	}

	transformedCapacity, err := expandPubsubLiteTopicPartitionConfigCapacity(original["capacity"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCapacity); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["capacity"] = transformedCapacity
	}

	return transformed, nil
}

func expandPubsubLiteTopicPartitionConfigCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubLiteTopicPartitionConfigCapacity(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPublishMibPerSec, err := expandPubsubLiteTopicPartitionConfigCapacityPublishMibPerSec(original["publish_mib_per_sec"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublishMibPerSec); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["publishMibPerSec"] = transformedPublishMibPerSec
	}

	transformedSubscribeMibPerSec, err := expandPubsubLiteTopicPartitionConfigCapacitySubscribeMibPerSec(original["subscribe_mib_per_sec"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubscribeMibPerSec); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["subscribeMibPerSec"] = transformedSubscribeMibPerSec
	}

	return transformed, nil
}

func expandPubsubLiteTopicPartitionConfigCapacityPublishMibPerSec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubLiteTopicPartitionConfigCapacitySubscribeMibPerSec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubLiteTopicRetentionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPerPartitionBytes, err := expandPubsubLiteTopicRetentionConfigPerPartitionBytes(original["per_partition_bytes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPerPartitionBytes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["perPartitionBytes"] = transformedPerPartitionBytes
	}

	transformedPeriod, err := expandPubsubLiteTopicRetentionConfigPeriod(original["period"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPeriod); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["period"] = transformedPeriod
	}

	return transformed, nil
}

func expandPubsubLiteTopicRetentionConfigPerPartitionBytes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubLiteTopicRetentionConfigPeriod(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandPubsubLiteTopicReservationConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedThroughputReservation, err := expandPubsubLiteTopicReservationConfigThroughputReservation(original["throughput_reservation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedThroughputReservation); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["throughputReservation"] = transformedThroughputReservation
	}

	return transformed, nil
}

func expandPubsubLiteTopicReservationConfigThroughputReservation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	f, err := tpgresource.ParseRegionalFieldValue("reservations", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for throughput_reservation: %s", err)
	}
	// Custom due to "locations" rather than "regions".
	return fmt.Sprintf("projects/%s/locations/%s/reservations/%s", f.Project, f.Region, f.Name), nil
}

func resourcePubsubLiteTopicEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*transport_tpg.Config)

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return nil, err
	}

	if zone == "" {
		return nil, fmt.Errorf("zone must be non-empty - set in resource or at provider-level")
	}

	// API Endpoint requires region in the URL. We infer it from the zone.

	region := tpgresource.GetRegionFromZone(zone)

	if region == "" {
		return nil, fmt.Errorf("invalid zone %q, unable to infer region from zone", zone)
	}

	return obj, nil
}
