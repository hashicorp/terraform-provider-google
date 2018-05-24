package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountRead,

		Schema: map[string]*schema.Schema{
			// Required
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Get the account id as the fully qualified name
	accountID := d.Get("account_id").(string)

	// If the account id isn't already the fully qualified name
	if !strings.HasPrefix(accountID, "projects/") {

		// If the account id is an email
		if strings.Contains(accountID, "@") {
			accountID = "projects/-/serviceAccounts/" + accountID
		} else {
			// Get the project from the resource or fallback to the project
			// in the provider configuration
			project, err := getProject(d, config)
			if err != nil {
				return err
			}
			// If the account id doesn't contain the email, build it
			accountID = fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", accountID, project)
		}
	}

	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(accountID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", accountID))
	}

	d.SetId(sa.Name)
	d.Set("email", sa.Email)
	d.Set("unique_id", sa.UniqueId)
	d.Set("project", sa.ProjectId)
	d.Set("account_id", strings.Split(sa.Email, "@")[0])
	d.Set("name", sa.Name)
	d.Set("display_name", sa.DisplayName)

	return nil
}
