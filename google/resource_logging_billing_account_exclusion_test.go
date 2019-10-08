package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingBillingAccountExclusion_basic(t *testing.T) {
	t.Parallel()

	billingAccount := getTestBillingAccountFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	var exclusion logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingBillingAccountExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountExclusion_basic(exclusionName, description, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountExclusionExists("google_logging_billing_account_exclusion.basic", &exclusion),
					testAccCheckLoggingBillingAccountExclusion(&exclusion, "google_logging_billing_account_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingBillingAccountExclusion_update(t *testing.T) {
	t.Parallel()

	billingAccount := getTestBillingAccountFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	descriptionBefore := "Basic BillingAccount Logging Exclusion" + acctest.RandString(10)
	descriptionAfter := "Updated Basic BillingAccount Logging Exclusion" + acctest.RandString(10)

	var exclusionBefore, exclusionAfter logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingBillingAccountExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountExclusion_basic(exclusionName, descriptionBefore, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountExclusionExists("google_logging_billing_account_exclusion.basic", &exclusionBefore),
					testAccCheckLoggingBillingAccountExclusion(&exclusionBefore, "google_logging_billing_account_exclusion.basic"),
				),
			},
			{
				Config: testAccLoggingBillingAccountExclusion_basic(exclusionName, descriptionAfter, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountExclusionExists("google_logging_billing_account_exclusion.basic", &exclusionAfter),
					testAccCheckLoggingBillingAccountExclusion(&exclusionAfter, "google_logging_billing_account_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_billing_account_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	// Description should have changed, but Filter and Disabled should be the same
	if exclusionBefore.Description == exclusionAfter.Description {
		t.Errorf("Expected Description to change, but it didn't: Description = %#v", exclusionBefore.Description)
	}
	if exclusionBefore.Filter != exclusionAfter.Filter {
		t.Errorf("Expected Filter to be the same, but it differs: before = %#v, after = %#v",
			exclusionBefore.Filter, exclusionAfter.Filter)
	}
	if exclusionBefore.Disabled != exclusionAfter.Disabled {
		t.Errorf("Expected Disabled to be the same, but it differs: before = %#v, after = %#v",
			exclusionBefore.Disabled, exclusionAfter.Disabled)
	}
}

func testAccCheckLoggingBillingAccountExclusionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_billing_account_exclusion" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.BillingAccounts.Exclusions.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("billingAccount exclusion still exists")
		}
	}

	return nil
}

func testAccCheckLoggingBillingAccountExclusionExists(n string, exclusion *logging.LogExclusion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.BillingAccounts.Exclusions.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*exclusion = *si

		return nil
	}
}

func testAccCheckLoggingBillingAccountExclusion(exclusion *logging.LogExclusion, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		if exclusion.Description != attributes["description"] {
			return fmt.Errorf("mismatch on description: api has %s but client has %s", exclusion.Description, attributes["description"])
		}

		if exclusion.Filter != attributes["filter"] {
			return fmt.Errorf("mismatch on filter: api has %s but client has %s", exclusion.Filter, attributes["filter"])
		}

		disabledAttribute, err := toBool(attributes["disabled"])
		if err != nil {
			return err
		}
		if exclusion.Disabled != disabledAttribute {
			return fmt.Errorf("mismatch on disabled: api has %t but client has %t", exclusion.Disabled, disabledAttribute)
		}

		return nil
	}
}

func testAccLoggingBillingAccountExclusion_basic(exclusionName, description, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_exclusion" "basic" {
	name             = "%s"
	billing_account  = "%s"
	description      = "%s"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, exclusionName, billingAccount, description, getTestProjectFromEnv())
}
