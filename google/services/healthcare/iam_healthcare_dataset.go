// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	healthcare "google.golang.org/api/healthcare/v1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamHealthcareDatasetSchema = map[string]*schema.Schema{
	"dataset_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type HealthcareDatasetIamUpdater struct {
	resourceId string
	d          tpgresource.TerraformResourceData
	Config     *transport_tpg.Config
}

func NewHealthcareDatasetIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	dataset := d.Get("dataset_id").(string)
	datasetId, err := ParseHealthcareDatasetId(dataset, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for %s: {{err}}", dataset), err)
	}

	return &HealthcareDatasetIamUpdater{
		resourceId: datasetId.DatasetId(),
		d:          d,
		Config:     config,
	}, nil
}

func DatasetIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	datasetId, err := ParseHealthcareDatasetId(d.Id(), config)
	if err != nil {
		return err
	}

	if err := d.Set("dataset_id", datasetId.DatasetId()); err != nil {
		return fmt.Errorf("Error setting dataset_id: %s", err)
	}
	d.SetId(datasetId.DatasetId())
	return nil
}

func (u *HealthcareDatasetIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.GetIamPolicy(u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := healthcareToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *HealthcareDatasetIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	healthcarePolicy, err := resourceManagerToHealthcarePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
		Policy: healthcarePolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *HealthcareDatasetIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *HealthcareDatasetIamUpdater) GetMutexKey() string {
	return u.resourceId
}

func (u *HealthcareDatasetIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Healthcare Dataset %q", u.resourceId)
}

func resourceManagerToHealthcarePolicy(p *cloudresourcemanager.Policy) (*healthcare.Policy, error) {
	out := &healthcare.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a healthcare policy: {{err}}", err)
	}
	return out, nil
}

func healthcareToResourceManagerPolicy(p *healthcare.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a healthcare policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
