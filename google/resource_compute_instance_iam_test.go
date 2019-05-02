package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeInstanceIamBinding(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.osLogin"
	region := getTestRegionFromEnv()
	zone := getTestZoneFromEnv()
	subnetwork := fmt.Sprintf("tf-test-net-%s", acctest.RandString(10))
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceIamBinding_basic(project, account, region, zone, subnetwork, instanceName, role),
			},
			{
				ResourceName:      "google_compute_instance_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s %s", project, zone, instanceName, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccComputeInstanceIamBinding_update(project, account, region, zone, subnetwork, instanceName, role),
			},
			{
				ResourceName:      "google_compute_instance_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s %s", project, zone, instanceName, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceIamMember(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.osLogin"
	region := getTestRegionFromEnv()
	zone := getTestZoneFromEnv()
	subnetwork := fmt.Sprintf("tf-test-net-%s", acctest.RandString(10))
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccComputeInstanceIamMember_basic(project, account, region, zone, subnetwork, instanceName, role),
			},
			{
				ResourceName:      "google_compute_instance_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", project, zone, instanceName, role, account, project),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceIamPolicy(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/compute.osLogin"
	region := getTestRegionFromEnv()
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))
	subnetwork := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceIamPolicy_basic(project, account, region, zone, subnetwork, instanceName, role),
			},
			// Test a few import formats
			{
				ResourceName:      "google_compute_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeInstanceIamMember_basic(project, account, region, zone, subnetworkName, instanceName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
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

resource "google_compute_instance" "test_vm" {
	project = "%s"
	zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    subnetwork ="${google_compute_subnetwork.subnetwork.self_link}"
  }
}

resource "google_compute_instance_iam_member" "foo" {
	project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance_name = "${google_compute_instance.test_vm.name}"
  role        = "%s"
  member      = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, subnetworkName, subnetworkName, region, project, zone, instanceName, roleId)
}

func testAccComputeInstanceIamPolicy_basic(project, account, region, zone, subnetworkName, instanceName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
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

resource "google_compute_instance" "test_vm" {
	project = "%s"
	zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    subnetwork ="${google_compute_subnetwork.subnetwork.self_link}"
  }
}

data "google_iam_policy" "foo" {
	binding {
		role = "%s"
		members = ["serviceAccount:${google_service_account.test_account.email}"]
	}
}

resource "google_compute_instance_iam_policy" "foo" {
	project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance_name = "${google_compute_instance.test_vm.name}"
	policy_data = "${data.google_iam_policy.foo.policy_data}"
}
`, account, subnetworkName, subnetworkName, region, project, zone, instanceName, roleId)
}

func testAccComputeInstanceIamBinding_basic(project, account, region, zone, subnetworkName, instanceName, roleId string) string {
	return fmt.Sprintf(`
	resource "google_service_account" "test_account" {
	  account_id   = "%s"
	  display_name = "Iam Testing Account"
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

	resource "google_compute_instance" "test_vm" {
		project = "%s"
		zone         = "%s"
	  name         = "%s"
	  machine_type = "n1-standard-1"
	  boot_disk {
	    initialize_params {
	      image = "debian-cloud/debian-9"
	    }
	  }
	  network_interface {
	    subnetwork ="${google_compute_subnetwork.subnetwork.self_link}"
	  }
	}

	resource "google_compute_instance_iam_binding" "foo" {
		project = "${google_compute_instance.test_vm.project}"
	  zone = "${google_compute_instance.test_vm.zone}"
	  instance_name = "${google_compute_instance.test_vm.name}"
	  role        = "%s"
	  members     = ["serviceAccount:${google_service_account.test_account.email}"]
	}
	`, account, subnetworkName, subnetworkName, region, project, zone, instanceName, roleId)
}

func testAccComputeInstanceIamBinding_update(project, account, region, zone, subnetworkName, instanceName, roleId string) string {
	return fmt.Sprintf(`
	resource "google_service_account" "test_account" {
	  account_id   = "%s"
	  display_name = "Iam Testing Account"
	}

	resource "google_service_account" "test_account_2" {
	  account_id   = "%s-2"
	  display_name = "Iam Testing Account"
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

	resource "google_compute_instance" "test_vm" {
		project = "%s"
		zone         = "%s"
	  name         = "%s"
	  machine_type = "n1-standard-1"
	  boot_disk {
	    initialize_params {
	      image = "debian-cloud/debian-9"
	    }
	  }
	  network_interface {
	    subnetwork ="${google_compute_subnetwork.subnetwork.self_link}"
	  }
	}

	resource "google_compute_instance_iam_binding" "foo" {
		project = "${google_compute_instance.test_vm.project}"
	  zone = "${google_compute_instance.test_vm.zone}"
	  instance_name = "${google_compute_instance.test_vm.name}"
	  role        = "%s"
	  members      = [
	    "serviceAccount:${google_service_account.test_account.email}",
	    "serviceAccount:${google_service_account.test_account_2.email}"
	  ]
	}
	`, account, account, subnetworkName, subnetworkName, region, project, zone, instanceName, roleId)
}
