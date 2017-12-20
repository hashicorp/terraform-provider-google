package google

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleBillingAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBillingAccountRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceBillingAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name, nameOk := d.GetOk("name")
	displayName, displayNameOk := d.GetOk("display_name")
	if nameOk == displayNameOk {
		return fmt.Errorf("One of ['name', 'display_name'] must be set to read billing accounts")
	}

	var billingAccount *cloudbilling.BillingAccount
	if nameOk {
		resp, err := config.clientBilling.BillingAccounts.Get(name.(string)).Do()
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
				return fmt.Errorf("Billing account not found: %s", name)
			}

			return fmt.Errorf("Error reading billing account: %s", err)
		}

		billingAccount = resp
	} else {
		resp, err := config.clientBilling.BillingAccounts.List().Do()
		if err != nil {
			return fmt.Errorf("Error reading billing account: %s", err)
		}

		for _, ba := range resp.BillingAccounts {
			if ba.DisplayName == displayName.(string) {
				if billingAccount != nil {
					return fmt.Errorf("More than one matching billing account found")
				}
				billingAccount = ba
			}
		}

		if billingAccount == nil {
			return fmt.Errorf("Billing account not found: %s", displayName)
		}
	}

	resp, err := config.clientBilling.BillingAccounts.Projects.List(billingAccount.Name).Do()
	if err != nil {
		return fmt.Errorf("Error reading billing account projects: %s", err)
	}
	projectIds := flattenBillingProjects(resp.ProjectBillingInfo)

	parts := strings.Split(billingAccount.Name, "/")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid billing account name. Expecting billingAccounts/{billing_account_id}")
	}

	d.SetId(parts[1])
	d.Set("name", billingAccount.Name)
	d.Set("display_name", billingAccount.DisplayName)
	d.Set("open", billingAccount.Open)
	d.Set("project_ids", projectIds)

	return nil
}

func flattenBillingProjects(billingProjects []*cloudbilling.ProjectBillingInfo) []string {
	projectIds := make([]string, len(billingProjects))
	for i, billingProject := range billingProjects {
		projectIds[i] = billingProject.ProjectId
	}

	return projectIds
}
