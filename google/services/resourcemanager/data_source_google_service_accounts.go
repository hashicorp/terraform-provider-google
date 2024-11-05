// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleServiceAccounts() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleServiceAccountsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleServiceAccountsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for service accounts: %s", err)
	}

	accounts := make([]map[string]interface{}, 0)

	accountList, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List("projects/" + project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Service accounts: %s", project))
	}

	for _, account := range accountList.Accounts {
		accounts = append(accounts, map[string]interface{}{
			"account_id":   strings.Split(account.Email, "@")[0],
			"disabled":     account.Disabled,
			"email":        account.Email,
			"display_name": account.DisplayName,
			"member":       "serviceAccount:" + account.Email,
			"name":         account.Name,
			"unique_id":    account.UniqueId,
		})
	}

	if err := d.Set("accounts", accounts); err != nil {
		return fmt.Errorf("Error retrieving service accounts: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s",
		project,
	))

	return nil
}
