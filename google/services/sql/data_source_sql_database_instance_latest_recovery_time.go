// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSqlDatabaseInstanceLatestRecoveryTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSqlDatabaseInstanceLatestRecoveryTimeRead,

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"latest_recovery_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Timestamp, identifies the latest recovery time of the source instance.`,
			},
		},
	}
}

func dataSourceSqlDatabaseInstanceLatestRecoveryTimeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	fv, err := tpgresource.ParseProjectFieldValue("instances", d.Get("instance").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	project := fv.Project
	instance := fv.Name

	latestRecoveryTime, err := config.NewSqlAdminClient(userAgent).Projects.Instances.GetLatestRecoveryTime(project, instance).Do()
	if err != nil {
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := d.Set("latest_recovery_time", latestRecoveryTime.LatestRecoveryTime); err != nil {
		return fmt.Errorf("Error setting latest_recovery_time: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/instance/%s", project, instance))
	return nil
}
