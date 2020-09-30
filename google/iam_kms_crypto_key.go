package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	d          *schema.ResourceData
	Config     *Config
}

func NewKmsCryptoKeyIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	cryptoKey := d.Get("crypto_key_id").(string)
	cryptoKeyId, err := parseKmsCryptoKeyId(cryptoKey, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", cryptoKey), err)
	}

	return &KmsCryptoKeyIamUpdater{
		resourceId: cryptoKeyId.cryptoKeyId(),
		d:          d,
		Config:     config,
	}, nil
}

func CryptoIdParseFunc(d *schema.ResourceData, config *Config) error {
	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}
	if err := d.Set("crypto_key_id", cryptoKeyId.cryptoKeyId()); err != nil {
		return fmt.Errorf("Error setting crypto_key_id: %s", err)
	}
	d.SetId(cryptoKeyId.cryptoKeyId())
	return nil
}

func (u *KmsCryptoKeyIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(u.resourceId).OptionsRequestedPolicyVersion(iamPolicyVersion).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := kmsToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *KmsCryptoKeyIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	kmsPolicy, err := resourceManagerToKmsPolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	_, err = u.Config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys.SetIamPolicy(u.resourceId, &cloudkms.SetIamPolicyRequest{
		Policy: kmsPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
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
