// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	resourceManager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestAccResourceManagerLien_basic(t *testing.T) {
	t.Parallel()

	projectName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	org := envvar.GetTestOrgFromEnv(t)
	var lien resourceManager.Lien

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckResourceManagerLienDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceManagerLien_basic(projectName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceManagerLienExists(
						t, "google_resource_manager_lien.lien", projectName, &lien),
				),
			},
			{
				ResourceName:      "google_resource_manager_lien.lien",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					// This has to be a function to close over lien.Name, which is necessary
					// because Name is a Computed attribute.
					return fmt.Sprintf("%s/%s",
						projectName,
						strings.Split(lien.Name, "/")[1]), nil
				},
			},
		},
	})
}

func testAccCheckResourceManagerLienExists(t *testing.T, n, projectName string, lien *resourceManager.Lien) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewResourceManagerClient(config.UserAgent).Liens.List().Parent(fmt.Sprintf("projects/%s", projectName)).Do()
		if err != nil {
			return err
		}
		if len(found.Liens) != 1 {
			return fmt.Errorf("Lien %s not found", rs.Primary.ID)
		}

		*lien = *found.Liens[0]

		return nil
	}
}

func testAccCheckResourceManagerLienDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_resource_manager_lien" {
				continue
			}

			_, err := config.NewResourceManagerClient(config.UserAgent).Liens.List().Parent(fmt.Sprintf("projects/%s", rs.Primary.Attributes["parent"])).Do()
			if err == nil {
				return fmt.Errorf("Lien %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccResourceManagerLien_basic(projectName, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "%s"
  name       = "some test project"
  org_id     = "%s"
}

resource "google_resource_manager_lien" "lien" {
  parent       = "projects/${google_project.project.project_id}"
  restrictions = ["resourcemanager.projects.delete"]
  origin       = "something"
  reason       = "something else"
}
`, projectName, org)
}
