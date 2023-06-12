// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigqueryDefaultServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleBigqueryDefaultServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleBigqueryDefaultServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	projectResource, err := config.NewBigQueryClient(userAgent).Projects.GetServiceAccount(project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "BigQuery service account not found")
	}

	d.SetId(projectResource.Email)
	if err := d.Set("email", projectResource.Email); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("member", "serviceAccount:"+projectResource.Email); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}
	return nil
}
