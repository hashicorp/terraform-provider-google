package google

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	resource.AddTestSweepers("FirebaseAndroidApp", &resource.Sweeper{
		Name: "FirebaseAndroidApp",
		F:    testSweepFirebaseAndroidApp,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepFirebaseAndroidApp(region string) error {
	resourceName := "FirebaseAndroidApp"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := acctest.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	t := &testing.T{}
	billingId := acctest.GetTestBillingAccountFromEnv(t)

	// Setup variables to replace in list template
	d := &acctest.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":         config.Project,
			"region":          region,
			"location":        region,
			"zone":            "-",
			"billing_account": billingId,
		},
	}

	listTemplate := strings.Split("https://firebase.googleapis.com/v1beta1/projects/{{project}}/androidApps", "?")[0]
	listUrl, err := ReplaceVars(d, config, listTemplate)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error preparing sweeper list url: %s", err)
		return nil
	}

	res, err := transport_tpg.SendRequest(config, "GET", config.Project, listUrl, config.UserAgent, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", listUrl, err)
		return nil
	}

	resourceList, ok := res["apps"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// Keep count of items that aren't sweepable for logging.
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["displayName"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource name was nil", resourceName)
			return nil
		}

		// Skip resources that shouldn't be sweeped
		if !acctest.IsSweepableTestResource(obj["displayName"].(string)) {
			nonPrefixCount++
			continue
		}

		name := obj["name"].(string)
		deleteUrl := fmt.Sprintf("https://firebase.googleapis.com/v1beta1/%s:remove", name)

		body := make(map[string]interface{})
		body["immediate"] = true

		// Don't wait on operations as we may have a lot to delete
		_, err = transport_tpg.SendRequest(config, "POST", config.Project, deleteUrl, config.UserAgent, body)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, name)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}
