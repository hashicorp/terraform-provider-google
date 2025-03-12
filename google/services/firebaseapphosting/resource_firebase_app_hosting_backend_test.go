// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firebaseapphosting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseAppHostingBackend_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirebaseAppHostingBackendDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseAppHostingBackend_firebaseAppHostingBackendBefore(context),
			},
			{
				ResourceName:            "google_firebase_app_hosting_backend.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "backend_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccFirebaseAppHostingBackend_firebaseAppHostingBackendAfter(context),
			},
			{
				ResourceName:            "google_firebase_app_hosting_backend.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "backend_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccFirebaseAppHostingBackend_firebaseAppHostingBackendBefore(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  project = "%{project_id}"

  # Must be firebase-app-hosting-compute
  account_id                   = "firebase-app-hosting-compute"
  display_name                 = "Firebase App Hosting compute service account"

  # Do not throw if already exists
  create_ignore_already_exists = true
}

resource "google_firebase_app_hosting_backend" "example" {
  project          = "%{project_id}"
  # Choose the region closest to your users
  location         = "us-central1"
  backend_id       = "tf-test-%{random_suffix}"
  app_id           = "1:0000000000:web:674cde32020e16fbce9dbe"
  display_name     = "My Backend After"
  serving_locality = "GLOBAL_ACCESS"
  service_account  = google_service_account.service_account.email
  environment      = "staging"

  annotations = {
    "key" = "before"
  }

  labels = {
    "key" = "before"
  }
}
`, context)
}

func testAccFirebaseAppHostingBackend_firebaseAppHostingBackendAfter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  project = "%{project_id}"

  # Must be firebase-app-hosting-compute
  account_id                   = "firebase-app-hosting-compute"
  display_name                 = "Firebase App Hosting compute service account"

  # Do not throw if already exists
  create_ignore_already_exists = true
}

resource "google_firebase_app_hosting_backend" "example" {
  project          = "%{project_id}"
  # Choose the region closest to your users
  location         = "us-central1"
  backend_id       = "tf-test-%{random_suffix}"
  app_id           = "1:0000000000:web:674cde32020e16fbce9dbd"
  display_name     = "My Backend After"
  serving_locality = "GLOBAL_ACCESS"
  service_account  = google_service_account.service_account.email
  environment      = "prod"

  annotations = {
    "key" = "after"
  }

  labels = {
    "key" = "after"
  }
}
`, context)
}
