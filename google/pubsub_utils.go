package google

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func comparePubsubSubscriptionBasename(_, old, new string, _ *schema.ResourceData) bool {
	oldStripped, err := getPubsubSubscriptionBasename(old)
	if err != nil {
		return false
	}

	newStripped, err := getPubsubSubscriptionBasename(new)
	if err != nil {
		return false
	}

	if oldStripped == newStripped {
		return true
	}

	return false
}

func comparePubsubTopicBasename(_, old, new string, _ *schema.ResourceData) bool {
	oldStripped, err := getPubsubTopicBasename(old)
	if err != nil {
		return false
	}

	newStripped, err := getPubsubTopicBasename(new)
	if err != nil {
		return false
	}

	if oldStripped == newStripped {
		return true
	}

	return false
}

func getComputedSubscriptionName(project, subscription string) string {
	match, _ := regexp.MatchString("projects\\/.*\\/subscriptions\\/.*", subscription)
	if match {
		return subscription
	}
	return fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription)
}

func getComputedTopicName(project, topic string) string {
	match, _ := regexp.MatchString("projects\\/.*\\/topics\\/.*", topic)
	if match {
		return topic
	}
	return fmt.Sprintf("projects/%s/topics/%s", project, topic)
}

func getPubsubSubscriptionBasename(subscription string) (string, error) {
	re, err := regexp.Compile("projects\\/(.*)\\/subscriptions\\/(.*)")
	if err != nil {
		return "", err
	}
	parts := re.FindStringSubmatch(subscription)
	if len(parts) == 0 {
		return subscription, nil
	} else {
		return parts[2], nil
	}
}

func getPubsubTopicBasename(topic string) (string, error) {
	re, err := regexp.Compile("projects\\/(.*)\\/topics\\/(.*)")
	if err != nil {
		return "", err
	}
	parts := re.FindStringSubmatch(topic)
	if len(parts) == 0 {
		return topic, nil
	} else {
		return parts[2], nil
	}
}
