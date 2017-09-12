package google

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeSharedVpc_basic(t *testing.T) {
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
				Config: testAccComputeSharedVpc_basic(pid, pid2, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHost("google_compute_shared_vpc.vpc", true),
					testAccCheckComputeSharedVpcResources("google_compute_shared_vpc.vpc", []string{pid2})),
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist
			// in order to check the XPN status.
			resource.TestStep{
				Config: testAccComputeSharedVpc_disabled(pid, pid2, pname, org, billingId),
				// Use the project ID since the google_compute_shared_vpc_host resource no longer exists
				Check: testAccCheckComputeSharedVpcHost("google_project.host", false),
			},
		},
	})
}

func TestAccComputeSharedVpc_update(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	pid := "terraform-" + acctest.RandString(10)
	pid2 := "terraform-" + acctest.RandString(10)
	pid3 := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpc_basic(pid, pid2, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHost("google_compute_shared_vpc.vpc", true),
					testAccCheckComputeSharedVpcResources("google_compute_shared_vpc.vpc", []string{pid2})),
			},
			resource.TestStep{
				Config: testAccComputeSharedVpc_addServiceProjects(pid, pid2, pid3, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHost("google_compute_shared_vpc.vpc", true),
					testAccCheckComputeSharedVpcResources("google_compute_shared_vpc.vpc", []string{pid2, pid3})),
			},
			resource.TestStep{
				Config: testAccComputeSharedVpc_removeServiceProjects(pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHost("google_compute_shared_vpc.vpc", true),
					testAccCheckComputeSharedVpcResources("google_compute_shared_vpc.vpc", []string{})),
			},
		},
	})
}

func testAccCheckComputeSharedVpcHost(n string, enabled bool) resource.TestCheckFunc {
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
			return fmt.Errorf("Project %q Shared VPC status was not expected, got %q", rs.Primary.ID, found.XpnProjectStatus)
		}

		return nil
	}
}

func testAccCheckComputeSharedVpcResources(n string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		tfServiceProjects := []string{}
		// We don't know the exact keys of the elements, so go through the whole list looking for matching ones
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, "service_projects") && k != "service_projects.#" {
				tfServiceProjects = append(tfServiceProjects, v)
			}
		}

		sort.Strings(tfServiceProjects)
		sort.Strings(expected)

		if !reflect.DeepEqual(expected, tfServiceProjects) {
			return fmt.Errorf("Service projects mismatch. Expected: %v, Actual: %v", expected, tfServiceProjects)
		}
		return nil
	}
}

func testAccComputeSharedVpc_basic(pid, pid2, name, org, billing string) string {
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
	service_projects = ["${google_project.service.project_id}"]

	depends_on = ["google_project_services.host", "google_project_services.service"]
}`, pid, name, org, billing, pid2, name, org, billing)
}

func testAccComputeSharedVpc_disabled(pid, pid2, name, org, billing string) string {
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
}`, pid, name, org, billing, pid2, name, org, billing)
}

func testAccComputeSharedVpc_addServiceProjects(pid, pid2, pid3, name, org, billing string) string {
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

resource "google_project" "service2" {
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

resource "google_project_services" "service2" {
	project = "${google_project.service2.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_compute_shared_vpc" "vpc" {
	host_project     = "${google_project.host.project_id}"
	service_projects = ["${google_project.service.project_id}", "${google_project.service2.project_id}"]

	depends_on = ["google_project_services.host", "google_project_services.service", "google_project_services.service2"]
}`, pid, name, org, billing,
		pid2, name, org, billing,
		pid3, name, org, billing)
}

func testAccComputeSharedVpc_removeServiceProjects(pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
	project_id      = "%s"
	name            = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "host" {
	project = "${google_project.host.project_id}"
	services = ["compute.googleapis.com"]
}

resource "google_compute_shared_vpc" "vpc" {
	host_project     = "${google_project.host.project_id}"

	depends_on = ["google_project_services.host"]
}`, pid, name, org, billing)
}
