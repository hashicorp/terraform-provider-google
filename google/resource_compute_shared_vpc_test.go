package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

func TestAccComputeSharedVpc_basic(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG", "GOOGLE_BILLING_ACCOUNT")
	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")

	hostProject := "xpn-host-" + acctest.RandString(10)
	serviceProject := "xpn-service-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpc_basic(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject("google_compute_shared_vpc_host_project.host"),
					testAccCheckComputeSharedVpcServiceProject("google_compute_shared_vpc_service_project.service"),
				),
			},
		},
	})
}

func testAccCheckComputeSharedVpcHostProject(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Projects.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Error reading project %s: %s", rs.Primary.ID, err)
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Project %s not found", rs.Primary.ID)
		}

		if found.XpnProjectStatus != "HOST" {
			return fmt.Errorf("Project %q Shared VPC status was not expected, got %q", rs.Primary.ID, found.XpnProjectStatus)
		}

		return nil
	}
}

func testAccCheckComputeSharedVpcServiceProject(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		hostProject := rs.Primary.Attributes["host_project"]
		serviceProject := rs.Primary.Attributes["service_project"]

		serviceHostProject, err := config.clientCompute.Projects.GetXpnHost(serviceProject).Do()
		if err != nil {
			return err
		}

		if serviceHostProject.Name != hostProject {
			return fmt.Errorf("Wrong host project for the given service project. Expected '%s', got '%s'", hostProject, serviceHostProject.Id)
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

resource "google_project_services" "host" {
	project  = "${google_project.host.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_project_services" "service" {
	project  = "${google_project.service.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_compute_shared_vpc_host_project" "host" {
	project    = "${google_project.host.project_id}"
	depends_on = ["google_project_services.host"]
}

resource "google_compute_shared_vpc_service_project" "service" {
	host_project    = "${google_project.host.project_id}"
	service_project = "${google_project.service.project_id}"
	depends_on      = ["google_compute_shared_vpc_host_project.host", "google_project_services.service"]
}`, hostProject, hostProject, org, billing, serviceProject, serviceProject, org, billing)
}
