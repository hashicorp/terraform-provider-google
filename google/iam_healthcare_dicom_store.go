package google

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1beta1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	Config     *Config
}

func NewHealthcareDicomStoreIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	dicomStore := d.Get("dicom_store_id").(string)
	dicomStoreId, err := parseHealthcareDicomStoreId(dicomStore, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", dicomStore), err)
	}

	return &HealthcareDicomStoreIamUpdater{
		resourceId: dicomStoreId.dicomStoreId(),
		Config:     config,
	}, nil
}

func DicomStoreIdParseFunc(d *schema.ResourceData, config *Config) error {
	dicomStoreId, err := parseHealthcareDicomStoreId(d.Id(), config)
	if err != nil {
		return err
	}
	d.Set("dicom_store_id", dicomStoreId.dicomStoreId())
	d.SetId(dicomStoreId.dicomStoreId())
	return nil
}

func (u *HealthcareDicomStoreIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientHealthcare.Projects.Locations.Datasets.DicomStores.GetIamPolicy(u.resourceId).Do()

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

	_, err = u.Config.clientHealthcare.Projects.Locations.Datasets.DicomStores.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
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
