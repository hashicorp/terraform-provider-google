package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamBillingAccountSchema = map[string]*schema.Schema{
	"billing_account_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type BillingAccountIamUpdater struct {
	billingAccountId string
	Config           *Config
}

func NewBillingAccountIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	return &BillingAccountIamUpdater{
		billingAccountId: canonicalBillingAccountId(d.Get("billing_account_id").(string)),
		Config:           config,
	}, nil
}

func BillingAccountIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("billing_account_id", d.Id())
	return nil
}

func (u *BillingAccountIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	return getBillingAccountIamPolicyByBillingAccountName(u.billingAccountId, u.Config)
}

func (u *BillingAccountIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	billingPolicy, err := resourceManagerToBillingPolicy(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.clientBilling.BillingAccounts.SetIamPolicy("billingAccounts/"+u.billingAccountId, &cloudbilling.SetIamPolicyRequest{
		Policy: billingPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountIamUpdater) GetResourceId() string {
	return u.billingAccountId
}

func (u *BillingAccountIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-billing-account-%s", u.billingAccountId)
}

func (u *BillingAccountIamUpdater) DescribeResource() string {
	return fmt.Sprintf("billingAccount %q", u.billingAccountId)
}

func canonicalBillingAccountId(resource string) string {
	return resource
}

func resourceManagerToBillingPolicy(p *cloudresourcemanager.Policy) (*cloudbilling.Policy, error) {
	out := &cloudbilling.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a billing policy: {{err}}", err)
	}
	return out, nil
}

func billingToResourceManagerPolicy(p *cloudbilling.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a billing policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}

// Retrieve the existing IAM Policy for a billing account
func getBillingAccountIamPolicyByBillingAccountName(resource string, config *Config) (*cloudresourcemanager.Policy, error) {
	p, err := config.clientBilling.BillingAccounts.GetIamPolicy("billingAccounts/" + resource).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for billing account %q: {{err}}", resource), err)
	}

	v1Policy, err := billingToResourceManagerPolicy(p)
	if err != nil {
		return nil, err
	}

	return v1Policy, nil
}
