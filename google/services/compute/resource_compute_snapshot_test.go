// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeSnapshot_encryption(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryption(snapshotName, diskName),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_encryption_key", "source_disk", "source_disk_encryption_key", "zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_encryptionCMEK(t *testing.T) {
	t.Parallel()
	// KMS causes errors due to rotation
	acctest.SkipIfVcr(t)

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryptionCMEK(snapshotName, diskName),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "snapshot_encryption_key", "source_disk_encryption_key"},
			},
		},
	})
}

func testAccComputeSnapshot_encryption(snapshotName string, diskName string) string {
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
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }

  source_disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
`, diskName, snapshotName)
}

func testAccComputeSnapshot_encryptionCMEK(snapshotName string, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-10"
  project = "debian-cloud"
}

resource "google_service_account" "test" {
	account_id   = "%s"
	display_name = "KMS Ops Account"
}

resource "google_kms_key_ring" "keyring" {
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "example-key" {
  name            = "%s"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"
}

resource "google_kms_crypto_key_iam_member" "example-key" {
  crypto_key_id = google_kms_crypto_key.example-key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${google_service_account.test.email}"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
	kms_key_self_link = google_kms_crypto_key_iam_member.example-key.crypto_key_id
	kms_key_service_account = google_service_account.test.email
  }
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
	kms_key_self_link = google_kms_crypto_key_iam_member.example-key.crypto_key_id
	kms_key_service_account = google_service_account.test.email
  }
}
`, diskName, diskName, diskName, diskName, snapshotName)
}
