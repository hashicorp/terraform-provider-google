package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/compute/v1"

	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func resourceComputeSharedVpcServiceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcServiceProjectCreate,
		Read:   resourceComputeSharedVpcServiceProjectRead,
		Delete: resourceComputeSharedVpcServiceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
	if err = computeOperationWait(config.clientCompute, op, hostProject, "Enabling Shared VPC Resource"); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", hostProject, serviceProject))

	return nil
}

func resourceComputeSharedVpcServiceProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	split := strings.Split(d.Id(), "/")
	if len(split) != 2 {
		return fmt.Errorf("Error parsing resource ID %s", d.Id())
	}
	hostProject := split[0]
	serviceProject := split[1]

	associatedHostProject, err := config.clientCompute.Projects.GetXpnHost(serviceProject).Do()
	if err != nil {
		log.Printf("[WARN] Removing shared VPC service. The service project is not associated with any host")

		d.SetId("")
		return nil
	}

	if hostProject != associatedHostProject.Name {
		log.Printf("[WARN] Removing shared VPC service. Expected associated host project to be '%s', got '%s'", hostProject, associatedHostProject.Name)
		d.SetId("")
		return nil
	}

	d.Set("host_project", hostProject)
	d.Set("service_project", serviceProject)

	return nil
}

func resourceComputeSharedVpcServiceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	hostProject := d.Get("host_project").(string)
	serviceProject := d.Get("service_project").(string)

	if err := disableXpnResource(config, hostProject, serviceProject); err != nil {
		// Don't fail if the service project is already disabled.
		if !isDisabledXpnResourceError(err) {
			return fmt.Errorf("Error disabling Shared VPC Resource %q: %s", serviceProject, err)
		}
	}

	return nil
}

func disableXpnResource(config *Config, hostProject, project string) error {
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
	if err = computeOperationWait(config.clientCompute, op, hostProject, "Disabling Shared VPC Resource"); err != nil {
		return err
	}
	return nil
}

func isDisabledXpnResourceError(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && len(gerr.Errors) > 0 && gerr.Errors[0].Reason == "invalidResourceUsage" {
			return true
		}
	}
	return false
}
