// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datalossprevention

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepers("DataLossPreventionDiscoveryConfig", testSweepDataLossPreventionDiscoveryConfig)
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepDataLossPreventionDiscoveryConfig(region string) error {
	resourceName := "DataLossPreventionDiscoveryConfig"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	// Setup variables to replace in list template
	d := &tpgresource.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":  config.Project,
			"region":   region,
			"location": region,
			"zone":     "-",
		},
	}

	listTemplate := strings.Split("https://dlp.googleapis.com/v2/projects/{{project}}/locations/{{location}}/discoveryConfigs", "?")[0]
	listUrl, err := tpgresource.ReplaceVars(d, config, listTemplate)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error preparing sweeper list url: %s", err)
		return nil
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    listUrl,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", listUrl, err)
		return nil
	}

	resourceList, ok := res["discoveryConfigs"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["name"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource name was nil", resourceName)
			return nil
		}

		// Note that we do not check for a sweepable prefix here.
		// We can have at most 1 DiscoveryConfig for a storage type in the same project/location, so ensure we delete everything.
		name := tpgresource.GetResourceNameFromSelfLink(obj["name"].(string))

		deleteTemplate := "https://dlp.googleapis.com/v2/projects/{{project}}/locations/{{location}}/discoveryConfigs/{{name}}"
		deleteUrl, err := tpgresource.ReplaceVars(d, config, deleteTemplate)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error preparing delete url: %s", err)
			return nil
		}
		deleteUrl = deleteUrl + name

		// Don't wait on operations as we may have a lot to delete
		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "DELETE",
			Project:   config.Project,
			RawURL:    deleteUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, name)
		}
	}
	return nil
}
