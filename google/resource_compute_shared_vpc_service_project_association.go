package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeSharedVpcServiceProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcServiceProjectAssociationCreate,
		Read:   resourceComputeSharedVpcServiceProjectAssociationRead,
		Delete: resourceComputeSharedVpcServiceProjectAssociationDelete,

		Schema: map[string]*schema.Schema{
			"host_project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeSharedVpcServiceProjectAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)
	if err := enableXpnResource(config, hostProject, serviceProject); err != nil {
		return fmt.Errorf("Error enabling Shared VPC service project %q: %s", serviceProject, err)
	}

	id := hostProjectServiceProjectHash(hostProject, serviceProject)
	log.Printf("[DEBUG] Shared VPC Service Project association %q created", id)
	d.SetId(id)

	return resourceComputeSharedVpcServiceProjectAssociationRead(d, meta)
}

func resourceComputeSharedVpcServiceProjectAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	project, err := config.clientCompute.Projects.Get(hostProject).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", hostProject))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing Shared VPC Service Project association because %q is no longer a Shared VPC Host", hostProject)
		d.SetId("")
		return nil
	}

	serviceProjects, err := findXpnResources(config, hostProject)
	if err != nil {
		return err
	}

	serviceProject := d.Get("service_project").(string)
	found := false
	for _, sp := range serviceProjects {
		if sp == serviceProject {
			found = true
			break
		}
	}
	if !found {
		// The association no longer exists.
		d.SetId("")
		return nil
	}

	id := hostProjectServiceProjectHash(hostProject, serviceProject)
	log.Printf("[DEBUG] Computed Shared VPC Service Project association %q", id)
	d.SetId(id)

	return nil
}

func resourceComputeSharedVpcServiceProjectAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)

	if err := disableXpnResource(config, hostProject, serviceProject); err != nil {
		if !isDisabledXpnResourceError(err) {
			return fmt.Errorf("Error disabling Shared VPC Resource %q: %s", serviceProject, err)
		}
	}

	log.Printf("[DEBUG] Shared VPC Service Project association %q deleted", d.Id())
	d.SetId("")
	return nil
}

func hostProjectServiceProjectHash(hostProject, serviceProject string) string {
	return fmt.Sprintf("a-%s%d", hostProject, hashcode.String(serviceProject))
}
