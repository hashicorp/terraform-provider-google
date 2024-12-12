// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRBackup() *schema.Resource {
	dsSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `Name of resource`,
		},
		"backups": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: `List of all backups under data source.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Name of the resource.`,
					},
					"location": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Location of the resource.`,
					},
					"backup_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Id of the requesting object, Backup.`,
					},
					"backup_vault_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Name of the Backup Vault associated with Backup.`,
					},
					"data_source_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Name of the Data Source associated with Backup.`,
					},
				},
			},
		},
		"location": {
			Type:     schema.TypeString,
			Required: true,
		},
		"project": {
			Type:     schema.TypeString,
			Required: true,
		},
		"data_source_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"backup_vault_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDrBackupRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudBackupDrBackupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	billingProject := project
	url, err := tpgresource.ReplaceVars(d, config, "{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{backup_vault_id}}/dataSources/{{data_source_id}}/backups")

	fmt.Sprintf("url: %s", url)

	if err != nil {
		return err
	}

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
		return fmt.Errorf("Error reading BackupVault: %s", err)
	}

	if err := d.Set("backups", flattenDataSourceBackupDRBackups(res["backups"], d, config)); err != nil {
		return fmt.Errorf("Error reading Backup: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/backupVaults/{{backup_vault_id}}/dataSources/{{data_source_id}}/backups")
	d.SetId(id)
	d.Set("name", id)

	return nil
}

func flattenDataSourceBackupDRBackups(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))

	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"name":            flattenDataSourceBackupDRBackupsName(original["name"], d, config),
			"location":        flattenDataSourceBackupDRBackupsLocation(original["location"], d, config),
			"backup_id":       flattenDataSourceBackupDRBackupsBackupId(original["backupId"], d, config),
			"backup_vault_id": flattenDataSourceBackupDRBackupsBackupVaultId(original["backupVaultId"], d, config),
			"data_source_id":  flattenDataSourceBackupDRBackupsDataSourceId(original["dataSourceId"], d, config),
		})
	}
	return transformed
}

func flattenDataSourceBackupDRBackupsName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataSourceBackupDRBackupsLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataSourceBackupDRBackupsBackupId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataSourceBackupDRBackupsBackupVaultId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataSourceBackupDRBackupsDataSourceId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
