package google

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	projectResManager, err := config.clientResourceManager.Projects.Get(project).Do()
	if err != nil {
		return fmt.Errorf("Error reading project resource: %s", err)
	}

	d.SetId(projectResManager.ProjectId)
	d.Set("project_number", strconv.FormatInt(projectResManager.ProjectNumber, 10))
	d.Set("name", projectResManager.Name)

	return nil
}
