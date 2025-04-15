// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepersLegacy("StorageBucket", testSweepStorageBucket)
}

func disableAnywhereCacheIfAny(config *transport_tpg.Config, bucket string) bool {
	// Define the cache list URL
	cacheListUrl := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s/anywhereCaches/", bucket)

	// Send request to get resource list
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    cacheListUrl,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error fetching caches from url %s: %s", cacheListUrl, err)
		return false
	}

	resourceList, ok := res["items"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] No caches found for %s.", bucket)
		return true
	}

	rl := resourceList.([]interface{})

	// Iterate over each object in the resource list
	for _, item := range rl {
		// Ensure the item is a map
		obj := item.(map[string]interface{})

		// Check the state of the object
		state := obj["state"].(string)
		if state != "running" && state != "paused" {
			continue
		}

		// Disable the cache if state is running or paused
		disableUrl := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s/anywhereCaches/%s/disable", obj["bucket"], obj["anywhereCacheId"])
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   config.Project,
			RawURL:    disableUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error disabling cache: %s", err)
		}
	}

	// Return true if no items were found, otherwise false
	return len(rl) == 0
}

func deleteObjects(config *transport_tpg.Config, bucket string) bool {
	userAgent := config.UserAgent
	var listError, deleteObjectError error
	for deleteObjectError == nil {
		res, err := config.NewStorageClient(userAgent).Objects.List(bucket).Versions(true).Do()
		if err != nil {
			log.Printf("Error listing contents of bucket %s: %v", bucket, err)
			listError = err
			break
		}
		if len(res.Items) == 0 {
			break
		}
		wp := workerpool.New(runtime.NumCPU() - 1)

		for _, object := range res.Items {
			log.Printf("[DEBUG] Found %s", object.Name)
			object := object

			wp.Submit(func() {
				log.Printf("[TRACE] Attempting to delete %s", object.Name)
				if err := config.NewStorageClient(userAgent).Objects.Delete(bucket, object.Name).Generation(object.Generation).Do(); err != nil {
					deleteObjectError = err
					log.Printf("[ERR] Failed to delete storage object %s: %s", object.Name, err)
				} else {
					log.Printf("[TRACE] Successfully deleted %s", object.Name)
				}
			})
		}
		wp.StopWait()
	}

	if listError != nil || deleteObjectError != nil {
		return false
	}

	return true
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepStorageBucket(region string) error {
	resourceName := "StorageBucket"
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

	params := map[string]string{
		"project":    config.Project,
		"projection": "noAcl", // returns 1000 items instead of 200
	}

	servicesUrl, err := transport_tpg.AddQueryParams("https://storage.googleapis.com/storage/v1/b", params)
	if err != nil {
		return err
	}

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

	resourceList, ok := res["items"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	bucketWithCaches := 0
	bucketWithObjects := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})

		id := obj["name"].(string)
		// Increment count and skip if resource is not sweepable.
		if !sweeper.IsSweepableTestResource(id) {
			nonPrefixCount++
			continue
		}

		readyToDeleteBucket := disableAnywhereCacheIfAny(config, id)
		if !readyToDeleteBucket {
			log.Printf("[INFO][SWEEPER_LOG] Bucket %s has anywhere caches, requests have been made to backend to disable them, The bucket would be automatically deleted once caches are deleted from bucket", id)
			bucketWithCaches++
			continue
		}

		emptyBucket := deleteObjects(config, id)
		if !emptyBucket {
			log.Printf("[INFO][SWEEPER_LOG] Deleting the objects failed in the bucket %v", id)
			bucketWithObjects++
			continue
		}

		deleteUrl := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", id)
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
			log.Printf("[INFO][SWEEPER_LOG] Deleted a %s resource: %s", resourceName, id)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without valid test prefixes remain.", nonPrefixCount)
	}
	if bucketWithCaches > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items with valid test prefixes remain, and can not be deleted due to their underlying resources", bucketWithCaches)
	}
	if bucketWithObjects > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items with valid test prefixes remain, and can not be deleted as buckets are non empty", bucketWithObjects)
	}

	return nil
}
