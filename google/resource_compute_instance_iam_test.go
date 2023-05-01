package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Even though the resource has generated tests, keep this one around until we are able to generate
// checking the different import formats
func TestAccComputeInstanceIamPolicy(t *testing.T) {
	t.Parallel()

	project := acctest.GetTestProjectFromEnv()
	role := "roles/compute.osLogin"
	zone := acctest.GetTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
			{
				ResourceName:      "google_compute_instance_iam_policy.foo",
				ImportStateId:     instanceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeInstanceIamPolicy_basic(zone, instanceName, roleId string) string {
	return fmt.Sprintf(`
  resource "google_compute_instance" "test_vm" {
    zone         = "%s"
    name         = "%s"
    machine_type = "e2-medium"

    boot_disk {
      initialize_params {
        image = "debian-cloud/debian-11"
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
