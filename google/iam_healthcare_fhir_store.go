package google

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	d          *schema.ResourceData
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
		d:          d,
		Config:     config,
	}, nil
}

func FhirStoreIdParseFunc(d *schema.ResourceData, config *Config) error {
	fhirStoreId, err := parseHealthcareFhirStoreId(d.Id(), config)
	if err != nil {
		return err
	}
	if err := d.Set("fhir_store_id", fhirStoreId.fhirStoreId()); err != nil {
		return fmt.Errorf("Error setting fhir_store_id: %s", err)
	}
	d.SetId(fhirStoreId.fhirStoreId())
	return nil
}

func (u *HealthcareFhirStoreIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.FhirStores.GetIamPolicy(u.resourceId).Do()

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

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewHealthcareClient(userAgent).Projects.Locations.Datasets.FhirStores.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
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
