package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	resourceManager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestAccResourceManagerLien_basic(t *testing.T) {
	t.Parallel()

	projectName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	org := getTestOrgFromEnv(t)
	var lien resourceManager.Lien

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceManagerLienDestroyProducer(t),
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

		config := googleProviderConfig(t)

		found, err := config.clientResourceManager.Liens.List().Parent(fmt.Sprintf("projects/%s", projectName)).Do()
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
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_resource_manager_lien" {
				continue
			}

			_, err := config.clientResourceManager.Liens.List().Parent(fmt.Sprintf("projects/%s", rs.Primary.Attributes["parent"])).Do()
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
