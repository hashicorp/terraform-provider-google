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
					testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_regional_secret_version.basic-1", "google_secret_manager_regional_secret_version.secret-version-basic"),
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
					testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_regional_secret_version.basic-2", "google_secret_manager_regional_secret_version.secret-version-basic"),
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
					testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_regional_secret_version.latest", "google_secret_manager_regional_secret_version.secret-version-basic-2"),
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
					testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_regional_secret_version.version", "google_secret_manager_regional_secret_version.secret-version-basic-1"),
				),
			},
		},
	})
}

func TestAccDataSourceSecretManagerRegionalRegionalSecretVersion_withBase64SecretData(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	data := "./test-fixtures/binary-file.pfx"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecretVersion_withBase64SecretData(randomString, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSecretManagerRegionalRegionalSecretVersion("data.google_secret_manager_regional_secret_version.basic-base64", "1"),
					testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_regional_secret_version.basic-base64", "google_secret_manager_regional_secret_version.secret-version-basic-base64"),
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

func testAccDataSourceSecretManagerRegionalRegionalSecretVersion_withBase64SecretData(randomString, data string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_regional_secret" "secret-basic-base64" {
  secret_id = "tf-test-secret-version-%s"
  location = "us-central1"
  labels = {
    label = "my-label"
  }
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic-base64" {
  secret = google_secret_manager_regional_secret.secret-basic-base64.name
  is_secret_data_base64 = true
  secret_data = filebase64("%s")
}

data "google_secret_manager_regional_secret_version" "basic-base64" {
  secret = google_secret_manager_regional_secret_version.secret-version-basic-base64.secret
  is_secret_data_base64 = true
}
`, randomString, data)
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

func testAccCheckSecretManagerRegionalRegionalSecretVersionSecretDataDatasourceMatchesResource(datasource, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("can't find Regional Secret Version resource: %s", resource)
		}

		ds, ok := s.RootModule().Resources[datasource]
		if !ok {
			return fmt.Errorf("can't find Regional Secret Version data source: %s", datasource)
		}

		if rs.Primary.ID == "" {
			return errors.New("resource ID not set.")
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		resourceSecretData, ok := rs.Primary.Attributes["secret_data"]
		if !ok {
			return errors.New("can't find 'secret_data' attribute in Regional Secret Version resource")
		}

		datasourceSecretData, ok := ds.Primary.Attributes["secret_data"]
		if !ok {
			return errors.New("can't find 'secret_data' attribute in Regional Secret Version data source")
		}

		if resourceSecretData != datasourceSecretData {
			return fmt.Errorf("expected %s, got %s, secret_data doesn't match", resourceSecretData, datasourceSecretData)
		}
		return nil
	}
}
