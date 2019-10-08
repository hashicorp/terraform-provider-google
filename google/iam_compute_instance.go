package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

var IamComputeInstanceSchema = map[string]*schema.Schema{
	"instance_name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},

	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},

	"zone": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type ComputeInstanceIamUpdater struct {
	project    string
	zone       string
	resourceId string
	Config     *Config
}

func NewComputeInstanceIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return nil, err
	}

	return &ComputeInstanceIamUpdater{
		project:    project,
		zone:       zone,
		resourceId: d.Get("instance_name").(string),
		Config:     config,
	}, nil
}

func ComputeInstanceIdParseFunc(d *schema.ResourceData, config *Config) error {
	parts := strings.Split(d.Id(), "/")
	var fv *ZonalFieldValue
	if len(parts) == 3 {
		// {project}/{zone}/{name} syntax
		fv = &ZonalFieldValue{
			Project:      parts[0],
			Zone:         parts[1],
			Name:         parts[2],
			resourceType: "instances",
		}
	} else if len(parts) == 2 {
		// /{zone}/{name} syntax
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		fv = &ZonalFieldValue{
			Project:      project,
			Zone:         parts[0],
			Name:         parts[1],
			resourceType: "instances",
		}
	} else {
		// We either have a name or a full self link, so use the field helper
		var err error
		fv, err = ParseInstanceFieldValue(d.Id(), d, config)
		if err != nil {
			return err
		}
	}

	d.Set("project", fv.Project)
	d.Set("zone", fv.Zone)
	d.Set("instance_name", fv.Name)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *ComputeInstanceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientCompute.Instances.GetIamPolicy(u.project, u.zone, u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := computeToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *ComputeInstanceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	computePolicy, err := resourceManagerToComputePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &compute.ZoneSetPolicyRequest{
		Policy: computePolicy,
	}
	_, err = u.Config.clientCompute.Instances.SetIamPolicy(u.project, u.zone, u.resourceId, req).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ComputeInstanceIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/zones/%s/instances/%s", u.project, u.zone, u.resourceId)
}

func (u *ComputeInstanceIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-compute-Instance-%s-%s-%s", u.project, u.zone, u.resourceId)
}

func (u *ComputeInstanceIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Compute Instance %s/%s/%s", u.project, u.zone, u.resourceId)
}
