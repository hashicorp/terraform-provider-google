package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingFolderExclusion_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	var exclusion logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_basic(exclusionName, description, folderName, "organizations/"+org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderExclusionExists("google_logging_folder_exclusion.basic", &exclusion),
					testAccCheckLoggingFolderExclusion(&exclusion, "google_logging_folder_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderExclusion_folderAcceptsFullFolderPath(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	description := "Description " + acctest.RandString(10)

	var exclusion logging.LogExclusion

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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderExclusionExists("google_logging_folder_exclusion.full-folder", &exclusion),
					testAccCheckLoggingFolderExclusion(&exclusion, "google_logging_folder_exclusion.full-folder"),
				),
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

func TestAccLoggingFolderExclusion_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)
	folderName := "tf-test-folder-" + acctest.RandString(10)
	parent := "organizations/" + org
	descriptionBefore := "Basic Folder Logging Exclusion" + acctest.RandString(10)
	descriptionAfter := "Updated Basic Folder Logging Exclusion" + acctest.RandString(10)

	var exclusionBefore, exclusionAfter logging.LogExclusion

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingFolderExclusionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_basic(exclusionName, descriptionBefore, folderName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderExclusionExists("google_logging_folder_exclusion.basic", &exclusionBefore),
					testAccCheckLoggingFolderExclusion(&exclusionBefore, "google_logging_folder_exclusion.basic"),
				),
			},
			{
				Config: testAccLoggingFolderExclusion_basic(exclusionName, descriptionAfter, folderName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderExclusionExists("google_logging_folder_exclusion.basic", &exclusionAfter),
					testAccCheckLoggingFolderExclusion(&exclusionAfter, "google_logging_folder_exclusion.basic"),
				),
			},
			{
				ResourceName:      "google_logging_folder_exclusion.basic",
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

func testAccCheckLoggingFolderExclusionExists(n string, exclusion *logging.LogExclusion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.Folders.Exclusions.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*exclusion = *si

		return nil
	}
}

func testAccCheckLoggingFolderExclusion(exclusion *logging.LogExclusion, n string) resource.TestCheckFunc {
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

func testAccLoggingFolderExclusion_basic(exclusionName, description, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_exclusion" "basic" {
	name             = "%s"
	folder           = "${element(split("/", google_folder.my-folder.name), 1)}"
	description      = "%s"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_folder" "my-folder" {
	display_name = "%s"
	parent       = "%s"
}`, exclusionName, description, getTestProjectFromEnv(), folderName, folderParent)
}

func testAccLoggingFolderExclusion_withFullFolderPath(exclusionName, description, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_exclusion" "full-folder" {
	name             = "%s"
	folder           = "${google_folder.my-folder.name}"
	description      = "%s"
	filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_folder" "my-folder" {
	display_name = "%s"
	parent       = "%s"
}`, exclusionName, description, getTestProjectFromEnv(), folderName, folderParent)
}
