package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeSharedVpcServiceProjectAssociation_basic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	pid := "terraform-" + acctest.RandString(10)
	pid2 := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpcServiceProjectAssociation_basic(pid, pid2, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHost("google_compute_shared_vpc.vpc", true),
				),
			},
		},
	})
}

func testAccComputeSharedVpcServiceProjectAssociation_basic(pid, pid2, name, org, billing string) string {
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
	project = "${google_project.host.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_project_services" "service" {
	project = "${google_project.service.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_compute_shared_vpc" "vpc" {
	host_project     = "${google_project.host.project_id}"

	depends_on = ["google_project_services.host"]
}

resource "google_compute_shared_vpc_service_project_association" "service" {
	host_project    = "${google_compute_shared_vpc.vpc.host_project}"
	service_project = "${google_project.service.project_id}"

	depends_on = ["google_project_services.service"]
  }
`, pid, name, org, billing, pid2, name, org, billing)
}
