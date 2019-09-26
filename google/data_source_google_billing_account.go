package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"google.golang.org/api/cloudbilling/v1"
)

func dataSourceGoogleBillingAccount() *schema.Resource {
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
		},
	}
}

func dataSourceBillingAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	open, openOk := d.GetOkExists("open")

	var billingAccount *cloudbilling.BillingAccount
	if v, ok := d.GetOk("billing_account"); ok {
		resp, err := config.clientBilling.BillingAccounts.Get(canonicalBillingAccountName(v.(string))).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Billing Account Not Found : %s", v))
		}

		if openOk && resp.Open != open.(bool) {
			return fmt.Errorf("Billing account not found: %s", v)
		}

		billingAccount = resp
	} else if v, ok := d.GetOk("display_name"); ok {
		token := ""
		for paginate := true; paginate; {
			resp, err := config.clientBilling.BillingAccounts.List().PageToken(token).Do()
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

	resp, err := config.clientBilling.BillingAccounts.Projects.List(billingAccount.Name).Do()
	if err != nil {
		return fmt.Errorf("Error reading billing account projects: %s", err)
	}
	projectIds := flattenBillingProjects(resp.ProjectBillingInfo)

	d.SetId(GetResourceNameFromSelfLink(billingAccount.Name))
	d.Set("name", billingAccount.Name)
	d.Set("display_name", billingAccount.DisplayName)
	d.Set("open", billingAccount.Open)
	d.Set("project_ids", projectIds)

	return nil
}

func canonicalBillingAccountName(ba string) string {
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
