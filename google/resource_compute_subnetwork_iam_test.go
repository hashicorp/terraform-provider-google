package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeSubnetworkIamBinding(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.networkUser"
	region := getTestRegionFromEnv()
	subnetwork := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetworkIamBinding_basic(account, region, subnetwork, role),
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", region, subnetwork, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccComputeSubnetworkIamBinding_update(account, region, subnetwork, role),
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", region, subnetwork, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSubnetworkIamMember(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.networkUser"
	region := getTestRegionFromEnv()
	subnetwork := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccComputeSubnetworkIamMember_basic(account, region, subnetwork, role),
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", region, subnetwork, role, account, project),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSubnetworkIamPolicy(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.networkUser"
	region := getTestRegionFromEnv()
	subnetwork := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetworkIamPolicy_basic(account, region, subnetwork, role),
			},
			// Test a few import formats
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", project, region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", project, region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", region, subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s", subnetwork),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSubnetworkIamBinding_basic(account, region, subnetworkName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_compute_network" "network" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name = "%s"
  region = "%s"
  ip_cidr_range = "10.1.0.0/16"
  network = "${google_compute_network.network.name}"
}

resource "google_compute_subnetwork_iam_binding" "foo" {
  project     = "${google_compute_subnetwork.subnetwork.project}"
  region      = "${google_compute_subnetwork.subnetwork.region}"
  subnetwork  = "${google_compute_subnetwork.subnetwork.name}"
  role        = "%s"
  members     = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account, subnetworkName, subnetworkName, region, roleId)
}

func testAccComputeSubnetworkIamBinding_update(account, region, subnetworkName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_service_account" "test_account_2" {
  account_id   = "%s-2"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_compute_network" "network" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name = "%s"
  region = "%s"
  ip_cidr_range = "10.1.0.0/16"
  network = "${google_compute_network.network.name}"
}

resource "google_compute_subnetwork_iam_binding" "foo" {
  project     = "${google_compute_subnetwork.subnetwork.project}"
  region      = "${google_compute_subnetwork.subnetwork.region}"
  subnetwork  = "${google_compute_subnetwork.subnetwork.name}"
  role         = "%s"
  members      = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}"
  ]
}
`, account, account, subnetworkName, subnetworkName, region, roleId)
}

func testAccComputeSubnetworkIamMember_basic(account, region, subnetworkName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_compute_network" "network" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name = "%s"
  region = "%s"
  ip_cidr_range = "10.1.0.0/16"
  network = "${google_compute_network.network.name}"
}

resource "google_compute_subnetwork_iam_member" "foo" {
  project     = "${google_compute_subnetwork.subnetwork.project}"
  region      = "${google_compute_subnetwork.subnetwork.region}"
  subnetwork  = "${google_compute_subnetwork.subnetwork.name}"
  role        = "%s"
  member      = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, subnetworkName, subnetworkName, region, roleId)
}

func testAccComputeSubnetworkIamPolicy_basic(account, region, subnetworkName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Subnetwork Iam Testing Account"
}

resource "google_compute_network" "network" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name = "%s"
  region = "%s"
  ip_cidr_range = "10.1.0.0/16"
  network = "${google_compute_network.network.name}"
}

data "google_iam_policy" "foo" {
	binding {
		role = "%s"

		members = ["serviceAccount:${google_service_account.test_account.email}"]
	}
}

resource "google_compute_subnetwork_iam_policy" "foo" {
  project     = "${google_compute_subnetwork.subnetwork.project}"
  region      = "${google_compute_subnetwork.subnetwork.region}"
  subnetwork  = "${google_compute_subnetwork.subnetwork.name}"
  policy_data = "${data.google_iam_policy.foo.policy_data}"
}
`, account, subnetworkName, subnetworkName, region, roleId)
}
