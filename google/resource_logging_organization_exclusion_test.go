package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Logging exclusions don't always work when making parallel requests, so run tests serially
func TestAccLoggingOrganizationExclusion(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"basic":    testAccLoggingOrganizationExclusion_basic,
		"update":   testAccLoggingOrganizationExclusion_update,
		"multiple": testAccLoggingOrganizationExclusion_multiple,
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

func testAccLoggingOrganizationExclusion_basic(t *testing.T) {
	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationExclusion_basicCfg(exclusionName, description, org),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingOrganizationExclusion_update(t *testing.T) {
	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	descriptionBefore := "Basic Organization Logging Exclusion" + acctest.RandString(10)
	descriptionAfter := "Updated Basic Organization Logging Exclusion" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationExclusion_basicCfg(exclusionName, descriptionBefore, org),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingOrganizationExclusion_basicCfg(exclusionName, descriptionAfter, org),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingOrganizationExclusion_multiple(t *testing.T) {
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationExclusion_multipleCfg(org),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic0",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingOrganizationExclusionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_organization_exclusion" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.Organizations.Exclusions.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("organization exclusion still exists")
		}
	}

	return nil
}

func testAccLoggingOrganizationExclusion_basicCfg(exclusionName, description, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_exclusion" "basic" {
  name        = "%s"
  org_id      = "%s"
  description = "%s"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, exclusionName, orgId, description, getTestProjectFromEnv())
}

func testAccLoggingOrganizationExclusion_multipleCfg(orgId string) string {
	s := ""
	for i := 0; i < 3; i++ {
		s += fmt.Sprintf(`
resource "google_logging_organization_exclusion" "basic%d" {
	name             = "%s"
	org_id           = "%s"
	description      = "Basic Organization Logging Exclusion"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, i, "tf-test-exclusion-"+acctest.RandString(10), orgId, getTestProjectFromEnv())
	}
	return s
}
