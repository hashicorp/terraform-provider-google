package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleStorageTransferProjectServiceAccount() *schema.Resource {
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
		},
	}
}

func dataSourceGoogleStorageTransferProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	serviceAccount, err := config.clientStorageTransfer.GoogleServiceAccounts.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "Google Cloud Storage Transfer service account not found")
	}

	d.SetId(serviceAccount.AccountEmail)
	d.Set("email", serviceAccount.AccountEmail)
	d.Set("project", project)
	return nil
}
