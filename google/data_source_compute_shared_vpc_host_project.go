package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceComputeSharedVpcHostProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeSharedVpcHostProjectRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_project": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeSharedVpcHostProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error getting service project for Shared VPC Host project: %s", err)
	}

	op, err := config.clientComputeBeta.Projects.GetXpnHost(project).Do()
	if err != nil {
		return fmt.Errorf("Error reading Shared VPC Host for project %s: %s", project, err)
	}

	d.SetId(fmt.Sprintf("projects/%s", project))
	d.Set("project", project)
	d.Set("host_project", op.Name)

	return nil
}
