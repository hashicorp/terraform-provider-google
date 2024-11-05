// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseAutonomousDatabases() *schema.Resource {
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
		"autonomous_databases": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseAutonomousDatabase().Schema),
			},
		},
	}
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseAutonomousDatabasesRead,
		Schema: dsSchema,
	}

}

func dataSourceOracleDatabaseAutonomousDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{OracleDatabaseBasePath}}projects/{{project}}/locations/{{location}}/autonomousDatabases")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	billingProject := ""
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for autonomousDatabases: %s", err)
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
		return fmt.Errorf("Error reading autonomousDatabases: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting autonomousDatabases project: %s", err)
	}

	if err := d.Set("autonomous_databases", flattenOracleDatabaseautonomousDatabases(res["autonomousDatabases"], d, config)); err != nil {
		return fmt.Errorf("Error setting autonomousDatabases: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/autonomousDatabases")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}

func flattenOracleDatabaseautonomousDatabases(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) []map[string]interface{} {
	if v == nil {
		return nil
	}
	l := v.([]interface{})
	transformed := make([]map[string]interface{}, 0)
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"name":             flattenOracleDatabaseAutonomousDatabaseName(original["name"], d, config),
			"database":         flattenOracleDatabaseAutonomousDatabaseDatabase(original["database"], d, config),
			"display_name":     flattenOracleDatabaseAutonomousDatabaseDisplayName(original["displayName"], d, config),
			"entitlement_id":   flattenOracleDatabaseAutonomousDatabaseEntitlementId(original["entitlementId"], d, config),
			"properties":       flattenOracleDatabaseAutonomousDatabaseProperties(original["properties"], d, config),
			"labels":           flattenOracleDatabaseAutonomousDatabaseLabels(original["labels"], d, config),
			"network":          flattenOracleDatabaseAutonomousDatabaseNetwork(original["network"], d, config),
			"cidr":             flattenOracleDatabaseAutonomousDatabaseCidr(original["cidr"], d, config),
			"create_time":      flattenOracleDatabaseAutonomousDatabaseCreateTime(original["createTime"], d, config),
			"terraform_labels": flattenOracleDatabaseAutonomousDatabaseTerraformLabels(original["labels"], d, config),
			"effective_labels": flattenOracleDatabaseAutonomousDatabaseEffectiveLabels(original["labels"], d, config),
		})
	}
	return transformed
}
