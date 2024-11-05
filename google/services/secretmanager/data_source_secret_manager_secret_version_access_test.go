// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version_access.basic", "1"),
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version_access.basic", "google_secret_manager_secret_version.secret-version-basic"),
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
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version_access.latest", "2"),
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version_access.latest", "google_secret_manager_secret_version.secret-version-basic-2"),
				),
			},
		},
	})
}

func TestAccDatasourceSecretManagerSecretVersionAccess_withBase64SecretData(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	data := "./test-fixtures/binary-file.pfx"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersionAccess_withBase64SecretData(randomString, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version_access.basic-base64", "1"),
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version_access.basic-base64", "google_secret_manager_secret_version.secret-version-basic-base64"),
				),
			},
		},
	})
}

func testAccDatasourceSecretManagerSecretVersionAccess_latest(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  labels = {
    label = "my-label"
  }
  replication {
    auto {}
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
    auto {}
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

func testAccDatasourceSecretManagerSecretVersionAccess_withBase64SecretData(randomString, data string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret-basic-base64" {
  secret_id = "tf-test-secret-version-%s"
  labels = {
    label = "my-label"
  }
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic-base64" {
  secret = google_secret_manager_secret.secret-basic-base64.name
  is_secret_data_base64 = true
  secret_data = filebase64("%s")
}

data "google_secret_manager_secret_version_access" "basic-base64" {
  secret = google_secret_manager_secret_version.secret-version-basic-base64.secret
  is_secret_data_base64 = true
}
`, randomString, data)
}
