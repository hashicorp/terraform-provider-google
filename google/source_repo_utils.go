package google

import "regexp"

func expandSourceRepoRepositoryPubsubConfigsTopic(v interface{}, d TerraformResourceData, config *Config) (string, error) {
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
