// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeSharedVpc_basic(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	hostProject := fmt.Sprintf("tf-test-h-%d", acctest.RandInt(t))
	serviceProject := fmt.Sprintf("tf-test-s-%d", acctest.RandInt(t))

	hostProjectResourceName := "google_compute_shared_vpc_host_project.host"
	serviceProjectResourceName := "google_compute_shared_vpc_service_project.service"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			//Create resources with the deletion_policy flag
			{
				Config: testAccComputeSharedVpc_SharedVPCServiceProjectWithDeletionPolicy(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject(t, hostProject, true),
					testAccCheckComputeSharedVpcServiceProject(t, hostProject, serviceProject, true),
				),
			},
			// Test that resource haven't been deleted after resource has been removed
			{
				Config: testAccComputeSharedVpc_SharedVPCServiceProjectWithDeletionPolicyDeleted(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject(t, hostProject, true),
					testAccCheckComputeSharedVpcServiceProject(t, hostProject, serviceProject, true),
				),
			},
			{
				Config: testAccComputeSharedVpc_basic(hostProject, serviceProject, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSharedVpcHostProject(t, hostProject, true),
					testAccCheckComputeSharedVpcServiceProject(t, hostProject, serviceProject, true),
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
					testAccCheckComputeSharedVpcHostProject(t, hostProject, false),
					testAccCheckComputeSharedVpcServiceProject(t, hostProject, serviceProject, false),
				),
			},
		},
	})
}

func testAccCheckComputeSharedVpcHostProject(t *testing.T, hostProject string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Projects.Get(hostProject).Do()
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

func testAccCheckComputeSharedVpcServiceProject(t *testing.T, hostProject, serviceProject string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		serviceHostProject, err := config.NewComputeClient(config.UserAgent).Projects.GetXpnHost(serviceProject).Do()
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

func testAccComputeSharedVpc_SharedVPCServiceProjectWithDeletionPolicy(hostProject, serviceProject, org, billing string) string {
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
  deletion_policy = "ABANDON"
  depends_on = [
    google_compute_shared_vpc_host_project.host,
    google_project_service.service,
  ]
}
`, hostProject, hostProject, org, billing, serviceProject, serviceProject, org, billing)
}

func testAccComputeSharedVpc_SharedVPCServiceProjectWithDeletionPolicyDeleted(hostProject, serviceProject, org, billing string) string {
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

`, hostProject, hostProject, org, billing, serviceProject, serviceProject, org, billing)
}
