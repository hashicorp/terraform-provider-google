package google

import (
	"regexp"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func expandSourceRepoRepositoryPubsubConfigsTopic(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	// short-circuit if the topic is a full uri so we don't need to getProject
	ok, err := regexp.MatchString(PubsubTopicRegex, v.(string))
	if err != nil {
		return "", err
	}

	if ok {
		return v.(string), nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	return getComputedTopicName(project, v.(string)), err
}
