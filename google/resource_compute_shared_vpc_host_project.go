package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeSharedVpcHostProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcHostProjectCreate,
		Read:   resourceComputeSharedVpcHostProjectRead,
		Delete: resourceComputeSharedVpcHostProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeSharedVpcHostProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("project").(string)
	op, err := config.clientCompute.Projects.EnableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error enabling Shared VPC Host %q: %s", hostProject, err)
	}

	d.SetId(hostProject)

	err = computeOperationWait(config.clientCompute, op, hostProject, "Enabling Shared VPC Host")
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func resourceComputeSharedVpcHostProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Id()

	project, err := config.clientCompute.Projects.Get(hostProject).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", hostProject))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing Shared VPC host resource %q because it's not enabled server-side", hostProject)
		d.SetId("")
	}

	d.Set("project", hostProject)

	return nil
}

func resourceComputeSharedVpcHostProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	hostProject := d.Get("project").(string)

	op, err := config.clientCompute.Projects.DisableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error disabling Shared VPC Host %q: %s", hostProject, err)
	}

	err = computeOperationWait(config.clientCompute, op, hostProject, "Disabling Shared VPC Host")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
