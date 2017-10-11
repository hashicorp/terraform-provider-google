package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeImage_basic(t *testing.T) {
	var image compute.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeImage_basic("image-test-" + acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						"google_compute_image.foobar", &image),
					testAccCheckComputeImageDescription(&image, "description-test"),
					testAccCheckComputeImageFamily(&image, "family-test"),
					testAccCheckComputeImageContainsLabel(&image, "my-label", "my-label-value"),
					testAccCheckComputeImageContainsLabel(&image, "empty-label", ""),
					testAccCheckComputeImageHasComputedFingerprint(&image, "google_compute_image.foobar"),
				),
			},
		},
	})
}

func TestAccComputeImage_update(t *testing.T) {
	var image compute.Image

	name := "image-test-" + acctest.RandString(10)
	// Only labels supports an update
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeImage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						"google_compute_image.foobar", &image),
					testAccCheckComputeImageContainsLabel(&image, "my-label", "my-label-value"),
					testAccCheckComputeImageContainsLabel(&image, "empty-label", ""),
					testAccCheckComputeImageHasComputedFingerprint(&image, "google_compute_image.foobar"),
				),
			},
			resource.TestStep{
				Config: testAccComputeImage_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						"google_compute_image.foobar", &image),
					testAccCheckComputeImageDoesNotContainLabel(&image, "my-label"),
					testAccCheckComputeImageContainsLabel(&image, "empty-label", "oh-look-theres-a-label-now"),
					testAccCheckComputeImageContainsLabel(&image, "new-field", "only-shows-up-when-updated"),
					testAccCheckComputeImageHasComputedFingerprint(&image, "google_compute_image.foobar"),
				),
			},
		},
	})
}

func TestAccComputeImage_basedondisk(t *testing.T) {
	var image compute.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeImage_basedondisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						"google_compute_image.foobar", &image),
					testAccCheckComputeImageHasSourceDisk(&image),
				),
			},
		},
	})
}

func testAccCheckComputeImageDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_image" {
			continue
		}

		_, err := config.clientCompute.Images.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckComputeImageExists(n string, image *compute.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Images.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		*image = *found

		return nil
	}
}

func testAccCheckComputeImageDescription(image *compute.Image, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if image.Description != description {
			return fmt.Errorf("Wrong image description: expected '%s' got '%s'", description, image.Description)
		}
		return nil
	}
}

func testAccCheckComputeImageFamily(image *compute.Image, family string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if image.Family != family {
			return fmt.Errorf("Wrong image family: expected '%s' got '%s'", family, image.Family)
		}
		return nil
	}
}

func testAccCheckComputeImageContainsLabel(image *compute.Image, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := image.Labels[key]
		if !ok {
			return fmt.Errorf("Expected label with key '%s' not found", key)
		}
		if v != value {
			return fmt.Errorf("Incorrect label value for key '%s': expected '%s' but found '%s'", key, value, v)
		}
		return nil
	}
}

func testAccCheckComputeImageDoesNotContainLabel(image *compute.Image, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := image.Labels[key]; ok {
			return fmt.Errorf("Expected no label for key '%s' but found one with value '%s'", key, v)
		}

		return nil
	}
}

func testAccCheckComputeImageHasComputedFingerprint(image *compute.Image, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// First ensure we actually have a fingerprint
		if image.LabelFingerprint == "" {
			return fmt.Errorf("No fingerprint set in API read result")
		}

		state := s.RootModule().Resources[resource]
		if state == nil {
			return fmt.Errorf("Unable to find resource named %s in resources", resource)
		}

		storedFingerprint := state.Primary.Attributes["label_fingerprint"]
		if storedFingerprint != image.LabelFingerprint {
			return fmt.Errorf("Stored fingerprint doesn't match fingerprint found on server; stored '%s', server '%s'",
				storedFingerprint, image.LabelFingerprint)
		}

		return nil
	}
}

func testAccCheckComputeImageHasSourceDisk(image *compute.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if image.SourceType == "" {
			return fmt.Errorf("No source disk")
		}
		return nil
	}
}

func testAccComputeImage_basic(name string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "foobar" {
	name = "%s"
	description = "description-test"
	family = "family-test"
	raw_disk {
	  source = "https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz"
	}
	create_timeout = 5
	labels = {
		my-label = "my-label-value"
		empty-label = ""
	}
}`, name)
}

func testAccComputeImage_update(name string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "foobar" {
	name = "%s"
	description = "description-test"
	family = "family-test"
	raw_disk {
	  source = "https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.4-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz"
	}
	create_timeout = 5
	labels = {
		empty-label = "oh-look-theres-a-label-now"
		new-field = "only-shows-up-when-updated"
	}
}`, name)
}

var testAccComputeImage_basedondisk = fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "disk-test-%s"
	zone = "us-central1-a"
	image = "debian-8-jessie-v20160803"
}
resource "google_compute_image" "foobar" {
	name = "image-test-%s"
	source_disk = "${google_compute_disk.foobar.self_link}"
}`, acctest.RandString(10), acctest.RandString(10))
