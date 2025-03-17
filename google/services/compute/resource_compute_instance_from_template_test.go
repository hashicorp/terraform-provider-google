// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeInstanceFromTemplate_basic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_basic(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "machine_type", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "attached_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scheduling.0.automatic_restart", "false"),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_self_link_unique(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_self_link_unique(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "machine_type", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "attached_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scheduling.0.automatic_restart", "false"),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_localSsdRecoveryTimeout(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	var expectedLocalSsdRecoveryTimeout = compute.Duration{}
	expectedLocalSsdRecoveryTimeout.Nanos = 0
	expectedLocalSsdRecoveryTimeout.Seconds = 3600

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_localSsdRecoveryTimeout(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					testAccCheckComputeInstanceLocalSsdRecoveryTimeout(&instance, expectedLocalSsdRecoveryTimeout),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplateWithOverride_localSsdRecoveryTimeout(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	var expectedLocalSsdRecoveryTimeout = compute.Duration{}
	expectedLocalSsdRecoveryTimeout.Nanos = 0
	expectedLocalSsdRecoveryTimeout.Seconds = 7200

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplateWithOverride_localSsdRecoveryTimeout(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					testAccCheckComputeInstanceLocalSsdRecoveryTimeout(&instance, expectedLocalSsdRecoveryTimeout),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_diskResourcePolicies(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_diskResourcePoliciesCreate(suffix, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstanceFromTemplate_diskResourcePoliciesUpdate(suffix, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.foobar", &instance),
				),
			},
			{
				Config:      testAccComputeInstanceFromTemplate_diskResourcePoliciesTwoPolicies(suffix, templateName),
				ExpectError: regexp.MustCompile("Too many list items"),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_overrideBootDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	overrideDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_overrideBootDisk(templateDisk, overrideDisk, templateName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "boot_disk.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "boot_disk.0.source", regexp.MustCompile(overrideDisk)),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_overrideAttachedDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	overrideDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_overrideAttachedDisk(templateDisk, overrideDisk, templateName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "attached_disk.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "attached_disk.0.source", regexp.MustCompile(overrideDisk)),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_overrideScratchDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	overrideDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_overrideScratchDisk(templateDisk, overrideDisk, templateName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.0.interface", "NVME"),
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.1.interface", "NVME"),
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.1.device_name", "override-local-ssd"),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_overrideScheduling(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_overrideScheduling(templateDisk, templateName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_TerminationTime(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateDisk := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"
	now := time.Now().UTC()
	terminationTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 9999, now.Location()).Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_terminationTime(templateDisk, templateName, terminationTime, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_overrideMetadataDotStartupScript(t *testing.T) {
	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_overrideMetadataDotStartupScript(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "metadata.startup-script", ""),
				),
			},
		},
	})
}

func TestAccComputeInstanceFromTemplate_useDiskSelfLink(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_regionalDisk(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceFromTemplateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance_from_template" {
				continue
			}

			_, err := config.NewComputeClient(config.UserAgent).Instances.Get(
				config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
			if err == nil {
				return fmt.Errorf("Instance still exists")
			}
		}

		return nil
	}
}

func TestAccComputeInstanceFromTemplate_confidentialInstanceConfigMain(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instance2 compute.Instance

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplate_confidentialInstanceConfigEnable(
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					"SEV"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst1", &instance),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance, true, "SEV"),
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst2", &instance2),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance2, true, ""),
				),
			},
			{
				Config: testAccComputeInstanceFromTemplate_confidentialInstanceConfigNoConfigSevSnp(
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					"SEV_SNP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst1", &instance),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance, false, "SEV_SNP"),
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst2", &instance2),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance2, false, "SEV_SNP"),
				),
			},
			{
				Config: testAccComputeInstanceFromTemplate_confidentialInstanceConfigNoConfigTdx(
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
					"TDX"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst1", &instance),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance, false, "TDX"),
					testAccCheckComputeInstanceExists(t, "google_compute_instance_from_template.inst2", &instance2),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance2, false, "TDX"),
				),
			},
		},
	})
}

func testAccComputeInstanceFromTemplate_basic(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  disk {
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    device_name  = "test-local-ssd"
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = false
    disk_type    = "pd-ssd"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = true
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

  // Overrides
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
  scheduling {
    automatic_restart = false
  }
}
`, template, template, instance)
}

func testAccComputeInstanceFromTemplate_localSsdRecoveryTimeout(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  disk {
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    device_name  = "test-local-ssd"
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = false
    disk_type    = "pd-ssd"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = true
    local_ssd_recovery_timeout {
			nanos = 0
			seconds = 3600
    }
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

  // Overrides
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
  scheduling {
    automatic_restart = false
  }
}
`, template, template, instance)
}

func testAccComputeInstanceFromTemplateWithOverride_localSsdRecoveryTimeout(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  disk {
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    device_name  = "test-local-ssd"
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = false
    disk_type    = "pd-ssd"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = true
    local_ssd_recovery_timeout {
			nanos = 0
			seconds = 3600
    }
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

  // Overrides
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
  scheduling {
    automatic_restart = false
    local_ssd_recovery_timeout {
			nanos = 0
			seconds = 7200
    }
  }
}
`, template, template, instance)
}

func testAccComputeInstanceFromTemplate_regionalDisk(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_region_disk" "foobar" {
  name          = "%s"
  size          = 10
  type          = "pd-ssd"
  region        = "us-central1"
  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
    disk_type    = "pd-ssd"
    type         = "PERSISTENT"
  }

  disk {
    source      = google_compute_region_disk.foobar.self_link
    auto_delete = false
    boot        = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = true
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.id

  // Overrides
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
  scheduling {
    automatic_restart = false
  }
}
`, template, template, instance)
}

func testAccComputeInstanceFromTemplate_self_link_unique(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  disk {
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    device_name  = "test-local-ssd"
    disk_type    = "local-ssd"
    type         = "SCRATCH"
    interface    = "NVME"
    disk_size_gb = 375
  }

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = false
    disk_type    = "pd-ssd"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = true
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link_unique

  // Overrides
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
  scheduling {
    automatic_restart = false
  }
}
`, template, template, instance)
}

func testAccComputeInstanceFromTemplate_overrideBootDisk(templateDisk, overrideDisk, template, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "template_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_disk" "override_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 20
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "template" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  disk {
    source      = google_compute_disk.template_disk.name
    auto_delete = false
    boot        = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.template.self_link

  // Overrides
  boot_disk {
    source = google_compute_disk.override_disk.self_link
  }
}
`, templateDisk, overrideDisk, template, instance)
}

func testAccComputeInstanceFromTemplate_overrideAttachedDisk(templateDisk, overrideDisk, template, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "template_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_disk" "override_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 20
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "template" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  disk {
    source      = google_compute_disk.template_disk.name
    auto_delete = false
    boot        = false
  }

  disk {
    source_image = "debian-cloud/debian-11"
    auto_delete  = true
    boot         = false
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.template.self_link

  // Overrides
  attached_disk {
    source = google_compute_disk.override_disk.name
  }
}
`, templateDisk, overrideDisk, template, instance)
}

func testAccComputeInstanceFromTemplate_overrideScratchDisk(templateDisk, overrideDisk, template, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "template_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_disk" "override_disk" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 20
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "template" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  disk {
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    disk_size_gb = 375
    interface    = "SCSI"
    auto_delete  = true
    boot         = false
  }

  disk {
    device_name  = "test-local-ssd"
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    disk_size_gb = 375
    interface    = "SCSI"
    auto_delete  = true
    boot         = false
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.template.self_link

  // Overrides
  scratch_disk {
    interface = "NVME"
  }

  scratch_disk {
    device_name = "override-local-ssd"
    interface   = "NVME"
  }
}
`, templateDisk, overrideDisk, template, instance)
}

func testAccComputeInstanceFromTemplate_overrideScheduling(templateDisk, template, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    preemptible = true
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link
}
`, templateDisk, template, instance)
}

func testAccComputeInstanceFromTemplate_terminationTime(templateDisk, template, termination_time, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    instance_termination_action = "STOP"
    termination_time = "%s"
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

  scheduling {
    instance_termination_action = "STOP"
    termination_time = "%s"
  }
}
`, templateDisk, template, termination_time, instance, termination_time)
}

func testAccComputeInstanceFromTemplate_overrideMetadataDotStartupScript(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    startup-script = "#!/bin/bash\necho Hello"
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

  // Overrides
  metadata = {
    startup-script = ""
  }
}
`, template, instance)
}

func testAccComputeInstanceFromTemplate_confidentialInstanceConfigEnable(templateDisk string, image string, template string, instance string, template2 string, instance2 string, confidentialInstanceType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image1" {
  family  = "ubuntu-2004-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_disk" "foobar1" {
  name  = "%s"
  image = data.google_compute_image.my_image1.self_link
  size  = 10
  type  = "pd-standard"
  zone  = "us-central1-a"
}

resource "google_compute_image" "foobar1" {
  name              = "%s"
  source_disk       = google_compute_disk.foobar1.self_link
}

resource "google_compute_instance_template" "foobar1" {
  name         = "%s"
  machine_type = "n2d-standard-2"

  disk {
    source_image = google_compute_image.foobar1.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    enable_confidential_compute  = true
    confidential_instance_type     = %q
  }
}

resource "google_compute_instance_from_template" "inst1" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar1.self_link
}

resource "google_compute_instance_template" "foobar2" {
  name         = "%s"
  machine_type = "n2d-standard-2"

  disk {
    source_image = google_compute_image.foobar1.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    enable_confidential_compute  = true
  }
}

resource "google_compute_instance_from_template" "inst2" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar2.self_link
}
`, templateDisk, image, template, confidentialInstanceType, instance, template2, instance2)
}

func testAccComputeInstanceFromTemplate_confidentialInstanceConfigNoConfigSevSnp(templateDisk string, image string, template string, instance string, template2 string, instance2 string, confidentialInstanceType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image1" {
  family  = "ubuntu-2004-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_disk" "foobar1" {
  name  = "%s"
  image = data.google_compute_image.my_image1.self_link
  size  = 10
  type  = "pd-standard"
  zone  = "us-central1-a"
}

resource "google_compute_image" "foobar1" {
  name              = "%s"
  source_disk       = google_compute_disk.foobar1.self_link
}

resource "google_compute_instance_template" "foobar3" {
  name         = "%s"
  machine_type = "n2d-standard-2"

  disk {
    source_image = google_compute_image.foobar1.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    enable_confidential_compute    = false
    confidential_instance_type     = %q
  }
}

resource "google_compute_instance_from_template" "inst1" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar3.self_link
}

resource "google_compute_instance_template" "foobar4" {
  name         = "%s"
  machine_type = "n2d-standard-2"

  disk {
    source_image = google_compute_image.foobar1.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    confidential_instance_type     = %q
  }
}

resource "google_compute_instance_from_template" "inst2" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar4.self_link
}
`, templateDisk, image, template, confidentialInstanceType, instance, template2, confidentialInstanceType, instance2)
}

func testAccComputeInstanceFromTemplate_confidentialInstanceConfigNoConfigTdx(templateDisk string, image string, template string, instance string, template2 string, instance2 string, confidentialInstanceType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image2" {
  family  = "ubuntu-2204-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_disk" "foobar2" {
  name  = "%s"
  image = data.google_compute_image.my_image2.self_link
  size  = 10
  type  = "pd-balanced"
  zone  = "us-central1-a"
}

resource "google_compute_image" "foobar2" {
  name              = "%s"
  source_disk       = google_compute_disk.foobar2.self_link
}

resource "google_compute_instance_template" "foobar5" {
  name         = "%s"
  machine_type = "c3-standard-4"

  disk {
    source_image = google_compute_image.foobar2.name
    auto_delete  = true
    boot         = true
    disk_type    = "pd-balanced"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    enable_confidential_compute    = false
    confidential_instance_type     = %q
  }
}

resource "google_compute_instance_from_template" "inst1" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar5.self_link
}

resource "google_compute_instance_template" "foobar6" {
  name         = "%s"
  machine_type = "c3-standard-4"

  disk {
    source_image = google_compute_image.foobar2.name
    auto_delete  = true
    boot         = true
    disk_type    = "pd-balanced"
    type         = "PERSISTENT"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  scheduling {
    automatic_restart = false
    on_host_maintenance = "TERMINATE"
  }

  confidential_instance_config {
    confidential_instance_type = %q
  }
}

resource "google_compute_instance_from_template" "inst2" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar6.self_link
}
`, templateDisk, image, template, confidentialInstanceType, instance, template2, confidentialInstanceType, instance2)
}

func TestAccComputeInstanceFromTemplateWithOverride_interface(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceFromTemplateWithOverride_interface(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "boot_disk.0.interface", "SCSI"),
				),
			},
		},
	})
}

func testAccComputeInstanceFromTemplateWithOverride_interface(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobarboot" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_disk" "foobarattach" {
  name = "%s"
  size = 100
  type = "pd-balanced"
  zone = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"  // can't be e2 because of local-ssd

  disk {
    source      = google_compute_disk.foobarboot.name
    auto_delete = false
    boot        = true
  }


  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.foobar.self_link

   attached_disk {
    source =  google_compute_disk.foobarattach.name
  }
   // Overrides
  boot_disk {
    interface = "SCSI"
    source =  google_compute_disk.foobarboot.name
  }
}
`, template, instance, template, instance)
}

func testAccComputeInstanceFromTemplate_diskResourcePoliciesCreate(suffix, template string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "test-snapshot-policy" {
  name    = "test-policy-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "11:00"
      }
    }
  }
}

resource "google_compute_resource_policy" "test-snapshot-policy2" {
  name    = "test-policy2-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "22:00"
      }
    }
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_region_instance_template" "foobar" {
  name         = "%s"
  region = "us-central1"
  machine_type = "n1-standard-1"
  disk {
    resource_policies = [ google_compute_resource_policy.test-snapshot-policy.name ]
    source_image = data.google_compute_image.my_image.self_link
  }
  network_interface {
      network = "default"
  }
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"
  source_instance_template = google_compute_region_instance_template.foobar.id
}
`, suffix, suffix, template, template)
}

func testAccComputeInstanceFromTemplate_diskResourcePoliciesUpdate(suffix, template string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "test-snapshot-policy" {
  name    = "test-policy-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "11:00"
      }
    }
  }
}

resource "google_compute_resource_policy" "test-snapshot-policy2" {
  name    = "test-policy2-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "22:00"
      }
    }
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_region_instance_template" "foobar" {
  name         = "%s"
  region = "us-central1"
  machine_type = "n1-standard-1"
  disk {
    resource_policies = [ google_compute_resource_policy.test-snapshot-policy2.name ]
    source_image = data.google_compute_image.my_image.self_link
  }
  network_interface {
      network = "default"
  }
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"
  source_instance_template = google_compute_region_instance_template.foobar.id
}
`, suffix, suffix, template, template)
}

func testAccComputeInstanceFromTemplate_diskResourcePoliciesTwoPolicies(suffix, template string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "test-snapshot-policy" {
  name    = "test-policy-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "11:00"
      }
    }
  }
}

resource "google_compute_resource_policy" "test-snapshot-policy2" {
  name    = "test-policy2-%s"
  snapshot_schedule_policy {
    schedule {
      hourly_schedule {
        hours_in_cycle = 1
        start_time     = "22:00"
      }
    }
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_region_instance_template" "foobar" {
  name         = "%s"
  region = "us-central1"
  machine_type = "n1-standard-1"
  disk {
    resource_policies = [ google_compute_resource_policy.test-snapshot-policy.name, google_compute_resource_policy.test-snapshot-policy2.name ]
    source_image = data.google_compute_image.my_image.self_link
  }
  network_interface {
      network = "default"
  }
}

resource "google_compute_instance_from_template" "foobar" {
  name = "%s"
  zone = "us-central1-a"
  source_instance_template = google_compute_region_instance_template.foobar.id
}
  `, suffix, suffix, template, template)
}
