// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package pubsub

import (
	"fmt"
	"regexp"
)

const PubsubTopicRegex = "projects\\/.*\\/topics\\/.*"

func GetComputedSubscriptionName(project, subscription string) string {
	match, _ := regexp.MatchString("projects\\/.*\\/subscriptions\\/.*", subscription)
	if match {
		return subscription
	}
	return fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription)
}

func GetComputedTopicName(project, topic string) string {
	match, _ := regexp.MatchString(PubsubTopicRegex, topic)
	if match {
		return topic
	}
	return fmt.Sprintf("projects/%s/topics/%s", project, topic)
}
