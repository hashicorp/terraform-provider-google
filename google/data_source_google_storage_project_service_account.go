package google

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleStorageProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageProjectServiceAccountRead,
		Schema: map[string]*schema.Schema{
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

	serviceAccount, err := config.clientStorage.Projects.ServiceAccount.Get(project).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return fmt.Errorf("GCS service account not found")
		}

		return fmt.Errorf("Error reading GCS service account: %s", err)
	}

	d.SetId(serviceAccount.EmailAddress)
	d.Set("email_address", serviceAccount.EmailAddress)

	return nil
}
