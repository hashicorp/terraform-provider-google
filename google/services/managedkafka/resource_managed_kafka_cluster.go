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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/managedkafka/Cluster.yaml
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

func ResourceManagedKafkaCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedKafkaClusterCreate,
		Read:   resourceManagedKafkaClusterRead,
		Update: resourceManagedKafkaClusterUpdate,
		Delete: resourceManagedKafkaClusterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceManagedKafkaClusterImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"capacity_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `A capacity configuration of a Kafka cluster.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory_bytes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The memory to provision for the cluster in bytes. The value must be between 1 GiB and 8 GiB per vCPU. Ex. 1024Mi, 4Gi.`,
						},
						"vcpu_count": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The number of vCPUs to provision for the cluster. The minimum is 3.`,
						},
					},
				},
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID to use for the cluster, which will become the final component of the cluster's name. The ID must be 1-63 characters long, and match the regular expression '[a-z]([-a-z0-9]*[a-z0-9])?' to comply with RFC 1035. This value is structured like: 'my-cluster-id'.`,
			},
			"gcp_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Configuration properties for a Kafka cluster deployed to Google Cloud Platform.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_config": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The configuration of access to the Kafka cluster.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_configs": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Virtual Private Cloud (VPC) subnets where IP addresses for the Kafka cluster are allocated. To make the cluster available in a VPC, you must specify at least one 'network_configs' block. Max of 10 subnets per cluster. Additional subnets may be specified with additional 'network_configs' blocks.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subnet": {
													Type:             schema.TypeString,
													Required:         true,
													DiffSuppressFunc: tpgresource.ProjectNumberDiffSuppress,
													Description:      `Name of the VPC subnet from which the cluster is accessible. Both broker and bootstrap server IP addresses and DNS entries are automatically created in the subnet. There can only be one subnet per network, and the subnet must be located in the same region as the cluster. The project may differ. The name of the subnet must be in the format 'projects/PROJECT_ID/regions/REGION/subnetworks/SUBNET'.`,
												},
											},
										},
									},
								},
							},
						},
						"kms_key": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.ProjectNumberDiffSuppress,
							Description:      `The Cloud KMS Key name to use for encryption. The key must be located in the same region as the cluster and cannot be changed. Must be in the format 'projects/PROJECT_ID/locations/LOCATION/keyRings/KEY_RING/cryptoKeys/KEY'.`,
						},
					},
				},
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the location of the Kafka resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `List of label KEY=VALUE pairs to add. Keys must start with a lowercase character and contain only hyphens (-), underscores ( ), lowercase characters, and numbers. Values must contain only hyphens (-), underscores ( ), lowercase characters, and numbers.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"rebalance_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Defines rebalancing behavior of a Kafka cluster.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The rebalance behavior for the cluster. When not specified, defaults to 'NO_REBALANCE'. Possible values: 'MODE_UNSPECIFIED', 'NO_REBALANCE', 'AUTO_REBALANCE_ON_SCALE_UP'.`,
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the cluster was created.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the cluster. Structured like: 'projects/PROJECT_ID/locations/LOCATION/clusters/CLUSTER_ID'.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current state of the cluster. Possible values: 'STATE_UNSPECIFIED', 'CREATING', 'ACTIVE', 'DELETING'.`,
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the cluster was last updated.`,
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

func resourceManagedKafkaClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	gcpConfigProp, err := expandManagedKafkaClusterGcpConfig(d.Get("gcp_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("gcp_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(gcpConfigProp)) && (ok || !reflect.DeepEqual(v, gcpConfigProp)) {
		obj["gcpConfig"] = gcpConfigProp
	}
	capacityConfigProp, err := expandManagedKafkaClusterCapacityConfig(d.Get("capacity_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("capacity_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(capacityConfigProp)) && (ok || !reflect.DeepEqual(v, capacityConfigProp)) {
		obj["capacityConfig"] = capacityConfigProp
	}
	rebalanceConfigProp, err := expandManagedKafkaClusterRebalanceConfig(d.Get("rebalance_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rebalance_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(rebalanceConfigProp)) && (ok || !reflect.DeepEqual(v, rebalanceConfigProp)) {
		obj["rebalanceConfig"] = rebalanceConfigProp
	}
	labelsProp, err := expandManagedKafkaClusterEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters?clusterId={{cluster_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Cluster: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Cluster: %s", err)
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
		return fmt.Errorf("Error creating Cluster: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = ManagedKafkaOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating Cluster", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create Cluster: %s", err)
	}

	if err := d.Set("name", flattenManagedKafkaClusterName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Cluster %q: %#v", d.Id(), res)

	return resourceManagedKafkaClusterRead(d, meta)
}

func resourceManagedKafkaClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Cluster: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ManagedKafkaCluster %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}

	if err := d.Set("gcp_config", flattenManagedKafkaClusterGcpConfig(res["gcpConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("name", flattenManagedKafkaClusterName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("create_time", flattenManagedKafkaClusterCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("update_time", flattenManagedKafkaClusterUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("labels", flattenManagedKafkaClusterLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("capacity_config", flattenManagedKafkaClusterCapacityConfig(res["capacityConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("rebalance_config", flattenManagedKafkaClusterRebalanceConfig(res["rebalanceConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("state", flattenManagedKafkaClusterState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("terraform_labels", flattenManagedKafkaClusterTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}
	if err := d.Set("effective_labels", flattenManagedKafkaClusterEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Cluster: %s", err)
	}

	return nil
}

func resourceManagedKafkaClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Cluster: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	gcpConfigProp, err := expandManagedKafkaClusterGcpConfig(d.Get("gcp_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("gcp_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, gcpConfigProp)) {
		obj["gcpConfig"] = gcpConfigProp
	}
	capacityConfigProp, err := expandManagedKafkaClusterCapacityConfig(d.Get("capacity_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("capacity_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, capacityConfigProp)) {
		obj["capacityConfig"] = capacityConfigProp
	}
	rebalanceConfigProp, err := expandManagedKafkaClusterRebalanceConfig(d.Get("rebalance_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rebalance_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, rebalanceConfigProp)) {
		obj["rebalanceConfig"] = rebalanceConfigProp
	}
	labelsProp, err := expandManagedKafkaClusterEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Cluster %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("gcp_config") {
		updateMask = append(updateMask, "gcpConfig")
	}

	if d.HasChange("capacity_config") {
		updateMask = append(updateMask, "capacityConfig")
	}

	if d.HasChange("rebalance_config") {
		updateMask = append(updateMask, "rebalanceConfig")
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
			return fmt.Errorf("Error updating Cluster %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Cluster %q: %#v", d.Id(), res)
		}

		err = ManagedKafkaOperationWaitTime(
			config, res, project, "Updating Cluster", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceManagedKafkaClusterRead(d, meta)
}

func resourceManagedKafkaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Cluster: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ManagedKafkaBasePath}}projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Cluster %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "Cluster")
	}

	err = ManagedKafkaOperationWaitTime(
		config, res, project, "Deleting Cluster", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Cluster %q: %#v", d.Id(), res)
	return nil
}

func resourceManagedKafkaClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/clusters/(?P<cluster_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cluster_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<cluster_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenManagedKafkaClusterGcpConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["access_config"] =
		flattenManagedKafkaClusterGcpConfigAccessConfig(original["accessConfig"], d, config)
	transformed["kms_key"] =
		flattenManagedKafkaClusterGcpConfigKmsKey(original["kmsKey"], d, config)
	return []interface{}{transformed}
}
func flattenManagedKafkaClusterGcpConfigAccessConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["network_configs"] =
		flattenManagedKafkaClusterGcpConfigAccessConfigNetworkConfigs(original["networkConfigs"], d, config)
	return []interface{}{transformed}
}
func flattenManagedKafkaClusterGcpConfigAccessConfigNetworkConfigs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"subnet": flattenManagedKafkaClusterGcpConfigAccessConfigNetworkConfigsSubnet(original["subnet"], d, config),
		})
	}
	return transformed
}
func flattenManagedKafkaClusterGcpConfigAccessConfigNetworkConfigsSubnet(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterGcpConfigKmsKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenManagedKafkaClusterCapacityConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["vcpu_count"] =
		flattenManagedKafkaClusterCapacityConfigVcpuCount(original["vcpuCount"], d, config)
	transformed["memory_bytes"] =
		flattenManagedKafkaClusterCapacityConfigMemoryBytes(original["memoryBytes"], d, config)
	return []interface{}{transformed}
}
func flattenManagedKafkaClusterCapacityConfigVcpuCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterCapacityConfigMemoryBytes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterRebalanceConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["mode"] =
		flattenManagedKafkaClusterRebalanceConfigMode(original["mode"], d, config)
	return []interface{}{transformed}
}
func flattenManagedKafkaClusterRebalanceConfigMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenManagedKafkaClusterTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenManagedKafkaClusterEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandManagedKafkaClusterGcpConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAccessConfig, err := expandManagedKafkaClusterGcpConfigAccessConfig(original["access_config"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccessConfig); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["accessConfig"] = transformedAccessConfig
	}

	transformedKmsKey, err := expandManagedKafkaClusterGcpConfigKmsKey(original["kms_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKmsKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["kmsKey"] = transformedKmsKey
	}

	return transformed, nil
}

func expandManagedKafkaClusterGcpConfigAccessConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedNetworkConfigs, err := expandManagedKafkaClusterGcpConfigAccessConfigNetworkConfigs(original["network_configs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetworkConfigs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["networkConfigs"] = transformedNetworkConfigs
	}

	return transformed, nil
}

func expandManagedKafkaClusterGcpConfigAccessConfigNetworkConfigs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedSubnet, err := expandManagedKafkaClusterGcpConfigAccessConfigNetworkConfigsSubnet(original["subnet"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubnet); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subnet"] = transformedSubnet
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandManagedKafkaClusterGcpConfigAccessConfigNetworkConfigsSubnet(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaClusterGcpConfigKmsKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaClusterCapacityConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedVcpuCount, err := expandManagedKafkaClusterCapacityConfigVcpuCount(original["vcpu_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedVcpuCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["vcpuCount"] = transformedVcpuCount
	}

	transformedMemoryBytes, err := expandManagedKafkaClusterCapacityConfigMemoryBytes(original["memory_bytes"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMemoryBytes); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["memoryBytes"] = transformedMemoryBytes
	}

	return transformed, nil
}

func expandManagedKafkaClusterCapacityConfigVcpuCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaClusterCapacityConfigMemoryBytes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaClusterRebalanceConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMode, err := expandManagedKafkaClusterRebalanceConfigMode(original["mode"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["mode"] = transformedMode
	}

	return transformed, nil
}

func expandManagedKafkaClusterRebalanceConfigMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandManagedKafkaClusterEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
