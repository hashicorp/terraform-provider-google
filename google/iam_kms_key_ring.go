package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	d          *schema.ResourceData
	Config     *Config
}

func NewKmsKeyRingIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	keyRing := d.Get("key_ring_id").(string)
	keyRingId, err := parseKmsKeyRingId(keyRing, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", keyRing), err)
	}

	return &KmsKeyRingIamUpdater{
		resourceId: keyRingId.keyRingId(),
		d:          d,
		Config:     config,
	}, nil
}

func KeyRingIdParseFunc(d *schema.ResourceData, config *Config) error {
	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return err
	}

	if err := d.Set("key_ring_id", keyRingId.keyRingId()); err != nil {
		return fmt.Errorf("Error setting key_ring_id: %s", err)
	}
	d.SetId(keyRingId.keyRingId())
	return nil
}

func (u *KmsKeyRingIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewKmsClient(userAgent).Projects.Locations.KeyRings.GetIamPolicy(u.resourceId).OptionsRequestedPolicyVersion(iamPolicyVersion).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := kmsToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *KmsKeyRingIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	kmsPolicy, err := resourceManagerToKmsPolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewKmsClient(userAgent).Projects.Locations.KeyRings.SetIamPolicy(u.resourceId, &cloudkms.SetIamPolicyRequest{
		Policy: kmsPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
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

func resourceManagerToKmsPolicy(p *cloudresourcemanager.Policy) (*cloudkms.Policy, error) {
	out := &cloudkms.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a kms policy: {{err}}", err)
	}
	return out, nil
}

func kmsToResourceManagerPolicy(p *cloudkms.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a kms policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
