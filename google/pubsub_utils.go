// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/pubsub"
)

const PubsubTopicRegex = pubsub.PubsubTopicRegex

func getComputedSubscriptionName(project, subscription string) string {
	return pubsub.GetComputedSubscriptionName(project, subscription)
}

func getComputedTopicName(project, topic string) string {
	return pubsub.GetComputedTopicName(project, topic)
}
