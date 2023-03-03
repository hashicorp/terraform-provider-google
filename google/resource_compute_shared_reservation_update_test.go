package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeSharedReservation_update(t *testing.T) {
	SkipIfVcr(t) // large number of parallel resources.
	t.Parallel()

	context := map[string]interface{}{
		"project":         GetTestProjectFromEnv(),
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservation_sharedReservation_basic(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "share_settings"},
			},
			{
				Config: testAccComputeReservation_sharedReservation_update(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "share_settings"},
			},
			{
				Config: testAccComputeReservation_sharedReservation_basic(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "share_settings"},
			},
		},
	})
}

func testAccComputeReservation_sharedReservation_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "owner_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}


resource "google_project_service" "compute" {
  project = google_project.owner_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project" "guest_project" {
  project_id      = "tf-test-2%{random_suffix}"
  name            = "tf-test-2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project" "guest_project_second" {
  project_id      = "tf-test-3%{random_suffix}"
  name            = "tf-test-3%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project" "guest_project_third" {
  project_id      = "tf-test-4%{random_suffix}"
  name            = "tf-test-4%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_organization_policy" "shared_reservation_org_policy" {
  org_id     = "%{org_id}"
  constraint = "constraints/compute.sharedReservationsOwnerProjects"
  list_policy {
    allow {
      values = ["projects/${google_project.owner_project.number}"]
    }
  }
}

resource "google_project_service" "compute_second_project" {
  project = google_project.guest_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute_third_project" {
  project = google_project.guest_project_second.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute_fourth_project" {
  project = google_project.guest_project_third.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_compute_reservation" "gce_reservation" {
  project = google_project.owner_project.project_id
  name = "my-reservation"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
  share_settings {
    share_type = "SPECIFIC_PROJECTS"
    project_map {
      id = google_project.guest_project.project_id
      project_id = google_project.guest_project.project_id
    }
  }
  depends_on = [google_organization_policy.shared_reservation_org_policy,google_project_service.compute,google_project_service.compute_second_project,google_project_service.compute_third_project]
}
`, context)
}

func testAccComputeReservation_sharedReservation_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "owner_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "compute" {
  project = google_project.owner_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project" "guest_project" {
  project_id      = "tf-test-2%{random_suffix}"
  name            = "tf-test-2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project" "guest_project_second" {
  project_id      = "tf-test-3%{random_suffix}"
  name            = "tf-test-3%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project" "guest_project_third" {
  project_id      = "tf-test-4%{random_suffix}"
  name            = "tf-test-4%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_organization_policy" "shared_reservation_org_policy" {
  org_id     = "%{org_id}"
  constraint = "constraints/compute.sharedReservationsOwnerProjects"
  list_policy {
    allow {
      values = ["projects/${google_project.owner_project.number}"]
    }
  }
}

resource "google_project_service" "compute_second_project" {
  project = google_project.guest_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute_third_project" {
  project = google_project.guest_project_second.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute_fourth_project" {
  project = google_project.guest_project_third.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_compute_reservation" "gce_reservation" {
  project = google_project.owner_project.project_id
  name = "my-reservation"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
  share_settings {
    share_type = "SPECIFIC_PROJECTS"
    project_map {
      id = google_project.guest_project.project_id
      project_id = google_project.guest_project.project_id
    }
    project_map {
      id = google_project.guest_project_second.project_id
      project_id = google_project.guest_project_second.project_id
    }
    project_map {
      id = google_project.guest_project_third.project_id
      project_id = google_project.guest_project_third.project_id
    }
  }
  depends_on = [google_organization_policy.shared_reservation_org_policy,google_project_service.compute,google_project_service.compute_second_project,google_project_service.compute_third_project]
}
`, context)
}
