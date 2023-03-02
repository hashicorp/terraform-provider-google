package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PerInstanceConfig needs both regular operation polling AND custom polling for deletion which is why this is not generated
func resourceComputePerInstanceConfigPollRead(d *schema.ResourceData, meta interface{}) PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.UserAgent)
		if err != nil {
			return nil, err
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/listPerInstanceConfigs")
		if err != nil {
			return nil, err
		}

		project, err := getProject(d, config)
		if err != nil {
			return nil, err
		}
		res, err := SendRequest(config, "POST", project, url, userAgent, nil)
		if err != nil {
			return res, err
		}
		res, err = flattenNestedComputePerInstanceConfig(d, meta, res)
		if err != nil {
			return nil, err
		}

		// Returns nil res if nested object is not found
		return res, nil
	}
}

// RegionPerInstanceConfig needs both regular operation polling AND custom polling for deletion which is why this is not generated
func resourceComputeRegionPerInstanceConfigPollRead(d *schema.ResourceData, meta interface{}) PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.UserAgent)
		if err != nil {
			return nil, err
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{region_instance_group_manager}}/listPerInstanceConfigs")
		if err != nil {
			return nil, err
		}

		project, err := getProject(d, config)
		if err != nil {
			return nil, err
		}
		res, err := SendRequest(config, "POST", project, url, userAgent, nil)
		if err != nil {
			return res, err
		}
		res, err = flattenNestedComputeRegionPerInstanceConfig(d, meta, res)
		if err != nil {
			return nil, err
		}

		// Returns nil res if nested object is not found
		return res, nil
	}
}

// Returns an instance name in the form zones/{zone}/instances/{instance} for the managed
// instance matching the name of a PerInstanceConfig
func findInstanceName(d *schema.ResourceData, config *Config) (string, error) {
	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{region_instance_group_manager}}/listManagedInstances")
	if err != nil {
		return "", err
	}

	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return "", err
	}

	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}
	instanceNameToFind := fmt.Sprintf("/%s", d.Get("name").(string))

	token := ""
	for paginate := true; paginate; {
		urlWithToken := ""
		if token != "" {
			urlWithToken = fmt.Sprintf("%s?maxResults=1&pageToken=%s", url, token)
		} else {
			urlWithToken = fmt.Sprintf("%s?maxResults=1", url)
		}
		res, err := SendRequest(config, "POST", project, urlWithToken, userAgent, nil)
		if err != nil {
			return "", err
		}

		managedInstances, ok := res["managedInstances"]
		if !ok {
			return "", fmt.Errorf("Failed to parse response for listManagedInstances for %s", d.Id())
		}

		managedInstancesArr := managedInstances.([]interface{})
		for _, managedInstanceRaw := range managedInstancesArr {
			instance := managedInstanceRaw.(map[string]interface{})
			name, ok := instance["instance"]
			if !ok {
				return "", fmt.Errorf("Failed to read instance name for managed instance: %#v", instance)
			}
			if strings.HasSuffix(name.(string), instanceNameToFind) {
				return name.(string), nil
			}
		}

		tokenRaw, paginate := res["nextPageToken"]
		if paginate {
			token = tokenRaw.(string)
		}
	}

	return "", fmt.Errorf("Failed to find managed instance with name: %s", instanceNameToFind)
}

func PollCheckInstanceConfigDeleted(resp map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		return ErrorPollResult(respErr)
	}

	// Nested object 404 appears as nil response
	if resp == nil {
		// Config no longer exists
		return SuccessPollResult()
	}

	// Read status
	status := resp["status"].(string)
	if status == "DELETING" {
		return PendingStatusPollResult("Still deleting")
	}
	return ErrorPollResult(fmt.Errorf("Expected PerInstanceConfig to be deleting but status is: %s", status))
}
