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

var IamProjectSchema = map[string]*schema.Schema{
	"project": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: CompareProjectName,
	},
}

type ProjectIamUpdater struct {
	resourceId string
	d          tpgresource.TerraformResourceData
	Config     *transport_tpg.Config
}

func NewProjectIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return &ProjectIamUpdater{
		resourceId: d.Get("project").(string),
		d:          d,
		Config:     config,
	}, nil
}

func ProjectIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	if err := d.Set("project", d.Id()); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	return nil
}

func (u *ProjectIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	projectId := tpgresource.GetResourceNameFromSelfLink(u.resourceId)

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewResourceManagerClient(userAgent).Projects.GetIamPolicy(projectId,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return p, nil
}

func (u *ProjectIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	projectId := tpgresource.GetResourceNameFromSelfLink(u.resourceId)

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewResourceManagerClient(userAgent).Projects.SetIamPolicy(projectId,
		&cloudresourcemanager.SetIamPolicyRequest{
			Policy:     policy,
			UpdateMask: "bindings,etag,auditConfigs",
		}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *ProjectIamUpdater) GetMutexKey() string {
	return getProjectIamPolicyMutexKey(u.resourceId)
}

func (u *ProjectIamUpdater) DescribeResource() string {
	return fmt.Sprintf("project %q", u.resourceId)
}

func CompareProjectName(_, old, new string, _ *schema.ResourceData) bool {
	// We can either get "projects/project-id" or "project-id", so strip any prefixes
	return tpgresource.GetResourceNameFromSelfLink(old) == tpgresource.GetResourceNameFromSelfLink(new)
}

func getProjectIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-project-%s", pid)
}
