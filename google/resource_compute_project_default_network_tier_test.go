package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeProjectDefaultNetworkTier_basic(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)
	projectID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_defaultNetworkTier_premium(projectID, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectDefaultNetworkTier_modify(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)
	projectID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_defaultNetworkTier_premium(projectID, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeProject_defaultNetworkTier_standard(projectID, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeProject_defaultNetworkTier_premium(projectID, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_default_network_tier" "fizzbuzz" {
  project      = google_project.project.project_id
  network_tier = "PREMIUM"
  depends_on   = [google_project_service.compute]
}
`, projectID, projectID, org, billing)
}

func testAccComputeProject_defaultNetworkTier_standard(projectID, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_default_network_tier" "fizzbuzz" {
  project      = google_project.project.project_id
  network_tier = "STANDARD"
  depends_on   = [google_project_service.compute]
}
`, projectID, projectID, org, billing)
}
