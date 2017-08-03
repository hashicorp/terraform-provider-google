package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSourceRepoRepository_basic(t *testing.T) {
	repositoryName := fmt.Sprintf("source-repos-repository-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceRepoRepositoryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSourceRepoRepository_basic(repositoryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceRepoRepositoryExists(
						"google_sourcerepos_repository.acceptance", repositoryName),
				),
			},
		},
	})
}

func testAccCheckSourceRepoRepositoryDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "google_sourcerepos_repository" {
			repositoryName := buildRepositoryName(config.Project, rs.Primary.Attributes["name"])

			_, err := config.clientSourceRepo.Projects.Repos.Get(repositoryName).Do()
			if err == nil {
				return fmt.Errorf(repositoryName + "Source Repos Repository still exists")
			}
		}
	}

	return nil
}

func testAccCheckSourceRepoRepositoryExists(resourceType, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceType]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		repositoryName := buildRepositoryName(config.Project, resourceName)

		resp, err := config.clientSourceRepo.Projects.Repos.Get(repositoryName).Do()

		if err != nil {
			return fmt.Errorf("Error confirming Source Repos Repository existence: %#v", err)
		}

		if resp.Name != repositoryName {
			return fmt.Errorf("Failed to verify Source Repos Repository by Name")
		}
		return nil
	}
}

func testAccSourceRepoRepository_basic(repositoryName string) string {
	return fmt.Sprintf(`
	resource "google_sourcerepos_repository" "acceptance" {
	  name = "%s"
	}
	`, repositoryName)
}
