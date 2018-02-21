package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleStorageProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageProjectServiceAccountRead,
	}
}

func dataSourceGoogleStorageProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	serviceAccount, err := config.clientStorage.Projects.ServiceAccount.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCS service account not found")
	}

	d.SetId(serviceAccount.EmailAddress)

	return nil
}
