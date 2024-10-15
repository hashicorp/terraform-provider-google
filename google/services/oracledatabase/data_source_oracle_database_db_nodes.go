// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseDbNodes() *schema.Resource {
	dsSchema := map[string]*schema.Schema{
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the project in which the dataset is located. If it is not provided, the provider project is used.",
		},
		"location": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "location",
		},
		"cloud_vm_cluster": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "vmcluster",
		},
		"db_nodes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The dbnode name",
					},
					"properties": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ocid": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Output only",
								},
								"ocpu_count": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"memory_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"db_node_storage_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"db_server_ocid": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Output only",
								},
								"hostname": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Output only",
								},
								"state": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Output only",
								},
								"total_cpu_core_count": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
							},
						},
					},
				},
			},
		},
	}
	return &schema.Resource{
		Read:   DataSourceOracleDatabaseDbNodesRead,
		Schema: dsSchema,
	}
}

func DataSourceOracleDatabaseDbNodesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	url, err := tpgresource.ReplaceVars(d, config, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/cloudVmClusters/{{cloud_vm_cluster}}/dbNodes")
	if err != nil {
		return err
	}
	billingProject := ""
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DbNode: %s", err)
	}
	billingProject = project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return fmt.Errorf("Error reading DbNode: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading DbNode: %s", err)
	}
	if err := d.Set("db_nodes", flattenOracleDatabaseDbNodes(res["dbNodes"], d, config)); err != nil {
		return fmt.Errorf("Error reading DbNode: %s", err)
	}
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/cloudVmClusters/{{cloud_vm_cluster}}/dbNodes")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}

func flattenOracleDatabaseDbNodes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if v == nil {
		return nil
	}
	l := v.([]interface{})
	transformed := make([]map[string]interface{}, 0)
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"name":       flattenOracleDatabaseDbNodeName(original["name"], d, config),
			"properties": flattenOracleDatabaseDbNodeProperties(original["properties"], d, config),
		})
	}

	return transformed
}

func flattenOracleDatabaseDbNodeName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodeProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["ocid"] = flattenOracleDatabaseDbNodePropertiesOcid(original["ocid"], d, config)
	transformed["ocpu_count"] = flattenOracleDatabaseDbNodePropertiesOcpuCount(original["ocpuCount"], d, config)
	transformed["memory_size_gb"] = flattenOracleDatabaseDbNodePropertiesMemorySizeGb(original["memorySizeGb"], d, config)
	transformed["db_node_storage_size_gb"] = flattenOracleDatabaseDbNodePropertiesDbNodeStorageSizeGb(original["dbNodeStorageSizeGb"], d, config)
	transformed["db_server_ocid"] = flattenOracleDatabaseDbNodePropertiesDbServerOcid(original["dbServerOcid"], d, config)
	transformed["hostname"] = flattenOracleDatabaseDbNodePropertiesHostname(original["hostname"], d, config)
	transformed["state"] = flattenOracleDatabaseDbNodePropertiesState(original["state"], d, config)
	transformed["total_cpu_core_count"] = flattenOracleDatabaseDbNodePropertiesTotalCpuCoreCount(original["totalCpuCoreCount"], d, config)

	return []interface{}{transformed}
}

func flattenOracleDatabaseDbNodePropertiesOcid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesOcpuCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesMemorySizeGb(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesDbNodeStorageSizeGb(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesDbServerOcid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesHostname(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbNodePropertiesTotalCpuCoreCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
