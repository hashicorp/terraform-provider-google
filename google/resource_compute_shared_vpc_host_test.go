package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeSharedVpcHost_basic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	pid := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpcHost_basic(pid, pname, org, billingId),
				Check:  testAccCheckComputeSharedVpcHostEnabled("google_compute_shared_vpc_host.host", true),
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist
			// in order to check the XPN status.
			resource.TestStep{
				Config: testAccComputeSharedVpcHost_disabled(pid, pname, org, billingId),
				// Use the project ID since the google_compute_shared_vpc_host resource no longer exists
				Check: testAccCheckComputeSharedVpcHostEnabled("google_project.project", false),
			},
		},
	})
}

func testAccCheckComputeSharedVpcHostEnabled(n string, enabled bool) resource.TestCheckFunc {
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

		if enabled != (found.XpnProjectStatus == "HOST") {
			return fmt.Errorf("Project %s XPN status was not expected, got %q", rs.Primary.ID, found.XpnProjectStatus)
		}

		return nil
	}
}

func testAccComputeSharedVpcHost_basic(pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id      = "%s"
	name            = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "project" {
	project = "${google_project.project.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_compute_shared_vpc_host" "host" {
	project = "${google_project.project.project_id}"

	depends_on = ["google_project_services.project"]
}`, pid, name, org, billing)
}

func testAccComputeSharedVpcHost_disabled(pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id      = "%s"
	name            = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "project" {
	project = "${google_project.project.project_id}"
	services = ["compute.googleapis.com"]
}`, pid, name, org, billing)
}
