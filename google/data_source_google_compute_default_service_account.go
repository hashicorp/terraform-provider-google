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
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

	if v, ok := d.GetOk("project_id"); ok {
		project = v.(string)
	}

	projectCompResource, err := config.clientCompute.Projects.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "GCE service account not found")
	}

	d.SetId(projectCompResource.DefaultServiceAccount)
	d.Set("email", projectCompResource.DefaultServiceAccount)
	d.Set("project_id", project)
	return nil
}
