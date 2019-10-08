package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
)

var IamComputeSubnetworkSchema = map[string]*schema.Schema{
	"subnetwork": {
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

	"region": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type ComputeSubnetworkIamUpdater struct {
	project    string
	region     string
	resourceId string
	Config     *Config
}

func NewComputeSubnetworkIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	return &ComputeSubnetworkIamUpdater{
		project:    project,
		region:     region,
		resourceId: d.Get("subnetwork").(string),
		Config:     config,
	}, nil
}

func ComputeSubnetworkIdParseFunc(d *schema.ResourceData, config *Config) error {
	parts := strings.Split(d.Id(), "/")
	var fv *RegionalFieldValue
	if len(parts) == 3 {
		// {project}/{region}/{name} syntax
		fv = &RegionalFieldValue{
			Project:      parts[0],
			Region:       parts[1],
			Name:         parts[2],
			resourceType: "subnetworks",
		}
	} else if len(parts) == 2 {
		// /{region}/{name} syntax
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		fv = &RegionalFieldValue{
			Project:      project,
			Region:       parts[0],
			Name:         parts[1],
			resourceType: "subnetworks",
		}
	} else {
		// We either have a name or a full self link, so use the field helper
		var err error
		fv, err = ParseSubnetworkFieldValue(d.Id(), d, config)
		if err != nil {
			return err
		}
	}
	d.Set("subnetwork", fv.Name)
	d.Set("project", fv.Project)
	d.Set("region", fv.Region)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *ComputeSubnetworkIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientCompute.Subnetworks.GetIamPolicy(u.project, u.region, u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := computeToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *ComputeSubnetworkIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	computePolicy, err := resourceManagerToComputePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &compute.RegionSetPolicyRequest{
		Policy: computePolicy,
	}
	_, err = u.Config.clientCompute.Subnetworks.SetIamPolicy(u.project, u.region, u.resourceId, req).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ComputeSubnetworkIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", u.project, u.region, u.resourceId)
}

func (u *ComputeSubnetworkIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-compute-subnetwork-%s-%s-%s", u.project, u.region, u.resourceId)
}

func (u *ComputeSubnetworkIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Compute Subnetwork %s/%s/%s", u.project, u.region, u.resourceId)
}

func resourceManagerToComputePolicy(p *cloudresourcemanager.Policy) (*compute.Policy, error) {
	out := &compute.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a resourcemanager policy to a compute policy: {{err}}", err)
	}
	return out, nil
}

func computeToResourceManagerPolicy(p *compute.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a compute policy to a resourcemanager policy: {{err}}", err)
	}
	return out, nil
}
