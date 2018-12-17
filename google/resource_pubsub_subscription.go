package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/pubsub/v1"
	"regexp"
)

func resourcePubsubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubSubscriptionCreate,
		Read:   resourcePubsubSubscriptionRead,
		Update: resourcePubsubSubscriptionUpdate,
		Delete: resourcePubsubSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePubsubSubscriptionStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"topic": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"ack_deadline_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"push_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"push_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourcePubsubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := getComputedSubscriptionName(project, d.Get("name").(string))
	computed_topic_name := getComputedTopicName(project, d.Get("topic").(string))

	//  process optional parameters
	var ackDeadlineSeconds int64
	ackDeadlineSeconds = 10
	if v, ok := d.GetOk("ack_deadline_seconds"); ok {
		ackDeadlineSeconds = int64(v.(int))
	}

	subscription := &pubsub.Subscription{
		AckDeadlineSeconds: ackDeadlineSeconds,
		Topic:              computed_topic_name,
		PushConfig:         expandPubsubSubscriptionPushConfig(d.Get("push_config").([]interface{})),
	}

	call := config.clientPubsub.Projects.Subscriptions.Create(name, subscription)
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return resourcePubsubSubscriptionRead(d, meta)
}

func getComputedTopicName(project, topic string) string {
	match, _ := regexp.MatchString("projects\\/.*\\/topics\\/.*", topic)
	if match {
		return topic
	}
	return fmt.Sprintf("projects/%s/topics/%s", project, topic)
}

func getComputedSubscriptionName(project, subscription string) string {
	match, _ := regexp.MatchString("projects\\/.*\\/subscriptions\\/.*", subscription)
	if match {
		return subscription
	}
	return fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription)
}

func resourcePubsubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Id()
	subscription, err := config.clientPubsub.Projects.Subscriptions.Get(name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Pubsub Subscription %q", name))
	}

	d.Set("name", GetResourceNameFromSelfLink(subscription.Name))
	d.Set("topic", subscription.Topic)
	d.Set("ack_deadline_seconds", subscription.AckDeadlineSeconds)
	d.Set("path", subscription.Name)
	d.Set("push_config", flattenPubsubSubscriptionPushConfig(subscription.PushConfig))
	d.Set("project", project)

	return nil
}

func resourcePubsubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("push_config") {
		_, err := config.clientPubsub.Projects.Subscriptions.ModifyPushConfig(d.Id(), &pubsub.ModifyPushConfigRequest{
			PushConfig: expandPubsubSubscriptionPushConfig(d.Get("push_config").([]interface{})),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating subscription %q: %s", d.Get("name"), err)
		}
	}

	d.Partial(false)

	return nil
}

func resourcePubsubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Subscriptions.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}

func resourcePubsubSubscriptionStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("projects/%s/subscriptions/%s", project, d.Id())

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenPubsubSubscriptionPushConfig(pushConfig *pubsub.PushConfig) []map[string]interface{} {
	configs := make([]map[string]interface{}, 0, 1)

	if pushConfig == nil || len(pushConfig.PushEndpoint) == 0 {
		return configs
	}

	configs = append(configs, map[string]interface{}{
		"push_endpoint": pushConfig.PushEndpoint,
		"attributes":    pushConfig.Attributes,
	})

	return configs
}

func expandPubsubSubscriptionPushConfig(configured []interface{}) *pubsub.PushConfig {
	if len(configured) == 0 || configured[0] == nil {
		// An empty `pushConfig` indicates that the Pub/Sub system should stop pushing messages
		// from the given subscription and allow messages to be pulled and acknowledged.
		return &pubsub.PushConfig{}
	}

	pushConfig := configured[0].(map[string]interface{})
	return &pubsub.PushConfig{
		PushEndpoint: pushConfig["push_endpoint"].(string),
		Attributes:   convertStringMap(pushConfig["attributes"].(map[string]interface{})),
	}
}
