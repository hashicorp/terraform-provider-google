// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/bigtableadmin/v2"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamBigtableInstanceSchema = map[string]*schema.Schema{
	"instance": {
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
}

type BigtableInstanceIamUpdater struct {
	project  string
	instance string
	d        tpgresource.TerraformResourceData
	Config   *transport_tpg.Config
}

func NewBigtableInstanceUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}

	return &BigtableInstanceIamUpdater{
		project:  project,
		instance: d.Get("instance").(string),
		d:        d,
		Config:   config,
	}, nil
}

func BigtableInstanceIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	fv, err := tpgresource.ParseProjectFieldValue("instances", d.Id(), "project", d, config, false)
	if err != nil {
		return err
	}

	if err := d.Set("project", fv.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("instance", fv.Name); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *BigtableInstanceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &bigtableadmin.GetIamPolicyRequest{}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewBigTableProjectsInstancesClient(userAgent).GetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := bigtableToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *BigtableInstanceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	bigtablePolicy, err := resourceManagerToBigtablePolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &bigtableadmin.SetIamPolicyRequest{Policy: bigtablePolicy}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewBigTableProjectsInstancesClient(userAgent).SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BigtableInstanceIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/instances/%s", u.project, u.instance)
}

func (u *BigtableInstanceIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-bigtable-instance-%s-%s", u.project, u.instance)
}

func (u *BigtableInstanceIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Bigtable Instance %s/%s", u.project, u.instance)
}

func resourceManagerToBigtablePolicy(p *cloudresourcemanager.Policy) (*bigtableadmin.Policy, error) {
	out := &bigtableadmin.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a bigtable policy to a cloudresourcemanager policy: {{err}}", err)
	}
	return out, nil
}

func bigtableToResourceManagerPolicy(p *bigtableadmin.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudresourcemanager policy to a bigtable policy: {{err}}", err)
	}
	return out, nil
}
