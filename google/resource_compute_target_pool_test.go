package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeTargetPool_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetPool_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetPoolExists(
						"google_compute_target_pool.foo"),
					testAccCheckComputeTargetPoolHealthCheck("google_compute_target_pool.foo", "google_compute_http_health_check.foobar"),
					testAccCheckComputeTargetPoolExists(
						"google_compute_target_pool.bar"),
					testAccCheckComputeTargetPoolHealthCheck("google_compute_target_pool.bar", "google_compute_http_health_check.foobar"),
				),
			},
			{
				ResourceName:      "google_compute_target_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeTargetPool_update(t *testing.T) {
	t.Parallel()

	tpname := fmt.Sprintf("tptest-%s", acctest.RandString(10))
	name1 := fmt.Sprintf("tptest-%s", acctest.RandString(10))
	name2 := fmt.Sprintf("tptest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetPoolDestroy,
		Steps: []resource.TestStep{
			{
				// Create target pool with no instances attached
				Config: testAccComputeTargetPool_update(tpname, "", name1, name2),
			},
			{
				ResourceName:      "google_compute_target_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add the two instances to the pool
				Config: testAccComputeTargetPool_update(tpname,
					`"${google_compute_instance.foo.self_link}", "${google_compute_instance.bar.self_link}"`,
					name1, name2),
			},
			{
				ResourceName:      "google_compute_target_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Reversing the order of instances or changing import format shouldn't matter
				Config: testAccComputeTargetPool_update(tpname,
					fmt.Sprintf(`"${google_compute_instance.bar.self_link}", "us-central1-a/%s"`, name1),
					name1, name2),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckComputeTargetPoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_pool" {
			continue
		}

		_, err := config.clientCompute.TargetPools.Get(
			config.Project, config.Region, rs.Primary.Attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("TargetPool still exists")
		}
	}

	return nil
}

func testAccCheckComputeTargetPoolExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.TargetPools.Get(
			config.Project, config.Region, rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("TargetPool not found")
		}

		return nil
	}
}

func testAccCheckComputeTargetPoolHealthCheck(targetPool, healthCheck string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetPoolRes, ok := s.RootModule().Resources[targetPool]
		if !ok {
			return fmt.Errorf("Not found: %s", targetPool)
		}

		healthCheckRes, ok := s.RootModule().Resources[healthCheck]
		if !ok {
			return fmt.Errorf("Not found: %s", healthCheck)
		}

		hcLink := healthCheckRes.Primary.Attributes["self_link"]
		if targetPoolRes.Primary.Attributes["health_checks.0"] != hcLink {
			return fmt.Errorf("Health check not set up. Expected %q to equal %q", targetPoolRes.Primary.Attributes["health_checks.0"], hcLink)
		}

		return nil
	}
}

func testAccComputeTargetPool_basic() string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_http_health_check" "foobar" {
	name = "healthcheck-test-%s"
	host = "example.com"
}

resource "google_compute_instance" "foobar" {
	name         = "inst-tp-test-%s"
	machine_type = "n1-standard-1"
	zone         = "us-central1-a"

	boot_disk {
		initialize_params{
			image = "${data.google_compute_image.my_image.self_link}"
		}
	}

	network_interface {
		network = "default"
	}
}

resource "google_compute_target_pool" "foo" {
	description = "Resource created for Terraform acceptance testing"
	instances = ["${google_compute_instance.foobar.self_link}", "us-central1-b/bar"]
	name = "tpool-test-%s"
	session_affinity = "CLIENT_IP_PROTO"
	health_checks = [
		"${google_compute_http_health_check.foobar.name}"
	]
}

resource "google_compute_target_pool" "bar" {
	description = "Resource created for Terraform acceptance testing"
	name = "tpool-test-%s"
	health_checks = [
		"${google_compute_http_health_check.foobar.self_link}"
	]
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeTargetPool_update(tpname, instances, name1, name2 string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo" {
	description = "Resource created for Terraform acceptance testing"
	name = "tpool-test-%s"
	instances = [%s]
}

resource "google_compute_instance" "foo" {
	name         = "%s"
	machine_type = "n1-standard-1"
	zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

	network_interface {
		network = "default"
	}
}

resource "google_compute_instance" "bar" {
	name         = "%s"
	machine_type = "n1-standard-1"
	zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

	network_interface {
		network = "default"
	}
}
`, tpname, instances, name1, name2)
}
