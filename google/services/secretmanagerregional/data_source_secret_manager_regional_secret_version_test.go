// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithResourceReference(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithResourceReference(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version.basic-1", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithSecretName(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithSecretName(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version.basic-2", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersion_latest(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersion_latest(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version.latest", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersion_versionField(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersion_versionField(randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version.version", "1"),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithResourceReference(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-%s"
}

data "google_secret_manager_regional_secret_version" "basic-1" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic.secret
}
`, randomString, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersion_basicWithSecretName(randomString string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "my-tf-test-secret-%s"
}

data "google_secret_manager_regional_secret_version" "basic-2" {
  secret = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret_version.secret-version-basic.location
}
`, randomString, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersion_latest(randomString string) string {
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

data "google_secret_manager_regional_secret_version" "latest" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic-2.secret
}
`, randomString)
}

func testAccDataSourceSecretManagerRegionalRegionalSecretVersion_versionField(randomString string) string {
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

data "google_secret_manager_regional_secret_version" "version" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic-2.secret
  version = "1"
}
`, randomString)
}

func testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion(n, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Regional Secret Version data source: %s", n)
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
