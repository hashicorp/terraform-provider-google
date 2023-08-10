// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// Logging exclusions don't always work when making parallel requests, so run tests serially
func TestAccLoggingBillingAccountExclusion(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"basic":    testAccLoggingBillingAccountExclusion_basic,
		"update":   testAccLoggingBillingAccountExclusion_update,
		"multiple": testAccLoggingBillingAccountExclusion_multiple,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccLoggingBillingAccountExclusion_basic(t *testing.T) {
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(t, 10)
	description := "Description " + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountExclusionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountExclusion_basicCfg(exclusionName, description, billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingBillingAccountExclusion_update(t *testing.T) {
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(t, 10)
	descriptionBefore := "Basic BillingAccount Logging Exclusion" + acctest.RandString(t, 10)
	descriptionAfter := "Updated Basic BillingAccount Logging Exclusion" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountExclusionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountExclusion_basicCfg(exclusionName, descriptionBefore, billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingBillingAccountExclusion_basicCfg(exclusionName, descriptionAfter, billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingBillingAccountExclusion_multiple(t *testing.T) {
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountExclusionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountExclusion_multipleCfg("tf-test-exclusion-"+acctest.RandString(t, 10), billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic0",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingBillingAccountExclusionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_billing_account_exclusion" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).BillingAccounts.Exclusions.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("billingAccount exclusion still exists")
			}
		}

		return nil
	}
}

func testAccLoggingBillingAccountExclusion_basicCfg(exclusionName, description, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_exclusion" "basic" {
  name            = "%s"
  billing_account = "%s"
  description     = "%s"
  filter          = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, exclusionName, billingAccount, description, envvar.GetTestProjectFromEnv())
}

func testAccLoggingBillingAccountExclusion_multipleCfg(exclusionName, billingAccount string) string {
	s := ""
	for i := 0; i < 3; i++ {
		s += fmt.Sprintf(`
resource "google_logging_billing_account_exclusion" "basic%d" {
	name             = "%s%d"
	billing_account  = "%s"
	description      = "Basic BillingAccount Logging Exclusion"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, i, exclusionName, i, billingAccount, envvar.GetTestProjectFromEnv())
	}
	return s
}
