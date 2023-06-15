// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleAppEngineDefaultServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleAppEngineDefaultServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleAppEngineDefaultServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	serviceAccountEmail := fmt.Sprintf("%s@appspot.gserviceaccount.com", project)

	serviceAccountName, err := tpgresource.ServiceAccountFQN(serviceAccountEmail, d, config)
	if err != nil {
		return err
	}

	sa, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Get(serviceAccountName).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Service Account %q", serviceAccountName))
	}

	d.SetId(sa.Name)
	if err := d.Set("email", sa.Email); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	if err := d.Set("unique_id", sa.UniqueId); err != nil {
		return fmt.Errorf("Error setting unique_id: %s", err)
	}
	if err := d.Set("project", sa.ProjectId); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("name", sa.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("display_name", sa.DisplayName); err != nil {
		return fmt.Errorf("Error setting display_name: %s", err)
	}
	if err := d.Set("member", "serviceAccount:"+sa.Email); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}

	return nil
}
