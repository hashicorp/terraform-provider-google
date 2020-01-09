package google

import (
	"fmt"
)

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
