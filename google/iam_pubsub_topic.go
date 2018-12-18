package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/pubsub/v1"
)

var IamPubsubTopicSchema = map[string]*schema.Schema{
	"topic": {
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

type PubsubTopicIamUpdater struct {
	topic  string
	Config *Config
}

func NewPubsubTopicIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	topic := getComputedTopicName(project, d.Get("topic").(string))

	return &PubsubTopicIamUpdater{
		topic:  topic,
		Config: config,
	}, nil
}

func PubsubTopicIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("topic", d.Id())
	return nil
}

func (u *PubsubTopicIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientPubsub.Projects.Topics.GetIamPolicy(u.topic).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	v1Policy, err := pubsubToResourceManagerPolicy(p)
	if err != nil {
		return nil, err
	}

	return v1Policy, nil
}

func (u *PubsubTopicIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	pubsubPolicy, err := resourceManagerToPubsubPolicy(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.clientPubsub.Projects.Topics.SetIamPolicy(u.topic, &pubsub.SetIamPolicyRequest{
		Policy: pubsubPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *PubsubTopicIamUpdater) GetResourceId() string {
	return u.topic
}

func (u *PubsubTopicIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-pubsub-topic-%s", u.topic)
}

func (u *PubsubTopicIamUpdater) DescribeResource() string {
	return fmt.Sprintf("pubsub topic %q", u.topic)
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
