resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudkms.googleapis.com",
  ]
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project_services.acceptance.project}"
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = "${google_kms_key_ring.key_ring.id}"
}

# rotation
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudkms.googleapis.com",
  ]
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project.acceptance.project_id}"
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name            = "%s"
  key_ring        = "${google_kms_key_ring.key_ring.id}"
  rotation_period = "%s"
}

# remove
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudkms.googleapis.com",
  ]
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project.acceptance.project_id}"
  name     = "%s"
  location = "us-central1"
}
