// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseCloudVmClusters() *schema.Resource {
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
		"cloud_vm_clusters": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseCloudVmCluster().Schema),
			},
		},
	}
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseCloudVmClustersRead,
		Schema: dsSchema,
	}

}

func dataSourceOracleDatabaseCloudVmClustersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/cloudVmClusters")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/cloudVmClusters")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	billingProject := ""
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for cloudVmClusters: %s", err)
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
		return fmt.Errorf("Error reading cloudVmClusters: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting cloudVmClusters project: %s", err)
	}

	if err := d.Set("cloud_vm_clusters", flattenOracleDatabaseCloudVmClusters(res["cloudVmClusters"], d, config)); err != nil {
		return fmt.Errorf("Error setting cloudVmClusters: %s", err)
	}

	return nil
}

func flattenOracleDatabaseCloudVmClusters(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if v == nil {
		return nil
	}
	l := v.([]interface{})
	transformed := make([]map[string]interface{}, 0)
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"name":                   flattenOracleDatabaseCloudVmClusterName(original["name"], d, config),
			"exadata_infrastructure": flattenOracleDatabaseCloudVmClusterExadataInfrastructure(original["exadataInfrastructure"], d, config),
			"display_name":           flattenOracleDatabaseCloudVmClusterDisplayName(original["displayName"], d, config),
			"gcp_oracle_zone":        flattenOracleDatabaseCloudVmClusterGcpOracleZone(original["gcpOracleZone"], d, config),
			"properties":             flattenOracleDatabaseCloudVmClusterProperties(original["properties"], d, config),
			"labels":                 flattenOracleDatabaseCloudVmClusterLabels(original["labels"], d, config),
			"create_time":            flattenOracleDatabaseCloudVmClusterCreateTime(original["createTime"], d, config),
			"cidr":                   flattenOracleDatabaseCloudVmClusterCidr(original["cidr"], d, config),
			"backup_subnet_cidr":     flattenOracleDatabaseCloudVmClusterBackupSubnetCidr(original["backupSubnetCidr"], d, config),
			"network":                flattenOracleDatabaseCloudVmClusterNetwork(original["network"], d, config),
			"terraform_labels":       flattenOracleDatabaseCloudVmClusterTerraformLabels(original["labels"], d, config),
			"effective_labels":       flattenOracleDatabaseCloudVmClusterEffectiveLabels(original["labels"], d, config),
		})
	}
	return transformed
}
