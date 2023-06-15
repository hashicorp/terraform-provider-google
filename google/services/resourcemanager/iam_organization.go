// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
	d          tpgresource.TerraformResourceData
	Config     *transport_tpg.Config
}

func NewOrganizationIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return &OrganizationIamUpdater{
		resourceId: d.Get("org_id").(string),
		d:          d,
		Config:     config,
	}, nil
}

func OrgIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	if err := d.Set("org_id", d.Id()); err != nil {
		return fmt.Errorf("Error setting org_id: %s", err)
	}
	return nil
}

func (u *OrganizationIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewResourceManagerClient(userAgent).Organizations.GetIamPolicy(
		"organizations/"+u.resourceId,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		},
	).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return p, nil
}

func (u *OrganizationIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewResourceManagerClient(userAgent).Organizations.SetIamPolicy("organizations/"+u.resourceId, &cloudresourcemanager.SetIamPolicyRequest{
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
