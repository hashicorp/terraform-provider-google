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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/managedkafka/Topic.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package managedkafka

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

func ResourceManagedKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedKafkaTopicCreate,
		Read:   resourceManagedKafkaTopicRead,
		Update: resourceManagedKafkaTopicUpdate,
		Delete: resourceManagedKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourceManagedKafkaTopicImport,
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
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The cluster name.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the location of the Kafka resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations.`,
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The number of replicas of each partition. A replication factor of 3 is recommended for high availability.`,
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID to use for the topic, which will become the final component of the topic's name. This value is structured like: 'my-topic-name'.`,
			},
			"configs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Configuration for the topic that are overridden from the cluster defaults. The key of the map is a Kafka topic property name, for example: 'cleanup.policy=compact', 'compression.type=producer'.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"partition_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of partitions in a topic. You can increase the partition count for a topic, but you cannot decrease it. Increasing partitions for a topic that uses a key might change how messages are distributed.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the topic. The 'topic' segment is used when connecting directly to the cluster. Must be in the format 'projects/PROJECT_ID/locations/LOCATION/clusters/CLUSTER_ID/topics/TOPIC_ID'.`,
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

func resourceManagedKafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	partitionCountProp, err := expandManagedKafkaTopicPartitionCount(d.Get("partition_count"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("partition_count"); !tpgresource.IsEmptyValue(reflect.ValueOf(partitionCountProp)) && (ok || !reflect.DeepEqual(v, partitionCountProp)) {
		obj["partitionCount"] = partitionCountProp
	}
	replicationFactorProp, err := expandManagedKafkaTopicReplicationFactor(d.Get("replication_factor"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("replication_factor"); !tpgresource.IsEmptyValue(reflect.ValueOf(replicationFactorProp)) && (ok || !reflect.DeepEqual(v, replicationFactorProp)) {
		obj["replicationFactor"] = replicationFactorProp
	}
	configsProp, err := expandManagedKafkaTopicConfigs(d.Get("configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(configsProp)) && (ok || !reflect.DeepEqual(v, configsProp)) {
		obj["configs"] = configsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics?topicId={{topic_id}}")
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
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics/{{topic_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)

	log.Printf("[DEBUG] Finished creating Topic %q: %#v", d.Id(), res)

	return resourceManagedKafkaTopicRead(d, meta)
}

func resourceManagedKafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics/{{topic_id}}")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ManagedKafkaTopic %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	if err := d.Set("name", flattenManagedKafkaTopicName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("partition_count", flattenManagedKafkaTopicPartitionCount(res["partitionCount"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("replication_factor", flattenManagedKafkaTopicReplicationFactor(res["replicationFactor"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}
	if err := d.Set("configs", flattenManagedKafkaTopicConfigs(res["configs"], d, config)); err != nil {
		return fmt.Errorf("Error reading Topic: %s", err)
	}

	return nil
}

func resourceManagedKafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
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
	partitionCountProp, err := expandManagedKafkaTopicPartitionCount(d.Get("partition_count"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("partition_count"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, partitionCountProp)) {
		obj["partitionCount"] = partitionCountProp
	}
	configsProp, err := expandManagedKafkaTopicConfigs(d.Get("configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("configs"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, configsProp)) {
		obj["configs"] = configsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics/{{topic_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Topic %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("partition_count") {
		updateMask = append(updateMask, "partitionCount")
	}

	if d.HasChange("configs") {
		updateMask = append(updateMask, "configs")
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

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)
	return resourceManagedKafkaTopicRead(d, meta)
}

func resourceManagedKafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
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

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics/{{topic_id}}")
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

func resourceManagedKafkaTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/clusters/(?P<cluster>[^/]+)/topics/(?P<topic_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<topic_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<topic_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/topics/{{topic_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenManagedKafkaTopicName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaTopicPartitionCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenManagedKafkaTopicReplicationFactor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenManagedKafkaTopicConfigs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandManagedKafkaTopicPartitionCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaTopicReplicationFactor(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaTopicConfigs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
