// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseCloudExadataInfrastructures() *schema.Resource {
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
		"cloud_exadata_infrastructures": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseCloudExadataInfrastructure().Schema),
			},
		},
	}
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseCloudExadataInfrastructuresRead,
		Schema: dsSchema,
	}

}

func dataSourceOracleDatabaseCloudExadataInfrastructuresRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	billingProject := ""
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for cloudExadataInfrastructures: %s", err)
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
		return fmt.Errorf("Error reading cloudExadataInfrastructures: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting cloudExadataInfrastructures project: %s", err)
	}

	if err := d.Set("cloud_exadata_infrastructures", flattenOracleDatabaseCloudExadataInfrastructures(res["cloudExadataInfrastructures"], d, config)); err != nil {
		return fmt.Errorf("Error setting cloudExadataInfrastructures: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/cloudExadataInfrastructures")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}

func flattenOracleDatabaseCloudExadataInfrastructures(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if v == nil {
		return nil
	}
	l := v.([]interface{})
	transformed := make([]map[string]interface{}, 0)
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"name":             flattenOracleDatabaseCloudExadataInfrastructureName(original["name"], d, config),
			"display_name":     flattenOracleDatabaseCloudExadataInfrastructureDisplayName(original["displayName"], d, config),
			"gcp_oracle_zone":  flattenOracleDatabaseCloudExadataInfrastructureGcpOracleZone(original["gcpOracleZone"], d, config),
			"entitlement_id":   flattenOracleDatabaseCloudExadataInfrastructureEntitlementId(original["entitlementId"], d, config),
			"properties":       flattenOracleDatabaseCloudExadataInfrastructureProperties(original["properties"], d, config),
			"labels":           flattenOracleDatabaseCloudExadataInfrastructureLabels(original["labels"], d, config),
			"create_time":      flattenOracleDatabaseCloudExadataInfrastructureCreateTime(original["createTime"], d, config),
			"terraform_labels": flattenOracleDatabaseCloudExadataInfrastructureTerraformLabels(original["labels"], d, config),
			"effective_labels": flattenOracleDatabaseCloudExadataInfrastructureEffectiveLabels(original["labels"], d, config),
		})
	}
	return transformed
}
