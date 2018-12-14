package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test that services can be enabled and disabled on a project
func TestAccProjectService_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	services := []string{"iam.googleapis.com", "cloudresourcemanager.googleapis.com"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_basic(services, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(services, pid, true),
				),
			},
			{
				ResourceName:            "google_project_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
			{
				ResourceName:            "google_project_service.test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(services, pid, false),
				),
			},
			// Create services with disabling turned off.
			{
				Config: testAccProjectService_noDisable(services, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(services, pid, true),
				),
			},
			// Check that services are still enabled even after the resources are deleted.
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(services, pid, true),
				),
			},
		},
	})
}

func TestAccProjectService_handleNotFound(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "iam.googleapis.com"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_handleNotFound(service, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService([]string{service}, pid, true),
				),
			},
			// Delete the project, implicitly deletes service, expect the plan to want to create the service again
			{
				Config:             testAccProjectService_handleNotFoundNoProject(service, pid),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckProjectService(services []string, pid string, expectEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		apiServices, err := getApiServices(pid, config, map[string]struct{}{})
		if err != nil {
			return fmt.Errorf("Error listing services for project %q: %v", pid, err)
		}

		for _, expected := range services {
			exists := false
			for _, actual := range apiServices {
				if expected == actual {
					exists = true
				}
			}
			if expectEnabled && !exists {
				return fmt.Errorf("Expected service %s is not enabled server-side (found %v)", expected, apiServices)
			}
			if !expectEnabled && exists {
				return fmt.Errorf("Expected disabled service %s is enabled server-side", expected)
			}
		}

		return nil
	}
}

func testAccProjectService_basic(services []string, pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "test" {
  project = "${google_project.acceptance.project_id}"
  service = "%s"
}

resource "google_project_service" "test2" {
  project = "${google_project.acceptance.project_id}"
  service = "%s"
}
`, pid, name, org, services[0], services[1])
}

func testAccProjectService_noDisable(services []string, pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "test" {
  project = "${google_project.acceptance.project_id}"
  service = "%s"
  disable_on_destroy = false
}

resource "google_project_service" "test2" {
  project = "${google_project.acceptance.project_id}"
  service = "%s"
  disable_on_destroy = false
}
`, pid, name, org, services[0], services[1])
}

func testAccProjectService_handleNotFound(service, pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

// by passing through locals, we break the dependency chain
// see terraform-provider-google#1292
locals {
  project_id = "${google_project.acceptance.project_id}"
}

resource "google_project_service" "test" {
  project = "${local.project_id}"
  service = "%s"
}
`, pid, name, org, service)
}

func testAccProjectService_handleNotFoundNoProject(service, pid string) string {
	return fmt.Sprintf(`
resource "google_project_service" "test" {
  project = "%s"
  service = "%s"
}
`, pid, service)
}
