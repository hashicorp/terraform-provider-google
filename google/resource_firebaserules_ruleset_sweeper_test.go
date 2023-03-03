// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"log"
	"testing"

	firebaserules "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/firebaserules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("FirebaserulesRuleset", &resource.Sweeper{
		Name: "FirebaserulesRuleset",
		F:    testSweepFirebaserulesRuleset,
	})
}

func testSweepFirebaserulesRuleset(region string) error {
	log.Print("[INFO][SWEEPER_LOG] Starting sweeper for FirebaserulesRuleset")

	config, err := SharedConfigForRegion(region)
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
	billingId := GetTestBillingAccountFromEnv(t)

	// Setup variables to be used for Delete arguments.
	d := map[string]string{
		"project":         config.Project,
		"region":          region,
		"location":        region,
		"zone":            "-",
		"billing_account": billingId,
	}

	client := NewDCLFirebaserulesClient(config, config.UserAgent, "", 0)
	err = client.DeleteAllRuleset(context.Background(), d["project"], isDeletableFirebaserulesRuleset)
	if err != nil {
		return err
	}
	return nil
}

func isDeletableFirebaserulesRuleset(r *firebaserules.Ruleset) bool {
	return IsSweepableTestResource(*r.Name)
}
