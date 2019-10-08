package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		},
	}
}

func dataSourceGoogleStorageProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	serviceAccountGetRequest := config.clientStorage.Projects.ServiceAccount.Get(project)

	if v, ok := d.GetOk("user_project"); ok {
		serviceAccountGetRequest = serviceAccountGetRequest.UserProject(v.(string))
	}

	serviceAccount, err := serviceAccountGetRequest.Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCS service account not found")
	}

	d.Set("project", project)
	d.Set("email_address", serviceAccount.EmailAddress)

	d.SetId(serviceAccount.EmailAddress)

	return nil
}
