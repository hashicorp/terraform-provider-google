package google

import (
	"fmt"
	"regexp"
)

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
