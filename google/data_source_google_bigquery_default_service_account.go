package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleBigqueryDefaultServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleBigqueryDefaultServiceAccountRead,
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

func dataSourceGoogleBigqueryDefaultServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	projectResource, err := config.clientBigQuery.Projects.GetServiceAccount(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "BigQuery service account not found")
	}

	d.SetId(projectResource.Email)
	if err := d.Set("email", projectResource.Email); err != nil {
		return fmt.Errorf("Error reading email: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	return nil
}
