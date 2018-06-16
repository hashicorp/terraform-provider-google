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
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
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

	serviceAccountName, err := serviceAccountFQN(d.Get("account_id").(string), d, config)
	if err != nil {
		return err
	}

	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(serviceAccountName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", serviceAccountName))
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
