package google

import (
	"fmt"
	"strings"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
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
			"kubernetes_subnetwork_access": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	d.SetId(fmt.Sprintf("%s/%s", hostProject, serviceProject))
	if err = computeOperationWait(config.clientCompute, op, hostProject, "Enabling Shared VPC Resource"); err != nil {
		return err
	}
	// k8s subnetwork access is configured through the cloud console on the same page as
	// shared VPC project links, so to meet user expectations, we're adding it here.
	// Behind the scenes, a k8s subnetwork access request just adds two service accounts
	// with the roles/compute.networkUser to the subnetworks in question.  This is
	// acheivable using the existing iam code, mostly.
	if v, ok := d.GetOk("kubernetes_subnetwork_access"); ok {
		proj, err := config.clientResourceManager.Projects.Get(serviceProject).Do()
		if err != nil {
			return err
		}
		b := &cloudresourcemanager.Binding{
			Members: []string{fmt.Sprintf("serviceAccount:%d@cloudservices.gserviceaccount.com", proj.ProjectNumber),
				fmt.Sprintf("serviceAccount:service-%d@container-engine-robot.iam.gserviceaccount.com", proj.ProjectNumber)},
			Role: "roles/compute.networkUser",
		}
		for _, subnet := range v.([]interface{}) {
			iamUpdater, err := getUpdater(subnet.(string), d, config)
			if err != nil {
				return err
			}
			err = iamPolicyReadModifyWrite(iamUpdater, func(ep *cloudresourcemanager.Policy) error {
				ep.Bindings = mergeBindings(append(ep.Bindings, b))
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
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

	// We don't want to iterate through all subnetworks, checking to see which
	// ones have k8s access.  In that sense, this resource isn't 'authoritative'.
	// Instead, we just check to make sure that the ones that are listed are still
	// enabled and working, catching out of band changes to those.
	if v, ok := d.GetOk("kubernetes_subnetwork_access"); ok {
		k8s_subnets := make([]string, 0, len(v.([]interface{})))
		proj, err := config.clientResourceManager.Projects.Get(serviceProject).Do()
		if err != nil {
			return err
		}
		for _, subnet := range v.([]interface{}) {
			iamUpdater, err := getUpdater(subnet.(string), d, config)
			if err != nil {
				return err
			}
			p, err := iamUpdater.GetResourceIamPolicy()
			if err == nil {
				for _, b := range p.Bindings {
					if b.Role == "roles/compute.networkUser" {
						for _, m := range b.Members {
							if m == fmt.Sprintf("serviceAccount:service-%d@container-engine-robot.iam.gserviceaccount.com", proj.ProjectNumber) {
								k8s_subnets = append(k8s_subnets, subnet.(string))
								break
							}
						}
					}
				}
			}
		}
		d.Set("kubernetes_subnetwork_access", k8s_subnets)
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
	if v, ok := d.GetOk("kubernetes_subnetwork_access"); ok {
		proj, err := config.clientResourceManager.Projects.Get(serviceProject).Do()
		if err != nil {
			return err
		}
		for _, subnet := range v.([]interface{}) {
			iamUpdater, err := getUpdater(subnet.(string), d, config)
			if err != nil {
				// Try to avoid failing if possible - ignore the error if, for instance, the subnet
				// has already been deleted.
				continue
			}

			err = iamPolicyReadModifyWrite(iamUpdater, func(ep *cloudresourcemanager.Policy) error {
				bindings := make([]*cloudresourcemanager.Binding, 0)
				for _, binding := range ep.Bindings {
					if binding.Role == "roles/compute.networkUser" {
						b := &cloudresourcemanager.Binding{
							Role:    binding.Role,
							Members: make([]string, 0),
						}
						for _, m := range binding.Members {
							if m == fmt.Sprintf("serviceAccount:service-%d@container-engine-robot.iam.gserviceaccount.com", proj.ProjectNumber) {
								continue
							} else if m == fmt.Sprintf("serviceAccount:%d@cloudservices.gserviceaccount.com", proj.ProjectNumber) {
								continue
							} else {
								b.Members = append(b.Members, m)
							}
						}
						bindings = append(bindings, b)
					} else {
						bindings = append(bindings, binding)
					}
				}
				ep.Bindings = mergeBindings(bindings)
				return nil
			})
			if err != nil {
				return err
			}
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

func getUpdater(subnet string, d *schema.ResourceData, config *Config) (*ComputeSubnetworkIamUpdater, error) {
	region, err := getRegionFromResourceReference(subnet, d, config)
	if err != nil {
		return nil, err
	}
	resourceIdParts := strings.Split(subnet, "/")
	return &ComputeSubnetworkIamUpdater{
		project:    d.Get("host_project").(string),
		region:     region,
		resourceId: resourceIdParts[len(resourceIdParts)-1],
		Config:     config,
	}, nil
}
