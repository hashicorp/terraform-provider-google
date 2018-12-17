package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
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
		Config:       config,
	}, nil
}

func PubsubSubscriptionIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("subscription", d.Id())
	return nil
}

func (u *PubsubSubscriptionIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientPubsub.Projects.Subscriptions.GetIamPolicy(u.subscription).Do()

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
	pubsubPolicy, err := resourceManagerToPubsubPolicy(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.clientPubsub.Projects.Subscriptions.SetIamPolicy(u.subscription, &pubsub.SetIamPolicyRequest{
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
