package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamProjectSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
}

type ProjectIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewProjectIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	pid, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return &ProjectIamUpdater{
		resourceId: pid,
		Config:     config,
	}, nil
}

func (u *ProjectIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientResourceManager.Projects.GetIamPolicy(u.resourceId,
		&cloudresourcemanager.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return p, nil
}

func (u *ProjectIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	_, err := u.Config.clientResourceManager.Projects.SetIamPolicy(u.resourceId, &cloudresourcemanager.SetIamPolicyRequest{
		Policy:     policy,
		UpdateMask: "bindings",
	}).Do()

	if err != nil {
		return fmt.Errorf("Error setting IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return nil
}

func (u *ProjectIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *ProjectIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-project-%s", u.resourceId)
}

func (u *ProjectIamUpdater) DescribeResource() string {
	return fmt.Sprintf("project %q", u.resourceId)
}
