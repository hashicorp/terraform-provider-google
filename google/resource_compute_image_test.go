package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeImage_withLicense(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_license("image-test-" + randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_image.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeImage_update(t *testing.T) {
	t.Parallel()

	var image compute.Image

	name := "image-test-" + randString(t, 10)
	// Only labels supports an update
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageContainsLabel(&image, "my-label", "my-label-value"),
					testAccCheckComputeImageContainsLabel(&image, "empty-label", ""),
				),
			},
			{
				Config: testAccComputeImage_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageDoesNotContainLabel(&image, "my-label"),
					testAccCheckComputeImageContainsLabel(&image, "empty-label", "oh-look-theres-a-label-now"),
					testAccCheckComputeImageContainsLabel(&image, "new-field", "only-shows-up-when-updated"),
				),
			},
			{
				ResourceName:            "google_compute_image.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"raw_disk"},
			},
		},
	})
}

func TestAccComputeImage_basedondisk(t *testing.T) {
	t.Parallel()

	var image compute.Image

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_basedondisk(randString(t, 10), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageHasSourceType(&image),
				),
			},
			{
				ResourceName:      "google_compute_image.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeImage_sourceImage(t *testing.T) {
	t.Parallel()

	var image compute.Image
	imageName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_sourceImage(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageHasSourceType(&image),
				),
			},
			{
				ResourceName:      "google_compute_image.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeImage_sourceSnapshot(t *testing.T) {
	t.Parallel()

	var image compute.Image

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	snapshotName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	imageName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_sourceSnapshot(diskName, snapshotName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageHasSourceType(&image),
				),
			},
			{
				ResourceName:      "google_compute_image.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeImageExists(t *testing.T, n string, image *compute.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No name is set")
		}

		config := googleProviderConfig(t)

		found, err := config.NewComputeClient(config.userAgent).Images.Get(
			config.Project, rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Image not found")
		}

		*image = *found

		return nil
	}
}

func TestAccComputeImage_resolveImage(t *testing.T) {
	t.Parallel()

	var image compute.Image
	rand := randString(t, 10)
	name := fmt.Sprintf("test-image-%s", rand)
	fam := fmt.Sprintf("test-image-family-%s", rand)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_resolving(name, fam),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeImageExists(
						t, "google_compute_image.foobar", &image),
					testAccCheckComputeImageResolution(t, "google_compute_image.foobar"),
				),
			},
		},
	})
}

func TestAccComputeImage_imageEncryptionKey(t *testing.T) {
	t.Parallel()

	kmsKey := BootstrapKMSKeyInLocation(t, "us-central1")
	kmsKeyName := GetResourceNameFromSelfLink(kmsKey.CryptoKey.Name)
	kmsRingName := GetResourceNameFromSelfLink(kmsKey.KeyRing.Name)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_imageEncryptionKey(kmsRingName, kmsKeyName, randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_image.image",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeImageResolution(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		project := config.Project

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No image name is set")
		}
		if rs.Primary.Attributes["family"] == "" {
			return fmt.Errorf("No image family is set")
		}
		if rs.Primary.Attributes["self_link"] == "" {
			return fmt.Errorf("No self_link is set")
		}

		name := rs.Primary.Attributes["name"]
		family := rs.Primary.Attributes["family"]
		link := rs.Primary.Attributes["self_link"]

		latestDebian, err := config.NewComputeClient(config.userAgent).Images.GetFromFamily("debian-cloud", "debian-11").Do()
		if err != nil {
			return fmt.Errorf("Error retrieving latest debian: %s", err)
		}

		images := map[string]string{
			"family/" + latestDebian.Family:                            "projects/debian-cloud/global/images/family/" + latestDebian.Family,
			"projects/debian-cloud/global/images/" + latestDebian.Name: "projects/debian-cloud/global/images/" + latestDebian.Name,
			latestDebian.Family:                                        "projects/debian-cloud/global/images/family/" + latestDebian.Family,
			latestDebian.Name:                                          "projects/debian-cloud/global/images/" + latestDebian.Name,
			latestDebian.SelfLink:                                      latestDebian.SelfLink,

			"global/images/" + name:          "global/images/" + name,
			"global/images/family/" + family: "global/images/family/" + family,
			name:                             "global/images/" + name,
			family:                           "global/images/family/" + family,
			"family/" + family:               "global/images/family/" + family,
			project + "/" + name:             "projects/" + project + "/global/images/" + name,
			project + "/" + family:           "projects/" + project + "/global/images/family/" + family,
			link:                             link,
		}

		for input, expectation := range images {
			result, err := resolveImage(config, project, input, config.userAgent)
			if err != nil {
				return fmt.Errorf("Error resolving input %s to image: %+v\n", input, err)
			}
			if result != expectation {
				return fmt.Errorf("Expected input '%s' to resolve to '%s', it resolved to '%s' instead.\n", input, expectation, result)
			}
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

func testAccCheckComputeImageHasSourceType(image *compute.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if image.SourceType == "" {
			return fmt.Errorf("No source disk")
		}
		return nil
	}
}

func testAccComputeImage_resolving(name, family string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "foobar" {
  name        = "%s"
  family      = "%s"
  source_disk = google_compute_disk.foobar.self_link
}
`, name, name, family)
}

func testAccComputeImage_basic(name string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "foobar" {
  name        = "%s"
  description = "description-test"
  family      = "family-test"
  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
  labels = {
    my-label    = "my-label-value"
    empty-label = ""
  }
}
`, name)
}

func testAccComputeImage_license(name string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "disk-test-%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "foobar" {
  name        = "%s"
  description = "description-test"
  source_disk = google_compute_disk.foobar.self_link

  labels = {
    my-label    = "my-label-value"
    empty-label = ""
  }
  licenses = [
    "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/licenses/debian-11-bullseye",
  ]
}
`, name, name)
}

func testAccComputeImage_update(name string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "foobar" {
  name        = "%s"
  description = "description-test"
  family      = "family-test"
  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
  labels = {
    empty-label = "oh-look-theres-a-label-now"
    new-field   = "only-shows-up-when-updated"
  }
}
`, name)
}

func testAccComputeImage_basedondisk(diskName, imageName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "disk-test-%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "foobar" {
  name        = "image-test-%s"
  source_disk = google_compute_disk.foobar.self_link
}
`, diskName, imageName)
}

func testAccComputeImage_sourceImage(imageName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_image" "foobar" {
  name         = "%s"
  source_image = data.google_compute_image.my_image.self_link
}
`, imageName)
}

func testAccComputeImage_sourceSnapshot(diskName, snapshotName, imageName string) string {
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

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
}

resource "google_compute_image" "foobar" {
  name            = "%s"
  source_snapshot = google_compute_snapshot.foobar.self_link
}
`, diskName, snapshotName, imageName)
}

func testAccComputeImage_imageEncryptionKey(kmsRingName, kmsKeyName, suffix string) string {
	return fmt.Sprintf(`
data "google_kms_key_ring" "ring" {
  name     = "%s"
  location = "us-central1"
}

data "google_kms_crypto_key" "key" {
  name     = "%s"
  key_ring = data.google_kms_key_ring.ring.id
}

resource "google_service_account" "test" {
  account_id   = "tf-test-sa-%s"
  display_name = "KMS Ops Account"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = data.google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${google_service_account.test.email}"
}

data "google_compute_image" "debian" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_image" "image" {
  name         = "tf-test-image-%s"
  source_image = data.google_compute_image.debian.self_link
  image_encryption_key {
    kms_key_self_link       = data.google_kms_crypto_key.key.id
    kms_key_service_account = google_service_account.test.email
  }
}
`, kmsRingName, kmsKeyName, suffix, suffix)
}
