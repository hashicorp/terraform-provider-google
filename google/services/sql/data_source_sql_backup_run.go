// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func DataSourceSqlBackupRun() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSqlBackupRunRead,

		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The identifier for this backup run. Unique only for a specific Cloud SQL instance. If left empty and multiple backups exist for the instance, most_recent must be set to true.`,
			},
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of the database instance.`,
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Location of the backups.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Project ID of the project that contains the instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time the backup operation actually started in UTC timezone in RFC 3339 format, for example 2012-11-15T16:19:00.094Z.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of this run.`,
			},
			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Toggles use of the most recent backup run if multiple backups exist for a Cloud SQL instance.`,
			},
		},
	}
}

func dataSourceSqlBackupRunRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)

	var backup *sqladmin.BackupRun
	if backupId, ok := d.GetOk("backup_id"); ok {
		backup, err = config.NewSqlAdminClient(userAgent).BackupRuns.Get(project, instance, int64(backupId.(int))).Do()
		if err != nil {
			return err
		}
	} else {
		res, err := config.NewSqlAdminClient(userAgent).BackupRuns.List(project, instance).Do()
		if err != nil {
			return err
		}
		backupsList := res.Items
		if len(backupsList) == 0 {
			return fmt.Errorf("No backups found for SQL Database Instance %s", instance)
		} else if len(backupsList) > 1 {
			mostRecent := d.Get("most_recent").(bool)
			if !mostRecent {
				return fmt.Errorf("Multiple SQL backup runs listed for Instance %s. Consider setting most_recent or specifying a backup_id", instance)
			}
		}
		backup = backupsList[0]
	}

	if err := d.Set("backup_id", backup.Id); err != nil {
		return fmt.Errorf("Error setting backup_id: %s", err)
	}
	if err := d.Set("location", backup.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("start_time", backup.StartTime); err != nil {
		return fmt.Errorf("Error setting start_time: %s", err)
	}
	if err := d.Set("status", backup.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/backupRuns/{{backup_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}
