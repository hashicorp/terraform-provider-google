// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

// This will sweep Compute Instance Templates
func init() {
	sweeper.AddTestSweepers("ComputeInstanceTemplate", testSweepComputeInstanceTemplate)
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepComputeInstanceTemplate(region string) error {
	resourceName := "ComputeInstanceTemplate"
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

	instanceTemplates, err := config.NewComputeClient(config.UserAgent).InstanceTemplates.List(config.Project).Do()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request instance templates LIST: %s", err)
		return nil
	}

	numTemplates := len(instanceTemplates.Items)
	if numTemplates == 0 {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", numTemplates, resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	for _, instanceTemplate := range instanceTemplates.Items {
		// Increment count and skip if resource is not sweepable.
		if !sweeper.IsSweepableTestResource(instanceTemplate.Name) {
			nonPrefixCount++
			continue
		}

		// Don't wait on operations as we may have a lot to delete
		_, err := config.NewComputeClient(config.UserAgent).InstanceTemplates.Delete(config.Project, instanceTemplate.Name).Do()
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting instance template: %s", instanceTemplate.Name)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, instanceTemplate.Name)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf-test prefix remain.", nonPrefixCount)
	}

	return nil
}
