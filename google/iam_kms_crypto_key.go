package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamKmsCryptoKeySchema = map[string]*schema.Schema{
	"crypto_key_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type KmsCryptoKeyIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewKmsCryptoKeyIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	cryptoKey := d.Get("crypto_key_id").(string)
	cryptoKeyId, err := parseKmsCryptoKeyId(cryptoKey, config)

	if err != nil {
		return nil, fmt.Errorf("Error parsing resource ID for for %s: %s", cryptoKey, err)
	}

	return &KmsCryptoKeyIamUpdater{
		resourceId: cryptoKeyId.cryptoKeyId(),
		Config:     config,
	}, nil
}

func (u *KmsCryptoKeyIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientKms.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(u.resourceId).Do()

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

func (u *KmsCryptoKeyIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	kmsPolicy := &cloudkms.Policy{}
	err := Convert(policy, kmsPolicy)

	if err != nil {
		return fmt.Errorf("Invalid IAM policy for %s: %s", u.DescribeResource(), err)
	}

	_, err = u.Config.clientKms.Projects.Locations.KeyRings.CryptoKeys.SetIamPolicy(u.resourceId, &cloudkms.SetIamPolicyRequest{
		Policy: kmsPolicy,
	}).Do()

	if err != nil {
		return fmt.Errorf("Error setting IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return nil
}

func (u *KmsCryptoKeyIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *KmsCryptoKeyIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-kms-crypto-key-%s", u.resourceId)
}

func (u *KmsCryptoKeyIamUpdater) DescribeResource() string {
	return fmt.Sprintf("KMS CryptoKey %q", u.resourceId)
}
