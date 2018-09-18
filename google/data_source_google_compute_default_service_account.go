package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeDefaultServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeDefaultServiceAccountRead,
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

func dataSourceGoogleComputeDefaultServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	projectCompResource, err := config.clientCompute.Projects.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCE service account not found")
	}

	d.SetId(projectCompResource.DefaultServiceAccount)
	d.Set("email", projectCompResource.DefaultServiceAccount)
	d.Set("project", project)
	return nil
}
