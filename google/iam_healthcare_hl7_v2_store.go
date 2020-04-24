package google

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamHealthcareHl7V2StoreSchema = map[string]*schema.Schema{
	"hl7_v2_store_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type HealthcareHl7V2StoreIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewHealthcareHl7V2StoreIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	hl7V2Store := d.Get("hl7_v2_store_id").(string)
	hl7V2StoreId, err := parseHealthcareHl7V2StoreId(hl7V2Store, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", hl7V2Store), err)
	}

	return &HealthcareHl7V2StoreIamUpdater{
		resourceId: hl7V2StoreId.hl7V2StoreId(),
		Config:     config,
	}, nil
}

func Hl7V2StoreIdParseFunc(d *schema.ResourceData, config *Config) error {
	hl7V2StoreId, err := parseHealthcareHl7V2StoreId(d.Id(), config)
	if err != nil {
		return err
	}
	d.Set("hl7_v2_store_id", hl7V2StoreId.hl7V2StoreId())
	d.SetId(hl7V2StoreId.hl7V2StoreId())
	return nil
}

func (u *HealthcareHl7V2StoreIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientHealthcare.Projects.Locations.Datasets.Hl7V2Stores.GetIamPolicy(u.resourceId).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := healthcareToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *HealthcareHl7V2StoreIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	healthcarePolicy, err := resourceManagerToHealthcarePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	_, err = u.Config.clientHealthcare.Projects.Locations.Datasets.Hl7V2Stores.SetIamPolicy(u.resourceId, &healthcare.SetIamPolicyRequest{
		Policy: healthcarePolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *HealthcareHl7V2StoreIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *HealthcareHl7V2StoreIamUpdater) GetMutexKey() string {
	return u.resourceId
}

func (u *HealthcareHl7V2StoreIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Healthcare Hl7V2Store %q", u.resourceId)
}
