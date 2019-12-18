package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Add two key value pairs
func TestAccComputeProjectMetadata_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectID := "terraform-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeProjectMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_basic0_metadata(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_metadata.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Add three key value pairs, then replace one and modify a second
func TestAccComputeProjectMetadata_modify_1(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectID := "terraform-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeProjectMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_modify0_metadata(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_metadata.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeProject_modify1_metadata(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_metadata.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Add two key value pairs, and replace both
func TestAccComputeProjectMetadata_modify_2(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectID := "terraform-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeProjectMetadataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_basic0_metadata(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_metadata.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeProject_basic1_metadata(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_metadata.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeProjectMetadataDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_project_metadata" {
			continue
		}

		project, err := config.clientCompute.Projects.Get(rs.Primary.ID).Do()
		if err == nil && len(project.CommonInstanceMetadata.Items) > 0 {
			return fmt.Errorf("Error, metadata items still exist in %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccComputeProject_basic0_metadata(projectID, name, org, billing string) string {
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

resource "google_compute_project_metadata" "fizzbuzz" {
  project = google_project.project.project_id
  metadata = {
    banana = "orange"
    sofa   = "darwinism"
  }
  depends_on = [google_project_service.compute]
}
`, projectID, name, org, billing)
}

func testAccComputeProject_basic1_metadata(projectID, name, org, billing string) string {
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

resource "google_compute_project_metadata" "fizzbuzz" {
  project = google_project.project.project_id
  metadata = {
    kiwi    = "papaya"
    finches = "darwinism"
  }
  depends_on = [google_project_service.compute]
}
`, projectID, name, org, billing)
}

func testAccComputeProject_modify0_metadata(projectID, name, org, billing string) string {
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

resource "google_compute_project_metadata" "fizzbuzz" {
  project = google_project.project.project_id
  metadata = {
    paper        = "pen"
    genghis_khan = "french bread"
    happy        = "smiling"
  }
  depends_on = [google_project_service.compute]
}
`, projectID, name, org, billing)
}

func testAccComputeProject_modify1_metadata(projectID, name, org, billing string) string {
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

resource "google_compute_project_metadata" "fizzbuzz" {
  project = google_project.project.project_id
  metadata = {
    paper = "pen"
    paris = "french bread"
    happy = "laughing"
  }
  depends_on = [google_project_service.compute]
}
`, projectID, name, org, billing)
}
