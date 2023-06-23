// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

var IamServiceAccountSchema = map[string]*schema.Schema{
	"service_account_id": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: verify.ValidateRegexp(verify.ServiceAccountLinkRegex),
	},
}

type ServiceAccountIamUpdater struct {
	serviceAccountId string
	d                tpgresource.TerraformResourceData
	Config           *transport_tpg.Config
}

func NewServiceAccountIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return &ServiceAccountIamUpdater{
		serviceAccountId: d.Get("service_account_id").(string),
		d:                d,
		Config:           config,
	}, nil
}

func ServiceAccountIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	if err := d.Set("service_account_id", d.Id()); err != nil {
		return fmt.Errorf("Error setting service_account_id: %s", err)
	}
	return nil
}

func (u *ServiceAccountIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewIamClient(userAgent).Projects.ServiceAccounts.GetIamPolicy(u.serviceAccountId).OptionsRequestedPolicyVersion(tpgiamresource.IamPolicyVersion).Do()

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

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewIamClient(userAgent).Projects.ServiceAccounts.SetIamPolicy(u.GetResourceId(), &iam.SetIamPolicyRequest{
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
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a iam policy: {{err}}", err)
	}
	return out, nil
}

func iamToResourceManagerPolicy(p *iam.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a iam policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
