package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeInstanceIamBinding(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	role := "roles/compute.osLogin"
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceIamBinding_basic(zone, instanceName, role),
			},
			{
				ResourceName:      "google_compute_instance_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s %s", project, zone, instanceName, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccComputeInstanceIamBinding_update(zone, instanceName, role),
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
	role := "roles/compute.osLogin"
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccComputeInstanceIamMember_basic(zone, instanceName, role),
			},
			{
				ResourceName:      "google_compute_instance_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s %s user:admin@hashicorptest.com", project, zone, instanceName, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceIamPolicy(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	role := "roles/compute.osLogin"
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceIamPolicy_basic(zone, instanceName, role),
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

func testAccComputeInstanceIamMember_basic(zone, instanceName, roleId string) string {
	return fmt.Sprintf(`
  resource "google_compute_instance" "test_vm" {
    zone         = "%s"
    name         = "%s"
    machine_type = "n1-standard-1"

    boot_disk {
      initialize_params {
        image = "debian-cloud/debian-9"
      }
    }

    network_interface {
      network = "default"
    }
  }

  resource "google_compute_instance_iam_member" "foo" {
    project       = "${google_compute_instance.test_vm.project}"
    zone          = "${google_compute_instance.test_vm.zone}"
    instance_name = "${google_compute_instance.test_vm.name}"
    role          = "%s"
    member        = "user:Admin@hashicorptest.com"
  }

`, zone, instanceName, roleId)
}

func testAccComputeInstanceIamPolicy_basic(zone, instanceName, roleId string) string {
	return fmt.Sprintf(`
  resource "google_compute_instance" "test_vm" {
    zone         = "%s"
    name         = "%s"
    machine_type = "n1-standard-1"

    boot_disk {
      initialize_params {
        image = "debian-cloud/debian-9"
      }
    }

    network_interface {
      network = "default"
    }
  }

  data "google_iam_policy" "foo" {
    binding {
      role    = "%s"
      members = ["user:Admin@hashicorptest.com"]
    }
  }

  resource "google_compute_instance_iam_policy" "foo" {
    project       = "${google_compute_instance.test_vm.project}"
    zone          = "${google_compute_instance.test_vm.zone}"
    instance_name = "${google_compute_instance.test_vm.name}"
    policy_data   = "${data.google_iam_policy.foo.policy_data}"
  }

`, zone, instanceName, roleId)
}

func testAccComputeInstanceIamBinding_basic(zone, instanceName, roleId string) string {
	return fmt.Sprintf(`
  resource "google_compute_instance" "test_vm" {
    zone         = "%s"
    name         = "%s"
    machine_type = "n1-standard-1"

    boot_disk {
      initialize_params {
        image = "debian-cloud/debian-9"
      }
    }

    network_interface {
      network = "default"
    }
  }

  resource "google_compute_instance_iam_binding" "foo" {
    project       = "${google_compute_instance.test_vm.project}"
    zone          = "${google_compute_instance.test_vm.zone}"
    instance_name = "${google_compute_instance.test_vm.name}"
    role          = "%s"
    members       = ["user:Admin@hashicorptest.com"]
  }

`, zone, instanceName, roleId)
}

func testAccComputeInstanceIamBinding_update(zone, instanceName, roleId string) string {
	return fmt.Sprintf(`
  resource "google_compute_instance" "test_vm" {
    zone         = "%s"
    name         = "%s"
    machine_type = "n1-standard-1"

    boot_disk {
      initialize_params {
        image = "debian-cloud/debian-9"
      }
    }

    network_interface {
      network = "default"
    }
  }

  resource "google_compute_instance_iam_binding" "foo" {
    project       = "${google_compute_instance.test_vm.project}"
    zone          = "${google_compute_instance.test_vm.zone}"
    instance_name = "${google_compute_instance.test_vm.name}"
    role          = "%s"
    members       = ["user:Admin@hashicorptest.com", "user:paddy@hashicorp.com"]
  }

`, zone, instanceName, roleId)
}
