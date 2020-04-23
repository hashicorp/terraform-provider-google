package google

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1beta1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	Config     *Config
}

func NewHealthcareDatasetIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	dataset := d.Get("dataset_id").(string)
	datasetId, err := parseHealthcareDatasetId(dataset, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", dataset), err)
	}

	return &HealthcareDatasetIamUpdater{
		resourceId: datasetId.datasetId(),
		Config:     config,
	}, nil
}

func DatasetIdParseFunc(d *schema.ResourceData, config *Config) error {
	datasetId, err := parseHealthcareDatasetId(d.Id(), config)
	if err != nil {
		return err
	}

	d.Set("dataset_id", datasetId.datasetId())
	d.SetId(datasetId.datasetId())
	return nil
}

func (u *HealthcareDatasetIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientHealthcare.Projects.Locations.Datasets.GetIamPolicy(u.resourceId).Do()

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

	_, err = u.Config.clientHealthcare.Projects.Locations.Datasets.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
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
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a healthcare policy: {{err}}", err)
	}
	return out, nil
}

func healthcareToResourceManagerPolicy(p *healthcare.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a healthcare policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
