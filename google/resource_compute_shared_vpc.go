package google

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeSharedVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcCreate,
		Read:   resourceComputeSharedVpcRead,
		Update: resourceComputeSharedVpcUpdate,
		Delete: resourceComputeSharedVpcDelete,

		Schema: map[string]*schema.Schema{
			"host_project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_projects": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceComputeSharedVpcCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	op, err := config.clientCompute.Projects.EnableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error enabling Shared VPC Host %q: %s", hostProject, err)
	}

	d.SetId(hostProject)

	err = computeOperationWait(config, op, hostProject, "Enabling Shared VPC Host")
	if err != nil {
		d.SetId("")
		return err
	}

	if v, ok := d.GetOk("service_projects"); ok {
		serviceProjects := convertStringArr(v.(*schema.Set).List())
		for _, project := range serviceProjects {
			if err = enableResource(config, hostProject, project); err != nil {
				return fmt.Errorf("Error enabling Shared VPC service project %q: %s", project, err)
			}
		}
	}

	return resourceComputeSharedVpcRead(d, meta)
}

func resourceComputeSharedVpcRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)

	project, err := config.clientCompute.Projects.Get(hostProject).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", hostProject))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing Shared VPC host resource %q because it's not enabled server-side", hostProject)
		d.SetId("")
	}

	serviceProjects := []string{}
	req := config.clientCompute.Projects.GetXpnResources(hostProject)
	if err := req.Pages(context.Background(), func(page *compute.ProjectsGetXpnResources) error {
		for _, xpnResourceId := range page.Resources {
			if xpnResourceId.Type == "PROJECT" {
				serviceProjects = append(serviceProjects, xpnResourceId.Id)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("Error reading Shared VPC service projects for host %q: %s", hostProject, err)
	}

	d.Set("service_projects", serviceProjects)

	return nil
}

func resourceComputeSharedVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	hostProject := d.Get("host_project").(string)

	if d.HasChange("service_projects") {
		old, new := d.GetChange("service_projects")
		oldMap := convertArrToMap(old.(*schema.Set).List())
		newMap := convertArrToMap(new.(*schema.Set).List())

		for project, _ := range oldMap {
			if _, ok := newMap[project]; !ok {
				// The project is in the old config but not the new one, disable it
				if err := disableResource(config, hostProject, project); err != nil {
					return fmt.Errorf("Error disabling Shared VPC service project %q: %s", project, err)
				}
			}
		}

		for project, _ := range newMap {
			if _, ok := oldMap[project]; !ok {
				// The project is in the new config but not the old one, enable it
				if err := enableResource(config, hostProject, project); err != nil {
					return fmt.Errorf("Error enabling Shared VPC service project %q: %s", project, err)
				}
			}
		}
	}

	return resourceComputeSharedVpcRead(d, meta)
}

func resourceComputeSharedVpcDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	hostProject := d.Get("host_project").(string)

	serviceProjects := convertStringArr(d.Get("service_projects").(*schema.Set).List())
	for _, project := range serviceProjects {
		if err := disableResource(config, hostProject, project); err != nil {
			return fmt.Errorf("Error disabling Shared VPC Resource %q: %s", project, err)
		}
	}

	op, err := config.clientCompute.Projects.DisableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error disabling Shared VPC Host %q: %s", hostProject, err)
	}

	err = computeOperationWait(config, op, hostProject, "Disabling Shared VPC Host")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func enableResource(config *Config, hostProject, project string) error {
	req := &compute.ProjectsEnableXpnResourceRequest{
		XpnResource: &compute.XpnResourceId{
			Id:   project,
			Type: "PROJECT",
		},
	}
	op, err := config.clientCompute.Projects.EnableXpnResource(hostProject, req).Do()
	if err != nil {
		return err
	}
	if err = computeOperationWait(config, op, hostProject, "Enabling Shared VPC Resource"); err != nil {
		return err
	}
	return nil
}

func disableResource(config *Config, hostProject, project string) error {
	req := &compute.ProjectsDisableXpnResourceRequest{
		XpnResource: &compute.XpnResourceId{
			Id:   project,
			Type: "PROJECT",
		},
	}
	op, err := config.clientCompute.Projects.DisableXpnResource(hostProject, req).Do()
	if err != nil {
		return err
	}
	if err = computeOperationWait(config, op, hostProject, "Disabling Shared VPC Resource"); err != nil {
		return err
	}
	return nil
}
