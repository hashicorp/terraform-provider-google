// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseDbServers() *schema.Resource {
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
		"cloud_exadata_infrastructure": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "exadata",
		},
		"db_servers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"display_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The Display name",
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
								"max_ocpu_count": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"memory_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"max_memory_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"db_node_storage_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"max_db_node_storage_size_gb": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"vm_count": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Output only",
								},
								"state": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Output only",
								},
								"db_node_ids": {
									Type:        schema.TypeList,
									Computed:    true,
									Description: "Output only",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &schema.Resource{
		Read:          DataSourceOracleDatabaseDbServersRead,
		Schema:        dsSchema,
		UseJSONNumber: true,
	}
}

func DataSourceOracleDatabaseDbServersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DbServer: %s", err)
	}
	billingProject = project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures/{{cloud_exadata_infrastructure}}/dbServers")
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading DbServer: %s", err)
	}
	if res["dbServers"] == nil {
		return fmt.Errorf("Error reading DbServer: %s", err)
	}
	dbServers, err := flattenOracleDatabaseDbServerList(config, res["dbServers"])
	if err != nil {
		return fmt.Errorf("error flattening dbserver list: %s", err)
	}

	if err := d.Set("db_servers", dbServers); err != nil {
		return fmt.Errorf("error setting dbserver: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures/{{cloud_exadata_infrastructure_id}}/dbServers")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}

func flattenOracleDatabaseDbServerList(config *transport_tpg.Config, dbServerList interface{}) (interface{}, error) {

	if dbServerList == nil {
		return nil, nil
	}

	l := dbServerList.([]interface{})
	transformed := make([]interface{}, 0)
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"display_name": flattenOracleDatabaseDbServerDisplayName(original["displayName"], config),
			"properties":   flattenOracleDatabaseDbServerProperties(original["properties"], config),
		})
	}
	return transformed, nil

}

func flattenOracleDatabaseDbServerDisplayName(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerProperties(v interface{}, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["ocid"] = flattenOracleDatabaseDbServerPropertiesOcid(original["ocid"], config)
	transformed["ocpu_count"] = flattenOracleDatabaseDbServerPropertiesOcpuCount(original["ocpuCount"], config)
	transformed["max_ocpu_count"] = flattenOracleDatabaseDbServerPropertiesMaxOcpuCount(original["maxOcpuCount"], config)
	transformed["memory_size_gb"] = flattenOracleDatabaseDbServerPropertiesMemorySizeGb(original["memorySizeGb"], config)
	transformed["max_memory_size_gb"] = flattenOracleDatabaseDbServerPropertiesMaxMemorySizeGb(original["maxMemorySizeGb"], config)
	transformed["db_node_storage_size_gb"] = flattenOracleDatabaseDbServerPropertiesDbNodeStorageSizeGb(original["dbNodeStorageSizeGb"], config)
	transformed["max_db_node_storage_size_gb"] = flattenOracleDatabaseDbServerPropertiesMaxDbNodeStorageSizeGb(original["maxDbNodeStorageSizeGb"], config)
	transformed["vm_count"] = flattenOracleDatabaseDbServerPropertiesVmcount(original["vmCount"], config)
	transformed["state"] = flattenOracleDatabaseDbServerPropertiesState(original["state"], config)
	transformed["db_node_ids"] = flattenOracleDatabaseDbServerPropertiesDbNodeIds(original["dbNodeIds"], config)

	return []interface{}{transformed}
}

func flattenOracleDatabaseDbServerPropertiesOcid(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesOcpuCount(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesMaxOcpuCount(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesMemorySizeGb(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesMaxMemorySizeGb(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesDbNodeStorageSizeGb(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesMaxDbNodeStorageSizeGb(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesVmcount(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesState(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOracleDatabaseDbServerPropertiesDbNodeIds(v interface{}, config *transport_tpg.Config) interface{} {
	return v
}
