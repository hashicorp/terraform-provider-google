package google

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeSharedVpcServiceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcServiceProjectCreate,
		Read:   resourceComputeSharedVpcServiceProjectRead,
		Delete: resourceComputeSharedVpcServiceProjectDelete,

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

func resourceComputeSharedVpcServiceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)

	req := &compute.ProjectsEnableXpnResourceRequest{
		XpnResource: &compute.XpnResourceId{
			Id:   serviceProject,
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

	d.SetId(fmt.Sprintf("%s/%s", hostProject, serviceProject))

	return nil
}

func resourceComputeSharedVpcServiceProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)

	req := config.clientCompute.Projects.GetXpnResources(hostProject)
	if err := req.Pages(context.Background(), func(page *compute.ProjectsGetXpnResources) error {
		for _, xpnResourceId := range page.Resources {
			if xpnResourceId.Type == "PROJECT" && xpnResourceId.Id == serviceProject {
				return nil
			}
		}
		return fmt.Errorf("%s is not a service project of %s", serviceProject, hostProject)
	}); err != nil {
		log.Printf("[WARN] %s", err)
		d.SetId("")
	}

	return nil
}

func resourceComputeSharedVpcServiceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)

	req := &compute.ProjectsDisableXpnResourceRequest{
		XpnResource: &compute.XpnResourceId{
			Id:   serviceProject,
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
