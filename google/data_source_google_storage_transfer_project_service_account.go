package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	if err := d.Set("email", serviceAccount.AccountEmail); err != nil {
		return fmt.Errorf("Error reading email: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	return nil
}
