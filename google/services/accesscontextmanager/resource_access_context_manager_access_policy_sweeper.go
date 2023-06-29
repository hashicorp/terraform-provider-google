// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager

import (
	"context"
	"fmt"
	"log"
	neturl "net/url"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepers("gcp_access_context_manager_policy", testSweepAccessContextManagerPolicies)
}

func testSweepAccessContextManagerPolicies(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Fatalf("error getting shared config for region %q: %s", region, err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading and validating shared config for region %q: %s", region, err)
	}

	testOrg := envvar.GetTestOrgFromEnv(nil)
	if testOrg == "" {
		log.Printf("test org not set for test environment, skip sweep")
		return nil
	}

	log.Printf("[DEBUG] Listing Access Policies for org %q", testOrg)

	parent := neturl.QueryEscape(fmt.Sprintf("organizations/%s", testOrg))
	listUrl := fmt.Sprintf("%saccessPolicies?parent=%s", config.AccessContextManagerBasePath, parent)

	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    listUrl,
		UserAgent: config.UserAgent,
	})
	if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("unable to list AccessPolicies for organization %q: %v", testOrg, err)
		return nil
	}
	var policies []interface{}
	if resp != nil {
		if v, ok := resp["accessPolicies"]; ok {
			policies = v.([]interface{})
		}
	}

	if len(policies) == 0 {
		log.Printf("[DEBUG] no access policies found, exiting sweeper")
		return nil
	}
	if len(policies) > 1 {
		log.Printf("unexpected - more than one access policies found, change the tests")
		return nil
	}

	policy := policies[0].(map[string]interface{})
	log.Printf("[DEBUG] Deleting test Access Policies %q", policy["name"])

	policyUrl := config.AccessContextManagerBasePath + policy["name"].(string)
	if _, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		RawURL:    policyUrl,
		UserAgent: config.UserAgent,
	}); err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("unable to delete access policy %q", policy["name"].(string))
		return nil
	}

	return nil
}
