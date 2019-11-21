package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Logging exclusions don't always work when making parallel requests, so run tests serially
func TestAccLoggingProjectExclusion(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"basic":                  testAccLoggingProjectExclusion_basic,
		"disablePreservesFilter": testAccLoggingProjectExclusion_disablePreservesFilter,
		"update":                 testAccLoggingProjectExclusion_update,
		"multiple":               testAccLoggingProjectExclusion_multiple,
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

func testAccLoggingProjectExclusion_basic(t *testing.T) {
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basicCfg(exclusionName),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingProjectExclusion_disablePreservesFilter(t *testing.T) {
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basicCfg(exclusionName),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectExclusion_basicDisabled(exclusionName),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingProjectExclusion_update(t *testing.T) {
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_basicCfg(exclusionName),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectExclusion_basicUpdated(exclusionName),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingProjectExclusion_multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectExclusion_multipleCfg(),
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic0",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_project_exclusion.basic2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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
			return fmt.Errorf("project exclusion %s still exists", attributes["id"])
		}
	}

	return nil
}

func testAccLoggingProjectExclusion_basicCfg(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
  name        = "%s"
  description = "Basic Project Logging Exclusion"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, name, getTestProjectFromEnv())
}

func testAccLoggingProjectExclusion_basicUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
  name        = "%s"
  description = "Basic Project Logging Exclusion"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=INFO"
}
`, name, getTestProjectFromEnv())
}

func testAccLoggingProjectExclusion_basicDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic" {
  name        = "%s"
  description = ""
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = true
}
`, name, getTestProjectFromEnv())
}

func testAccLoggingProjectExclusion_multipleCfg() string {
	s := ""
	for i := 0; i < 3; i++ {
		s += fmt.Sprintf(`
resource "google_logging_project_exclusion" "basic%d" {
	name = "%s"
	description = "Basic Project Logging Exclusion"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, i, "tf-test-exclusion-"+acctest.RandString(10), getTestProjectFromEnv())
	}
	return s
}
