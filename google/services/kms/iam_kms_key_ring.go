// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
	d          tpgresource.TerraformResourceData
	Config     *transport_tpg.Config
}

func NewKmsKeyRingIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	keyRing := d.Get("key_ring_id").(string)
	keyRingId, err := parseKmsKeyRingId(keyRing, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for %s: {{err}}", keyRing), err)
	}

	return &KmsKeyRingIamUpdater{
		resourceId: keyRingId.KeyRingId(),
		d:          d,
		Config:     config,
	}, nil
}

func KeyRingIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return err
	}

	if err := d.Set("key_ring_id", keyRingId.KeyRingId()); err != nil {
		return fmt.Errorf("Error setting key_ring_id: %s", err)
	}
	d.SetId(keyRingId.KeyRingId())
	return nil
}

func (u *KmsKeyRingIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewKmsClient(userAgent).Projects.Locations.KeyRings.GetIamPolicy(u.resourceId).OptionsRequestedPolicyVersion(tpgiamresource.IamPolicyVersion).Do()

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

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
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
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a kms policy: {{err}}", err)
	}
	return out, nil
}

func kmsToResourceManagerPolicy(p *cloudkms.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a kms policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
