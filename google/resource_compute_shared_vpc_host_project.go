package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeSharedVpcHostProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcHostProjectCreate,
		Read:   resourceComputeSharedVpcHostProjectRead,
		Delete: resourceComputeSharedVpcHostProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the project that will serve as a Shared VPC host project`,
			},
		},
	}
}

func resourceComputeSharedVpcHostProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	hostProject := d.Get("project").(string)
	op, err := config.NewComputeBetaClient(userAgent).Projects.EnableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error enabling Shared VPC Host %q: %s", hostProject, err)
	}

	d.SetId(hostProject)

	err = computeOperationWaitTime(config, op, hostProject, "Enabling Shared VPC Host", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func resourceComputeSharedVpcHostProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	hostProject := d.Id()

	project, err := config.NewComputeBetaClient(userAgent).Projects.Get(hostProject).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", hostProject))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing Shared VPC host resource %q because it's not enabled server-side", hostProject)
		d.SetId("")
	}

	if err := d.Set("project", hostProject); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceComputeSharedVpcHostProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	hostProject := d.Get("project").(string)

	op, err := config.NewComputeBetaClient(userAgent).Projects.DisableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error disabling Shared VPC Host %q: %s", hostProject, err)
	}

	err = computeOperationWaitTime(config, op, hostProject, "Disabling Shared VPC Host", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
