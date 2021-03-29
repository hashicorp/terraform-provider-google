package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// List of prefixes used for test resource names
var testResourcePrefixes = []string{
	"tf-test",
	"tfgen",
	"gke-us-central1-tf", // composer-created disks which are abandoned by design (https://cloud.google.com/composer/pricing)
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedConfigForRegion returns a common config setup needed for the sweeper
// functions for a given region
func sharedConfigForRegion(region string) (*Config, error) {
	project := getTestProjectFromEnv()
	if project == "" {
		return nil, fmt.Errorf("set project using any of these env variables %v", projectEnvVars)
	}

	creds := getTestCredsFromEnv()
	if creds == "" {
		return nil, fmt.Errorf("set credentials using any of these env variables %v", credsEnvVars)
	}

	conf := &Config{
		Credentials: creds,
		Region:      region,
		Project:     project,
	}

	ConfigureBasePaths(conf)

	return conf, nil
}

func isSweepableTestResource(resourceName string) bool {
	for _, p := range testResourcePrefixes {
		if strings.HasPrefix(resourceName, p) {
			return true
		}
	}
	return false
}
