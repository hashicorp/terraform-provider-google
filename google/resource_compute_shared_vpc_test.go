package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeSharedVpc_basic(t *testing.T) {
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)

	hostProject := "xpn-host-" + acctest.RandString(10)
	serviceProject := "xpn-service-" + acctest.RandString(10)

	hostProjectResourceName := "google_compute_shared_vpc_host_project.host"
	serviceProjectResourceName := "google_compute_shared_vpc_service_project.service"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSharedVpc_basic(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject(hostProject, true),
					testAccCheckComputeSharedVpcServiceProject(hostProject, serviceProject, true),
				),
			},
			// Test import.
			{
				ResourceName:      hostProjectResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      serviceProjectResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testAccComputeSharedVpc_disabled(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject(hostProject, false),
					testAccCheckComputeSharedVpcServiceProject(hostProject, serviceProject, false),
				),
			},
		},
	})
}

func testAccCheckComputeSharedVpcHostProject(hostProject string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Projects.Get(hostProject).Do()
		if err != nil {
			return fmt.Errorf("Error reading project %s: %s", hostProject, err)
		}

		if found.Name != hostProject {
			return fmt.Errorf("Project %s not found", hostProject)
		}

		if enabled != (found.XpnProjectStatus == "HOST") {
			return fmt.Errorf("Project %q shared VPC status was not expected, got %q", hostProject, found.XpnProjectStatus)
		}

		return nil
	}
}

func testAccCheckComputeSharedVpcServiceProject(hostProject, serviceProject string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		serviceHostProject, err := config.clientCompute.Projects.GetXpnHost(serviceProject).Do()
		if err != nil {
			if enabled {
				return fmt.Errorf("Expected service project to be enabled.")
			}
			return nil
		}

		if enabled != (serviceHostProject.Name == hostProject) {
			return fmt.Errorf("Wrong host project for the given service project. Expected '%s', got '%s'", hostProject, serviceHostProject.Name)
		}

		return nil
	}
}

func testAccComputeSharedVpc_basic(hostProject, serviceProject, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project" "service" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "host" {
  project = google_project.host.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "service" {
  project = google_project.service.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_shared_vpc_host_project" "host" {
  project    = google_project.host.project_id
  depends_on = [google_project_service.host]
}

resource "google_compute_shared_vpc_service_project" "service" {
  host_project    = google_project.host.project_id
  service_project = google_project.service.project_id
  depends_on = [
    google_compute_shared_vpc_host_project.host,
    google_project_service.service,
  ]
}
`, hostProject, hostProject, org, billing, serviceProject, serviceProject, org, billing)
}

func testAccComputeSharedVpc_disabled(hostProject, serviceProject, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project" "service" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "host" {
  project = google_project.host.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "service" {
  project = google_project.service.project_id
  service = "compute.googleapis.com"
}
`, hostProject, hostProject, org, billing, serviceProject, serviceProject, org, billing)
}
