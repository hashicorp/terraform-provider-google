// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepers("AlloydbCluster", testSweepAlloydbCluster)
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepAlloydbCluster(region string) error {
	resourceName := "AlloydbCluster"
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

	t := &testing.T{}
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	// Setup variables to replace in list template
	d := &tpgresource.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":         config.Project,
			"region":          region,
			"location":        region,
			"zone":            "-",
			"billing_account": billingId,
		},
	}

	// manual patch: use aggregated list instead of sweeper-specific location. This will clear secondary clusters.
	listTemplate := strings.Split("https://alloydb.googleapis.com/v1/projects/{{project}}/locations/-/clusters", "?")[0]
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

	resourceList, ok := res["clusters"]
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

		// manual patch: use raw name for url instead of constructing it, so that resource locations are supplied through aggregated list
		// manual patch: Using the force=true ensures that we delete instances as well.
		name := obj["name"].(string)
		shortname := tpgresource.GetResourceNameFromSelfLink(name)
		// Skip resources that shouldn't be sweeped
		if !sweeper.IsSweepableTestResource(shortname) {
			nonPrefixCount++
			continue
		}

		deleteTemplate := "https://alloydb.googleapis.com/v1/%s?force=true"
		deleteUrl := fmt.Sprintf(deleteTemplate, name)

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
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", name, shortname)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}
