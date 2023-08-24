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

func TestAccDatasourceSecretManagerSecretVersion_basic(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersion_basic(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version.basic", "1"),
				),
			},
		},
	})
}

func TestAccDatasourceSecretManagerSecretVersion_latest(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersion_latest(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version.latest", "2"),
				),
			},
		},
	})
}

func testAccCheckDatasourceSecretManagerSecretVersion(n, expected string) resource.TestCheckFunc {
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

func testAccDatasourceSecretManagerSecretVersion_latest(randomString string) string {
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

data "google_secret_manager_secret_version" "latest" {
  secret = google_secret_manager_secret_version.secret-version-basic-2.secret
}
`, randomString)
}

func testAccDatasourceSecretManagerSecretVersion_basic(randomString string) string {
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

data "google_secret_manager_secret_version" "basic" {
  secret = google_secret_manager_secret_version.secret-version-basic.secret
  version = 1
}
`, randomString, randomString)
}
