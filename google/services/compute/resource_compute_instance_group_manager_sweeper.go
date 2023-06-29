// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func init() {
	sweeper.AddTestSweepers("ComputeInstanceGroupManager", testSweepComputeInstanceGroupManager)
}

// At the time of writing, the CI only passes us-central1 as the region.
// Since we can read all instances across zones, we don't really use this param.
func testSweepComputeInstanceGroupManager(region string) error {
	resourceName := "ComputeInstanceGroupManager"
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

	found, err := config.NewComputeClient(config.UserAgent).InstanceGroupManagers.AggregatedList(config.Project).Do()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request: %s", err)
		return nil
	}

	// Keep count of items that aren't sweepable for logging.
	nonPrefixCount := 0
	for zone, itemList := range found.Items {
		for _, igm := range itemList.InstanceGroupManagers {
			if !sweeper.IsSweepableTestResource(igm.Name) {
				nonPrefixCount++
				continue
			}

			// Don't wait on operations as we may have a lot to delete
			_, err := config.NewComputeClient(config.UserAgent).InstanceGroupManagers.Delete(config.Project, tpgresource.GetResourceNameFromSelfLink(zone), igm.Name).Do()
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error deleting %s resource %s : %s", resourceName, igm.Name, err)
			} else {
				log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, igm.Name)
			}
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}
