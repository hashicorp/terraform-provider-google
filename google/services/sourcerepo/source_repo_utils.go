// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sourcerepo

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/services/pubsub"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func expandSourceRepoRepositoryPubsubConfigsTopic(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	// short-circuit if the topic is a full uri so we don't need to GetProject
	ok, err := regexp.MatchString(pubsub.PubsubTopicRegex, v.(string))
	if err != nil {
		return "", err
	}

	if ok {
		return v.(string), nil
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return "", err
	}

	return pubsub.GetComputedTopicName(project, v.(string)), err
}
