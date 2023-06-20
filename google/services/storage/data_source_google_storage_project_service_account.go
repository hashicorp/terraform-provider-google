// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageProjectServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"user_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_address": {
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

func dataSourceGoogleStorageProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	serviceAccountGetRequest := config.NewStorageClient(userAgent).Projects.ServiceAccount.Get(project)

	if v, ok := d.GetOk("user_project"); ok {
		serviceAccountGetRequest = serviceAccountGetRequest.UserProject(v.(string))
	}

	serviceAccount, err := serviceAccountGetRequest.Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "GCS service account not found")
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("email_address", serviceAccount.EmailAddress); err != nil {
		return fmt.Errorf("Error setting email_address: %s", err)
	}
	if err := d.Set("member", "serviceAccount:"+serviceAccount.EmailAddress); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}

	d.SetId(serviceAccount.EmailAddress)

	return nil
}
