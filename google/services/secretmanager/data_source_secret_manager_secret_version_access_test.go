// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatasourceSecretManagerSecretVersionAccess_basic(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersionAccess_basic(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersionAccess("data.google_secret_manager_secret_version_access.basic", "1"),
				),
			},
		},
	})
}

func TestAccDatasourceSecretManagerSecretVersionAccess_latest(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersionAccess_latest(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersionAccess("data.google_secret_manager_secret_version_access.latest", "2"),
				),
			},
		},
	})
}

func testAccCheckDatasourceSecretManagerSecretVersionAccess(n, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Secret Version data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		version, ok := rs.Primary.Attributes["version"]
		if !ok {
			return errors.New("can't find 'version' attribute")
		}

		if version != expected {
			return fmt.Errorf("expected %s, got %s, version not found", expected, version)
		}
		return nil
	}
}

func testAccDatasourceSecretManagerSecretVersionAccess_latest(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  labels = {
    label = "my-label"
  }
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic-1" {
  secret = google_secret_manager_secret.secret-basic.name
  secret_data = "my-tf-test-secret-first"
}

resource "google_secret_manager_secret_version" "secret-version-basic-2" {
  secret = google_secret_manager_secret.secret-basic.name
  secret_data = "my-tf-test-secret-second"

  depends_on = [google_secret_manager_secret_version.secret-version-basic-1]
}

data "google_secret_manager_secret_version_access" "latest" {
  secret = google_secret_manager_secret_version.secret-version-basic-2.secret
}
`, randomString)
}

func testAccDatasourceSecretManagerSecretVersionAccess_basic(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  labels = {
    label = "my-label"
  }
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.name
  secret_data = "my-tf-test-secret-%s"
}

data "google_secret_manager_secret_version_access" "basic" {
  secret = google_secret_manager_secret_version.secret-version-basic.secret
  version = 1
}
`, randomString, randomString)
}
