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
	"github.com/hashicorp/terraform-provider-google/google/services/logging"
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
	org := envvar.GetTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)
	description := "Description " + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderExclusionDestroyProducer(t),
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
	org := envvar.GetTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)
	description := "Description " + acctest.RandString(t, 10)

	checkFn := func(s []*terraform.InstanceState) error {
		loggingExclusionId, err := logging.ParseLoggingExclusionId(s[0].ID)
		if err != nil {
			return err
		}

		folderAttribute := s[0].Attributes["folder"]
		if loggingExclusionId.ResourceId != folderAttribute {
			return fmt.Errorf("imported folder id does not match: actual = %#v expected = %#v", folderAttribute, loggingExclusionId.ResourceId)
		}

		return nil
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderExclusionDestroyProducer(t),
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
	org := envvar.GetTestOrgFromEnv(t)
	exclusionName := "tf-test-exclusion-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)
	parent := "organizations/" + org
	descriptionBefore := "Basic Folder Logging Exclusion" + acctest.RandString(t, 10)
	descriptionAfter := "Updated Basic Folder Logging Exclusion" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderExclusionDestroyProducer(t),
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
	org := envvar.GetTestOrgFromEnv(t)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)
	parent := "organizations/" + org

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderExclusionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderExclusion_multipleCfg(folderName, parent, "tf-test-exclusion-"+acctest.RandString(t, 10)),
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

func testAccCheckLoggingFolderExclusionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_folder_exclusion" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).Folders.Exclusions.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("folder exclusion still exists")
			}
		}

		return nil
	}
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
`, exclusionName, description, envvar.GetTestProjectFromEnv(), folderName, folderParent)
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
`, exclusionName, description, envvar.GetTestProjectFromEnv(), folderName, folderParent)
}

func testAccLoggingFolderExclusion_multipleCfg(folderName, folderParent, exclusionName string) string {
	s := fmt.Sprintf(`
resource "google_folder" "my-folder" {
	display_name = "%s"
	parent       = "%s"
}
`, folderName, folderParent)

	for i := 0; i < 3; i++ {
		s += fmt.Sprintf(`
resource "google_logging_folder_exclusion" "basic%d" {
  name        = "%s%d"
  folder      = element(split("/", google_folder.my-folder.name), 1)
  description = "Basic Folder Logging Exclusion"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}
`, i, exclusionName, i, envvar.GetTestProjectFromEnv())
	}
	return s
}
