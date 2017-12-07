package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

var IamServiceAccountSchema = map[string]*schema.Schema{
	"service_account_id": &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateRegexp(ServiceAccountLinkRegex),
	},
}

type ServiceAccountIamUpdater struct {
	serviceAccountId string
	Config           *Config
}

func NewServiceAccountIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	return &ServiceAccountIamUpdater{
		serviceAccountId: d.Get("service_account_id").(string),
		Config:           config,
	}, nil
}

func (u *ServiceAccountIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientIAM.Projects.ServiceAccounts.GetIamPolicy(u.serviceAccountId).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for %s: %s", u.DescribeResource(), err)
	}

	cloudResourcePolicy, err := iamToResourceManagerPolicy(p)
	if err != nil {
		return nil, err
	}

	return cloudResourcePolicy, nil
}

func (u *ServiceAccountIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	iamPolicy, err := resourceManagerToIamPolicy(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.clientIAM.Projects.ServiceAccounts.SetIamPolicy(u.GetResourceId(), &iam.SetIamPolicyRequest{
		Policy: iamPolicy,
	}).Do()

	if err != nil {
		return fmt.Errorf("Error setting IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return nil
}

func (u *ServiceAccountIamUpdater) GetResourceId() string {
	return u.serviceAccountId
}

func (u *ServiceAccountIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-service-account-%s", u.serviceAccountId)
}

func (u *ServiceAccountIamUpdater) DescribeResource() string {
	return fmt.Sprintf("service account '%s'", u.serviceAccountId)
}

func resourceManagerToIamPolicy(p *cloudresourcemanager.Policy) (policy *iam.Policy, err error) {
	policy = &iam.Policy{}

	err = Convert(p, policy)

	return
}

func iamToResourceManagerPolicy(p *iam.Policy) (policy *cloudresourcemanager.Policy, err error) {
	policy = &cloudresourcemanager.Policy{}

	err = Convert(p, policy)

	return
}
