package google

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("StorageBucket", &resource.Sweeper{
		Name: "StorageBucket",
		F:    testSweepStorageBucket,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepStorageBucket(region string) error {
	resourceName := "StorageBucket"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
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

	servicesUrl, err := addQueryParams("https://storage.googleapis.com/storage/v1/b", params)
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", config.Project, servicesUrl, config.userAgent, nil)
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
	for _, ri := range rl {
		obj := ri.(map[string]interface{})

		id := obj["name"].(string)
		// Increment count and skip if resource is not sweepable.
		if !isSweepableTestResource(id) {
			nonPrefixCount++
			continue
		}

		deleteUrl := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", id)
		_, err = sendRequest(config, "DELETE", config.Project, deleteUrl, config.userAgent, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Deleted a %s resource: %s", resourceName, id)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf-test prefix remain.", nonPrefixCount)
	}

	return nil
}
