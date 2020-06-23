package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamOrganizationSchema = map[string]*schema.Schema{
	"org_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The numeric ID of the organization in which you want to manage the audit logging config.`,
	},
}

type OrganizationIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewOrganizationIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	return &OrganizationIamUpdater{
		resourceId: d.Get("org_id").(string),
		Config:     config,
	}, nil
}

func OrgIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("org_id", d.Id())
	return nil
}

func (u *OrganizationIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientResourceManager.Organizations.GetIamPolicy("organizations/"+u.resourceId, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return p, nil
}

func (u *OrganizationIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	_, err := u.Config.clientResourceManager.Organizations.SetIamPolicy("organizations/"+u.resourceId, &cloudresourcemanager.SetIamPolicyRequest{
		Policy:     policy,
		UpdateMask: "bindings,etag,auditConfigs",
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *OrganizationIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *OrganizationIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-organization-%s", u.resourceId)
}

func (u *OrganizationIamUpdater) DescribeResource() string {
	return fmt.Sprintf("organization %q", u.resourceId)
}
