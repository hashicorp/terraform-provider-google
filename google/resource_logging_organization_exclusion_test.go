package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingOrganizationExclusion_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	var exclusion logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationExclusion_basic(exclusionName, description, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationExclusionExists("google_logging_organization_exclusion.basic", &exclusion),
					testAccCheckLoggingOrganizationExclusion(&exclusion, "google_logging_organization_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingOrganizationExclusion_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	descriptionBefore := "Basic Organization Logging Exclusion" + acctest.RandString(10)
	descriptionAfter := "Updated Basic Organization Logging Exclusion" + acctest.RandString(10)

	var exclusionBefore, exclusionAfter logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationExclusion_basic(exclusionName, descriptionBefore, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationExclusionExists("google_logging_organization_exclusion.basic", &exclusionBefore),
					testAccCheckLoggingOrganizationExclusion(&exclusionBefore, "google_logging_organization_exclusion.basic"),
				),
			},
			{
				Config: testAccLoggingOrganizationExclusion_basic(exclusionName, descriptionAfter, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationExclusionExists("google_logging_organization_exclusion.basic", &exclusionAfter),
					testAccCheckLoggingOrganizationExclusion(&exclusionAfter, "google_logging_organization_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_organization_exclusion.basic",
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

func testAccCheckLoggingOrganizationExclusionExists(n string, exclusion *logging.LogExclusion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.Organizations.Exclusions.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*exclusion = *si

		return nil
	}
}

func testAccCheckLoggingOrganizationExclusion(exclusion *logging.LogExclusion, n string) resource.TestCheckFunc {
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

func testAccLoggingOrganizationExclusion_basic(exclusionName, description, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_exclusion" "basic" {
	name             = "%s"
	org_id           = "%s"
	description      = "%s"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, exclusionName, orgId, description, getTestProjectFromEnv())
}
