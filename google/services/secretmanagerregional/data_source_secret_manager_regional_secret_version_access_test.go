// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithResourceReference(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithResourceReference(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version_access.basic-1", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithSecretName(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithSecretName(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version_access.basic-2", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_latest(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_latest(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version_access.latest-1", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_versionField(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_versionField(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version_access.version-access", "1"),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithResourceReference(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-%s"
}

data "google_secret_manager_regional_secret_version_access" "basic-1" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic.secret
}
`, randomString, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_basicWithSecretName(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-%s"
}

data "google_secret_manager_regional_secret_version_access" "basic-2" {
  secret = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret_version.secret-version-basic.location
}
`, randomString, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_latest(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic-1" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-first"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic-2" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-second"

  depends_on = [google_secret_manager_regional_secret_version.secret-version-basic-1]
}

data "google_secret_manager_regional_secret_version_access" "latest-1" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic-2.secret
}
`, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersionAccess_versionField(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic-1" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-first"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic-2" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-second"

  depends_on = [google_secret_manager_regional_secret_version.secret-version-basic-1]
}

data "google_secret_manager_regional_secret_version_access" "version-access" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic-2.secret
  version = "1"
}
`, randomString)
}
