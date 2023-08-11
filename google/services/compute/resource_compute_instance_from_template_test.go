// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.0.interface", "NVME"),
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

func TestAccComputeInstanceFromTemplate_012_removableFields(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	templateName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_instance_from_template.inst"

	// First config is a basic instance from template, second tests the empty list syntax
	config1 := testAccComputeInstanceFromTemplate_012_removableFieldsTpl(templateName) +
		testAccComputeInstanceFromTemplate_012_removableFields1(instanceName)
	config2 := testAccComputeInstanceFromTemplate_012_removableFieldsTpl(templateName) +
		testAccComputeInstanceFromTemplate_012_removableFields2(instanceName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceFromTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					resource.TestCheckResourceAttr(resourceName, "service_account.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_account.0.scopes.#", "3"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, resourceName, &instance),

					// Check that fields were able to be removed
					resource.TestCheckResourceAttr(resourceName, "scratch_disk.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "attached_disk.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.alias_ip_range.#", "0"),
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

func testAccComputeInstanceFromTemplate_012_removableFieldsTpl(template string) string {

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
    disk_size_gb = 20
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  can_ip_forward = true
}
`, template)
}

func testAccComputeInstanceFromTemplate_012_removableFields1(instance string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  allow_stopping_for_update = true

  source_instance_template = google_compute_instance_template.foobar.self_link
}
`, instance)
}

func testAccComputeInstanceFromTemplate_012_removableFields2(instance string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_from_template" "inst" {
  name = "%s"
  zone = "us-central1-a"

  allow_stopping_for_update = true

  source_instance_template = google_compute_instance_template.foobar.self_link

  // Overrides
  network_interface {
    alias_ip_range = []
  }

  service_account = []

  scratch_disk = []

  attached_disk = []

  timeouts {
    create = "10m"
    update = "10m"
  }
}
`, instance)
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
