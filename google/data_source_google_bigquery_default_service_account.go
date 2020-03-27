package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	d.Set("email", projectResource.Email)
	d.Set("project", project)
	return nil
}
