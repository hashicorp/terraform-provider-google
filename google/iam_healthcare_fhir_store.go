package google

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1beta1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamHealthcareFhirStoreSchema = map[string]*schema.Schema{
	"fhir_store_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type HealthcareFhirStoreIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewHealthcareFhirStoreIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	fhirStore := d.Get("fhir_store_id").(string)
	fhirStoreId, err := parseHealthcareFhirStoreId(fhirStore, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", fhirStore), err)
	}

	return &HealthcareFhirStoreIamUpdater{
		resourceId: fhirStoreId.fhirStoreId(),
		Config:     config,
	}, nil
}

func FhirStoreIdParseFunc(d *schema.ResourceData, config *Config) error {
	fhirStoreId, err := parseHealthcareFhirStoreId(d.Id(), config)
	if err != nil {
		return err
	}
	d.Set("fhir_store_id", fhirStoreId.fhirStoreId())
	d.SetId(fhirStoreId.fhirStoreId())
	return nil
}

func (u *HealthcareFhirStoreIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientHealthcare.Projects.Locations.Datasets.FhirStores.GetIamPolicy(u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := healthcareToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *HealthcareFhirStoreIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	healthcarePolicy, err := resourceManagerToHealthcarePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	_, err = u.Config.clientHealthcare.Projects.Locations.Datasets.FhirStores.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
		Policy: healthcarePolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *HealthcareFhirStoreIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *HealthcareFhirStoreIamUpdater) GetMutexKey() string {
	return u.resourceId
}

func (u *HealthcareFhirStoreIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Healthcare FhirStore %q", u.resourceId)
}
