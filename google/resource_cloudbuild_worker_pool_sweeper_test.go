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

	cloudbuild "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("cloudbuildWorker_pool", &resource.Sweeper{
		Name: "cloudbuildWorker_pool",
		F:    testSweepcloudbuildWorker_pool,
	})
}

func testSweepcloudbuildWorker_pool(region string) error {
	log.Print("[INFO][SWEEPER_LOG] Starting sweeper for cloudbuildWorker_pool")

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

	t := &testing.T{}
	billingId := getTestBillingAccountFromEnv(t)

	// Setup variables to be used for Delete arguments.
	d := map[string]string{
		"project":         config.Project,
		"region":          region,
		"location":        region,
		"zone":            "-",
		"billing_account": billingId,
	}

	client := NewDCLcloudbuildClient(config, config.userAgent, "", 0)
	err = client.DeleteAllWorkerPool(context.Background(), d["project"], d["location"], isDeletablecloudbuildWorker_pool)
	if err != nil {
		return err
	}
	return nil
}

func isDeletablecloudbuildWorker_pool(r *cloudbuild.WorkerPool) bool {
	return isSweepableTestResource(*r.Name)
}
