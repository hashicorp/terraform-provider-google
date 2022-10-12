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
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	serviceAccount, err := config.NewStorageTransferClient(userAgent).GoogleServiceAccounts.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "Google Cloud Storage Transfer service account not found")
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
