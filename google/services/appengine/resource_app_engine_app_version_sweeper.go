// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// This will sweep both Standard and Flexible App Engine App Versions
func init() {
	sweeper.AddTestSweepers("AppEngineAppVersion", testSweepAppEngineAppVersion)
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepAppEngineAppVersion(region string) error {
	resourceName := "AppEngineAppVersion"
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

	servicesUrl := "https://appengine.googleapis.com/v1/apps/" + config.Project + "/services"
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    servicesUrl,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", servicesUrl, err)
		return nil
	}

	resourceList, ok := res["services"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["id"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource id was nil", resourceName)
			return nil
		}

		id := obj["id"].(string)
		// Increment count and skip if resource is not sweepable.
		if !sweeper.IsSweepableTestResource(id) {
			nonPrefixCount++
			continue
		}

		deleteUrl := servicesUrl + "/" + id
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
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, id)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf_test prefix remain.", nonPrefixCount)
	}

	return nil
}
