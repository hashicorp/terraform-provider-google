package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/pubsub/v1"
)

var IamPubsubSubscriptionSchema = map[string]*schema.Schema{
	"subscription": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareSelfLinkOrResourceName,
	},
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type PubsubSubscriptionIamUpdater struct {
	subscription string
	d            *schema.ResourceData
	Config       *Config
}

func NewPubsubSubscriptionIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	subscription := getComputedSubscriptionName(project, d.Get("subscription").(string))

	return &PubsubSubscriptionIamUpdater{
		subscription: subscription,
		d:            d,
		Config:       config,
	}, nil
}

func PubsubSubscriptionIdParseFunc(d *schema.ResourceData, _ *Config) error {
	if err := d.Set("subscription", d.Id()); err != nil {
		return fmt.Errorf("Error setting subscription: %s", err)
	}
	return nil
}

func (u *PubsubSubscriptionIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewPubsubClient(userAgent).Projects.Subscriptions.GetIamPolicy(u.subscription).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	v1Policy, err := pubsubToResourceManagerPolicy(p)
	if err != nil {
		return nil, err
	}

	return v1Policy, nil
}

func (u *PubsubSubscriptionIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	pubsubPolicy, err := resourceManagerToPubsubPolicy(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.NewPubsubClient(userAgent).Projects.Subscriptions.SetIamPolicy(u.subscription, &pubsub.SetIamPolicyRequest{
		Policy: pubsubPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *PubsubSubscriptionIamUpdater) GetResourceId() string {
	return u.subscription
}

func (u *PubsubSubscriptionIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-pubsub-subscription-%s", u.subscription)
}

func (u *PubsubSubscriptionIamUpdater) DescribeResource() string {
	return fmt.Sprintf("pubsub subscription %q", u.subscription)
}

// v1 and v2beta policy are identical
func resourceManagerToPubsubPolicy(in *cloudresourcemanager.Policy) (*pubsub.Policy, error) {
	out := &pubsub.Policy{}
	err := Convert(in, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a pubsub policy: {{err}}", err)
	}
	return out, nil
}

func pubsubToResourceManagerPolicy(in *pubsub.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(in, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a pubsub policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
