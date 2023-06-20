// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagetransfer

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageTransferProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageTransferProjectServiceAccountRead,
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
			"subject_id": {
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

func dataSourceGoogleStorageTransferProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	serviceAccount, err := config.NewStorageTransferClient(userAgent).GoogleServiceAccounts.Get(project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Google Cloud Storage Transfer service account not found")
	}

	d.SetId(serviceAccount.AccountEmail)
	if err := d.Set("email", serviceAccount.AccountEmail); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	if err := d.Set("subject_id", serviceAccount.SubjectId); err != nil {
		return fmt.Errorf("Error setting subject_id: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("member", "serviceAccount:"+serviceAccount.AccountEmail); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}
	return nil
}
