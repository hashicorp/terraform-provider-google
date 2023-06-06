// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// PerInstanceConfig needs both regular operation polling AND custom polling for deletion which is why this is not generated
func resourceComputePerInstanceConfigPollRead(d *schema.ResourceData, meta interface{}) transport_tpg.PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*transport_tpg.Config)
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return nil, err
		}

		url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/listPerInstanceConfigs")
		if err != nil {
			return nil, err
		}

		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return nil, err
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
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
func resourceComputeRegionPerInstanceConfigPollRead(d *schema.ResourceData, meta interface{}) transport_tpg.PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*transport_tpg.Config)
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return nil, err
		}

		url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{region_instance_group_manager}}/listPerInstanceConfigs")
		if err != nil {
			return nil, err
		}

		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return nil, err
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
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
func findInstanceName(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceGroupManagers/{{region_instance_group_manager}}/listManagedInstances")
	if err != nil {
		return "", err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return "", err
	}

	project, err := tpgresource.GetProject(d, config)
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
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    urlWithToken,
			UserAgent: userAgent,
		})
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

func PollCheckInstanceConfigDeleted(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	if respErr != nil {
		return transport_tpg.ErrorPollResult(respErr)
	}

	// Nested object 404 appears as nil response
	if resp == nil {
		// Config no longer exists
		return transport_tpg.SuccessPollResult()
	}

	// Read status
	status := resp["status"].(string)
	if status == "DELETING" {
		return transport_tpg.PendingStatusPollResult("Still deleting")
	}
	return transport_tpg.ErrorPollResult(fmt.Errorf("Expected PerInstanceConfig to be deleting but status is: %s", status))
}
