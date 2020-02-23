package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Logging exclusions don't always work when making parallel requests, so run tests serially
func TestAccLoggingFolderExclusion(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"basic":                       testAccLoggingFolderExclusion_basic,
		"folderAcceptsFullFolderPath": testAccLoggingFolderExclusion_folderAcceptsFullFolderPath,
		"update":                      testAccLoggingFolderExclusion_update,
		"multiple":                    testAccLoggingFolderExclusion_multiple,
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

func testAccLoggingFolderExclusion_basic(t *testing.T) {
	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_basicCfg(exclusionName, description, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingFolderExclusion_folderAcceptsFullFolderPath(t *testing.T) {
	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	checkFn := func(s []*terraform.InstanceState) error {
		loggingExclusionId, err := parseLoggingExclusionId(s[0].ID)
		if err != nil {
			return err
		}

		folderAttribute := s[0].Attributes["folder"]
		if loggingExclusionId.resourceId != folderAttribute {
			return fmt.Errorf("imported folder id does not match: actual = %#v expected = %#v", folderAttribute, loggingExclusionId.resourceId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_withFullFolderPath(exclusionName, description, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.full-folder",
				ImportState:       true,
				ImportStateVerify: true,
				// We support both notations: folder/[FOLDER_ID] and plain [FOLDER_ID] however the
				// importer will always use the plain [FOLDER_ID] notation which will differ from
				// the schema if the schema has used the prefixed notation. We have to check this in
				// a checkFn instead.
				ImportStateVerifyIgnore: []string{"folder"},
				ImportStateCheck:        checkFn,
			},
		},
	})
}

func testAccLoggingFolderExclusion_update(t *testing.T) {
	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	parent := "organizations/" + org
	descriptionBefore := "Basic Folder Logging Exclusion" + acctest.RandString(10)
	descriptionAfter := "Updated Basic Folder Logging Exclusion" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_basicCfg(exclusionName, descriptionBefore, folderName, parent),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingFolderExclusion_basicCfg(exclusionName, descriptionAfter, folderName, parent),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoggingFolderExclusion_multiple(t *testing.T) {
	org := getTestOrgFromEnv(t)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	parent := "organizations/" + org

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_multipleCfg(folderName, parent),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic0",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingFolderExclusionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_folder_exclusion" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.Folders.Exclusions.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("folder exclusion still exists")
		}
	}

	return nil
}

func testAccLoggingFolderExclusion_basicCfg(exclusionName, description, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_exclusion" "basic" {
  name        = "%s"
  folder      = element(split("/", google_folder.my-folder.name), 1)
  description = "%s"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, exclusionName, description, getTestProjectFromEnv(), folderName, folderParent)
}

func testAccLoggingFolderExclusion_withFullFolderPath(exclusionName, description, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_exclusion" "full-folder" {
  name        = "%s"
  folder      = google_folder.my-folder.name
  description = "%s"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, exclusionName, description, getTestProjectFromEnv(), folderName, folderParent)
}

func testAccLoggingFolderExclusion_multipleCfg(folderName, folderParent string) string {
	s := fmt.Sprintf(`
resource "google_folder" "my-folder" {
	display_name = "%s"
	parent       = "%s"
}
`, folderName, folderParent)

	for i := 0; i < 3; i++ {
		s += fmt.Sprintf(`
resource "google_logging_folder_exclusion" "basic%d" {
  name        = "%s"
  folder      = element(split("/", google_folder.my-folder.name), 1)
  description = "Basic Folder Logging Exclusion"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, i, "tf-test-exclusion-"+acctest.RandString(10), getTestProjectFromEnv())
	}
	return s
}
