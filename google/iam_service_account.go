package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

var IamServiceAccountSchema = map[string]*schema.Schema{
	"service_account_id": {
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

func ServiceAccountIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("service_account_id", d.Id())
	return nil
}

func (u *ServiceAccountIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientIAM.Projects.ServiceAccounts.GetIamPolicy(u.serviceAccountId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
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
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
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

func resourceManagerToIamPolicy(p *cloudresourcemanager.Policy) (*iam.Policy, error) {
	out := &iam.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a iam policy: {{err}}", err)
	}
	return out, nil
}

func iamToResourceManagerPolicy(p *iam.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a iam policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
