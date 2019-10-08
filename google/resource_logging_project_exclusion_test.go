package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingProjectExclusion_basic(t *testing.T) {
	t.Parallel()

	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	var exclusion logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basic(exclusionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectExclusionExists("google_logging_project_exclusion.basic", &exclusion),
					testAccCheckLoggingProjectExclusion(&exclusion, "google_logging_project_exclusion.basic")),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectExclusion_disablePreservesFilter(t *testing.T) {
	t.Parallel()

	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	var exclusionBefore, exclusionAfter logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basic(exclusionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectExclusionExists("google_logging_project_exclusion.basic", &exclusionBefore),
					testAccCheckLoggingProjectExclusion(&exclusionBefore, "google_logging_project_exclusion.basic"),
				),
			},
			{
				Config: testAccLoggingProjectExclusion_basicDisabled(exclusionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectExclusionExists("google_logging_project_exclusion.basic", &exclusionAfter),
					testAccCheckLoggingProjectExclusion(&exclusionAfter, "google_logging_project_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	// Description and Disabled should have changed, but Filter should be the same
	if exclusionBefore.Description == exclusionAfter.Description {
		t.Errorf("Expected Description to change, but it didn't: Description = %#v", exclusionBefore.Description)
	}
	if exclusionBefore.Filter != exclusionAfter.Filter {
		t.Errorf("Expected Filter to be the same, but it differs: before = %#v, after = %#v",
			exclusionBefore.Filter, exclusionAfter.Filter)
	}
	if exclusionBefore.Disabled == exclusionAfter.Disabled {
		t.Errorf("Expected Disabled to change, but it didn't: Disabled = %#v", exclusionBefore.Disabled)
	}
}

func TestAccLoggingProjectExclusion_update(t *testing.T) {
	t.Parallel()

	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	var exclusionBefore, exclusionAfter logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basic(exclusionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectExclusionExists("google_logging_project_exclusion.basic", &exclusionBefore),
					testAccCheckLoggingProjectExclusion(&exclusionBefore, "google_logging_project_exclusion.basic"),
				),
			},
			{
				Config: testAccLoggingProjectExclusion_basicUpdated(exclusionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectExclusionExists("google_logging_project_exclusion.basic", &exclusionAfter),
					testAccCheckLoggingProjectExclusion(&exclusionAfter, "google_logging_project_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	// Filter should have changed, but Description and Disabled should be the same
	if exclusionBefore.Description != exclusionAfter.Description {
		t.Errorf("Expected Description to be the same, but it differs: before = %#v, after = %#v",
			exclusionBefore.Description, exclusionAfter.Description)
	}
	if exclusionBefore.Filter == exclusionAfter.Filter {
		t.Errorf("Expected Filter to change, but it didn't: Filter = %#v", exclusionBefore.Filter)
	}
	if exclusionBefore.Disabled != exclusionAfter.Disabled {
		t.Errorf("Expected Disabled to be the same, but it differs: before = %#v, after = %#v",
			exclusionBefore.Disabled, exclusionAfter.Disabled)
	}
}

func testAccCheckLoggingProjectExclusionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_project_exclusion" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.Projects.Exclusions.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("project exclusion still exists")
		}
	}

	return nil
}

func testAccCheckLoggingProjectExclusionExists(n string, exclusion *logging.LogExclusion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.Projects.Exclusions.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*exclusion = *si

		return nil
	}
}

func testAccCheckLoggingProjectExclusion(exclusion *logging.LogExclusion, n string) resource.TestCheckFunc {
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

func testAccLoggingProjectExclusion_basic(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
	name = "%s"
	description = "Basic Project Logging Exclusion"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}`, name, getTestProjectFromEnv())
}

func testAccLoggingProjectExclusion_basicUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
	name = "%s"
	description = "Basic Project Logging Exclusion"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=INFO"
}`, name, getTestProjectFromEnv())
}

func testAccLoggingProjectExclusion_basicDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
	name = "%s"
	description = ""
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
	disabled = true
}`, name, getTestProjectFromEnv())
}
