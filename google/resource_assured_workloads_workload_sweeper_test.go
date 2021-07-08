package google

import (
	"context"
	"log"
	"testing"

	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("AssuredWorkloadsWorkload", &resource.Sweeper{
		Name: "AssuredWorkloadsWorkload",
		F:    testSweepAssuredWorkloadsWorkload,
	})
}

func testSweepAssuredWorkloadsWorkload(region string) error {
	resourceName := "AssuredWorkloadsWorkload"
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

	client := NewDCLAssuredWorkloadsClient(config, config.userAgent, "")
	err = client.DeleteAllWorkload(context.Background(), d["organization"], d["location"], isDeletableAssuredWorkloadsWorkload)
	if err != nil {
		return err
	}
	return nil
}

func isDeletableAssuredWorkloadsWorkload(r *assuredworkloads.Workload) bool {
	return isSweepableTestResource(*r.Name)
}
