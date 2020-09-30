package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeDefaultServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeDefaultServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeDefaultServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	projectCompResource, err := config.NewComputeClient(userAgent).Projects.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCE default service account")
	}

	serviceAccountName, err := serviceAccountFQN(projectCompResource.DefaultServiceAccount, d, config)
	if err != nil {
		return err
	}

	sa, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Get(serviceAccountName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", serviceAccountName))
	}

	d.SetId(sa.Name)
	if err := d.Set("email", sa.Email); err != nil {
		return fmt.Errorf("Error setting email: %s", err)
	}
	if err := d.Set("unique_id", sa.UniqueId); err != nil {
		return fmt.Errorf("Error setting unique_id: %s", err)
	}
	if err := d.Set("project", sa.ProjectId); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("name", sa.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("display_name", sa.DisplayName); err != nil {
		return fmt.Errorf("Error setting display_name: %s", err)
	}

	return nil
}
