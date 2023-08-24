// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGoogleIamTestablePermissions_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name = "//cloudresourcemanager.googleapis.com/projects/%s"
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						[]string{"GA"},
						"",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name = "//cloudresourcemanager.googleapis.com/projects/%s"
				stages = ["GA"]
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						[]string{"GA"},
						"",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name   = "//cloudresourcemanager.googleapis.com/projects/%s"
				custom_support_level = "NOT_SUPPORTED"
				stages               = ["BETA"]
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						[]string{"BETA"},
						"NOT_SUPPORTED",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name   = "//cloudresourcemanager.googleapis.com/projects/%s"
				custom_support_level = "not_supported"
				stages               = ["beta"]
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						[]string{"BETA"},
						"NOT_SUPPORTED",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name   = "//cloudresourcemanager.googleapis.com/projects/%s"
				stages               = ["ga", "beta"]
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						[]string{"GA", "BETA"},
						"",
					),
				),
			},
		},
	})
}

func testAccCheckGoogleIamTestablePermissionsMeta(project string, n string, expectedStages []string, expectedSupportLevel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find perms data source: %s", n)
		}
		expectedId := fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", project)
		if rs.Primary.ID != expectedId {
			return fmt.Errorf("perms data source ID not set.")
		}
		attrs := rs.Primary.Attributes
		count, ok := attrs["permissions.#"]
		if !ok {
			return fmt.Errorf("can't find 'permsissions' attribute")
		}
		permCount, err := strconv.Atoi(count)
		if err != nil {
			return err
		}
		if permCount < 2 {
			return fmt.Errorf("count should be greater than 2")
		}
		foundStageCounter := len(expectedStages)
		foundSupport := false

		for i := 0; i < permCount; i++ {
			for s := 0; s < len(expectedStages); s++ {
				stageKey := "permissions." + strconv.Itoa(i) + ".stage"
				supportKey := "permissions." + strconv.Itoa(i) + ".custom_support_level"
				if tpgresource.StringInSlice(expectedStages, attrs[stageKey]) {
					foundStageCounter -= 1
				}
				if attrs[supportKey] == expectedSupportLevel {
					foundSupport = true
				}
				if foundSupport && foundStageCounter == 0 {
					return nil
				}
			}
		}

		if foundSupport { // This means we didn't find a stage
			return fmt.Errorf("Could not find stages %v in output", expectedStages)
		}
		if foundStageCounter == 0 { // This meads we didn't fins a custom_support_level
			return fmt.Errorf("Could not find custom_support_level %s in output", expectedSupportLevel)
		}
		return fmt.Errorf("Unable to find customSupportLevel or stages in output")
	}
}
