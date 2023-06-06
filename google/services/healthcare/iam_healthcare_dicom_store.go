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

var IamHealthcareDicomStoreSchema = map[string]*schema.Schema{
	"dicom_store_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type HealthcareDicomStoreIamUpdater struct {
	resourceId string
	d          tpgresource.TerraformResourceData
	Config     *transport_tpg.Config
}

func NewHealthcareDicomStoreIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	dicomStore := d.Get("dicom_store_id").(string)
	dicomStoreId, err := ParseHealthcareDicomStoreId(dicomStore, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for %s: {{err}}", dicomStore), err)
	}

	return &HealthcareDicomStoreIamUpdater{
		resourceId: dicomStoreId.DicomStoreId(),
		d:          d,
		Config:     config,
	}, nil
}

func DicomStoreIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	dicomStoreId, err := ParseHealthcareDicomStoreId(d.Id(), config)
	if err != nil {
		return err
	}
	if err := d.Set("dicom_store_id", dicomStoreId.DicomStoreId()); err != nil {
		return fmt.Errorf("Error setting dicom_store_id: %s", err)
	}
	d.SetId(dicomStoreId.DicomStoreId())
	return nil
}

func (u *HealthcareDicomStoreIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.DicomStores.GetIamPolicy(u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := healthcareToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *HealthcareDicomStoreIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	healthcarePolicy, err := resourceManagerToHealthcarePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.DicomStores.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
		Policy: healthcarePolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *HealthcareDicomStoreIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *HealthcareDicomStoreIamUpdater) GetMutexKey() string {
	return u.resourceId
}

func (u *HealthcareDicomStoreIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Healthcare DicomStore %q", u.resourceId)
}
