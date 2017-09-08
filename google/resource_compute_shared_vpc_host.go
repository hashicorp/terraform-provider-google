package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeSharedVpcHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcHostCreate,
		Read:   resourceComputeSharedVpcHostRead,
		Delete: resourceComputeSharedVpcHostDelete,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeSharedVpcHostCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.Projects.EnableXpnHost(projectID).Do()
	if err != nil {
		return fmt.Errorf("Error disabling XPN Host: %s", err)
	}

	d.SetId(projectID)

	err = computeOperationWait(config, op, projectID, "Enabling XPN Host")
	if err != nil {
		d.SetId("")
		return err
	}

	return resourceComputeSharedVpcHostRead(d, meta)
}

func resourceComputeSharedVpcHostRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	project, err := config.clientCompute.Projects.Get(projectID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", projectID))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing %s VPC host resource because it's not enabled server-side", projectID)
		d.SetId("")
	}

	return nil
}

func resourceComputeSharedVpcHostDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.Projects.DisableXpnHost(projectID).Do()
	if err != nil {
		return fmt.Errorf("Error disabling XPN Host: %s", err)
	}

	err = computeOperationWait(config, op, projectID, "Disabling XPN Host")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
