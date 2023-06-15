// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package billing

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/cloudbilling/v1"
)

func DataSourceGoogleBillingAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBillingAccountRead,
		Schema: map[string]*schema.Schema{
			"billing_account": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"display_name"},
			},
			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"billing_account"},
			},
			"open": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lookup_projects": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceBillingAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	open, openOk := d.GetOkExists("open")

	var billingAccount *cloudbilling.BillingAccount
	if v, ok := d.GetOk("billing_account"); ok {
		resp, err := config.NewBillingClient(userAgent).BillingAccounts.Get(CanonicalBillingAccountName(v.(string))).Do()
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Billing Account Not Found : %s", v))
		}

		if openOk && resp.Open != open.(bool) {
			return fmt.Errorf("Billing account not found: %s", v)
		}

		billingAccount = resp
	} else if v, ok := d.GetOk("display_name"); ok {
		token := ""
		for paginate := true; paginate; {
			resp, err := config.NewBillingClient(userAgent).BillingAccounts.List().PageToken(token).Do()
			if err != nil {
				return fmt.Errorf("Error reading billing accounts: %s", err)
			}

			for _, ba := range resp.BillingAccounts {
				if ba.DisplayName == v.(string) {
					if openOk && ba.Open != open.(bool) {
						continue
					}
					if billingAccount != nil {
						return fmt.Errorf("More than one matching billing account found")
					}
					billingAccount = ba
				}
			}

			token = resp.NextPageToken
			paginate = token != ""
		}

		if billingAccount == nil {
			return fmt.Errorf("Billing account not found: %s", v)
		}
	} else {
		return fmt.Errorf("one of billing_account or display_name must be set")
	}

	if d.Get("lookup_projects").(bool) {
		resp, err := config.NewBillingClient(userAgent).BillingAccounts.Projects.List(billingAccount.Name).Do()
		if err != nil {
			return fmt.Errorf("Error reading billing account projects: %s", err)
		}
		projectIds := flattenBillingProjects(resp.ProjectBillingInfo)

		if err := d.Set("project_ids", projectIds); err != nil {
			return fmt.Errorf("Error setting project_ids: %s", err)
		}
	}

	d.SetId(tpgresource.GetResourceNameFromSelfLink(billingAccount.Name))
	if err := d.Set("name", billingAccount.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("display_name", billingAccount.DisplayName); err != nil {
		return fmt.Errorf("Error setting display_name: %s", err)
	}
	if err := d.Set("open", billingAccount.Open); err != nil {
		return fmt.Errorf("Error setting open: %s", err)
	}

	return nil
}

func CanonicalBillingAccountName(ba string) string {
	if strings.HasPrefix(ba, "billingAccounts/") {
		return ba
	}

	return "billingAccounts/" + ba
}

func flattenBillingProjects(billingProjects []*cloudbilling.ProjectBillingInfo) []string {
	projectIds := make([]string, len(billingProjects))
	for i, billingProject := range billingProjects {
		projectIds[i] = billingProject.ProjectId
	}

	return projectIds
}
