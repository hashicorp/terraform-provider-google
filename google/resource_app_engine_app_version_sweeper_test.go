package google

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// This will sweep both Standard and Flexible App Engine App Versions
func init() {
	resource.AddTestSweepers("AppEngineAppVersion", &resource.Sweeper{
		Name: "AppEngineAppVersion",
		F:    testSweepAppEngineAppVersion,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepAppEngineAppVersion(region string) error {
	resourceName := "AppEngineAppVersion"
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

	servicesUrl := "https://appengine.googleapis.com/v1/apps/" + config.Project + "/services"
	res, err := sendRequest(config, "GET", config.Project, servicesUrl, nil)
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
	// items who don't match the tf-test prefix
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["id"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource id was nil", resourceName)
			return nil
		}

		id := obj["id"].(string)
		// Only sweep resources with the test prefix
		if !strings.HasPrefix(id, "tf-test") {
			nonPrefixCount++
			continue
		}

		deleteUrl := servicesUrl + "/" + id
		// Don't wait on operations as we may have a lot to delete
		_, err = sendRequest(config, "DELETE", config.Project, deleteUrl, nil)
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
