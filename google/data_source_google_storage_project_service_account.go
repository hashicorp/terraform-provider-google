package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleStorageProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageProjectServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"user_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_address": {
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

func dataSourceGoogleStorageProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	serviceAccountGetRequest := config.NewStorageClient(userAgent).Projects.ServiceAccount.Get(project)

	if v, ok := d.GetOk("user_project"); ok {
		serviceAccountGetRequest = serviceAccountGetRequest.UserProject(v.(string))
	}

	serviceAccount, err := serviceAccountGetRequest.Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCS service account not found")
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("email_address", serviceAccount.EmailAddress); err != nil {
		return fmt.Errorf("Error setting email_address: %s", err)
	}
	if err := d.Set("member", "serviceAccount:"+serviceAccount.EmailAddress); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}

	d.SetId(serviceAccount.EmailAddress)

	return nil
}
