package google

import (
	"fmt"
	"os"
	"testing"

	compute "google.golang.org/api/compute/v1"

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

	var project compute.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpcHost_basic(pid, pname, org, billingId),
				Check:  testAccCheckComputeSharedVpcHostEnabled("google_compute_shared_vpc_host.host", &project),
			},
		},
	})

	// There doesn't seem to be a way to test disabling. Since disabling removes the resource, the only way
	// to read it is to know the project id already. We can get the project id from testAccCheckComputeSharedVpcHostEnabled,
	// but we can't use it in another TestCheckFunc because it won't have been initialized (there's a real reason for this
	// but I don't understand it enough to explain here). If we try to do the check here, the project will already have
	// been deleted.
}

func testAccCheckComputeSharedVpcHostEnabled(n string, project *compute.Project) resource.TestCheckFunc {
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
			return fmt.Errorf("Project %s XPN status was not host, got %q", rs.Primary.ID, found.XpnProjectStatus)
		}

		*project = *found

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
