package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamKmsKeyRingSchema = map[string]*schema.Schema{
	"key_ring_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type KmsKeyRingIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewKmsKeyRingIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	keyRing := d.Get("key_ring_id").(string)
	keyRingId, err := parseKmsKeyRingId(keyRing, config)

	if err != nil {
		return nil, fmt.Errorf("Error parsing resource ID for for %s: %s", keyRing, err)
	}

	return &KmsKeyRingIamUpdater{
		resourceId: keyRingId.keyRingId(),
		Config:     config,
	}, nil
}

func (u *KmsKeyRingIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientKms.Projects.Locations.KeyRings.GetIamPolicy(u.resourceId).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for %s: %s", u.DescribeResource(), err)
	}

	cloudResourcePolicy := &cloudresourcemanager.Policy{}

	err = Convert(p, cloudResourcePolicy)

	if err != nil {
		return nil, fmt.Errorf("Invalid IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return cloudResourcePolicy, nil
}

func (u *KmsKeyRingIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	kmsPolicy := &cloudkms.Policy{}
	err := Convert(policy, kmsPolicy)

	if err != nil {
		return fmt.Errorf("Invalid IAM policy for %s: %s", u.DescribeResource(), err)
	}

	_, err = u.Config.clientKms.Projects.Locations.KeyRings.SetIamPolicy(u.resourceId, &cloudkms.SetIamPolicyRequest{
		Policy: kmsPolicy,
	}).Do()

	if err != nil {
		return fmt.Errorf("Error setting IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return nil
}

func (u *KmsKeyRingIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *KmsKeyRingIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-kms-key-ring-%s", u.resourceId)
}

func (u *KmsKeyRingIamUpdater) DescribeResource() string {
	return fmt.Sprintf("KMS KeyRing %q", u.resourceId)
}
