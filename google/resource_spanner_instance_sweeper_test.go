package google

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("SpannerInstance", &resource.Sweeper{
		Name: "SpannerInstance",
		F:    testSweepSpannerInstance,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepSpannerInstance(region string) error {
	resourceName := "SpannerInstance"
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

	spannerUrl := "https://spanner.googleapis.com/v1"
	listUrl := spannerUrl + "/projects/" + config.Project + "/instances"
	res, err := sendRequest(config, "GET", config.Project, listUrl, config.userAgent, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", listUrl, err)
		return nil
	}

	resourceList, ok := res["instances"]
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
		if obj["name"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource name was nil", resourceName)
			return nil
		}

		name := obj["name"].(string)
		shortName := name[strings.LastIndex(name, "/")+1:]

		// Increment count and skip if resource is not sweepable.
		if !isSweepableTestResource(shortName) {
			nonPrefixCount++
			continue
		}

		deleteUrl := spannerUrl + "/" + name
		// Don't wait on operations as we may have a lot to delete
		_, err = sendRequest(config, "DELETE", config.Project, deleteUrl, config.userAgent, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, shortName)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf_test prefix remain.", nonPrefixCount)
	}

	return nil
}
