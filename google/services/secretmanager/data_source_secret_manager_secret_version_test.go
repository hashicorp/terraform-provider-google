// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version.basic", "google_secret_manager_secret_version.secret-version-basic"),
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
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version.latest", "google_secret_manager_secret_version.secret-version-basic-2"),
				),
			},
		},
	})
}

func TestAccDatasourceSecretManagerSecretVersion_withBase64SecretData(t *testing.T) {
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	data := "./test-fixtures/binary-file.pfx"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecretManagerSecretVersion_withBase64SecretData(randomString, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceSecretManagerSecretVersion("data.google_secret_manager_secret_version.basic-base64", "1"),
					testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource("data.google_secret_manager_secret_version.basic-base64", "google_secret_manager_secret_version.secret-version-basic-base64"),
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

func testAccCheckSecretManagerSecretVersionSecretDataDatasourceMatchesResource(datasource, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("can't find Secret Version resource: %s", resource)
		}

		ds, ok := s.RootModule().Resources[datasource]
		if !ok {
			return fmt.Errorf("can't find Secret Version data source: %s", datasource)
		}

		if rs.Primary.ID == "" {
			return errors.New("resource ID not set.")
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		resourceSecretData, ok := rs.Primary.Attributes["secret_data"]
		if !ok {
			return errors.New("can't find 'secret_data' attribute in Secret Version resource")
		}

		datasourceSecretData, ok := ds.Primary.Attributes["secret_data"]
		if !ok {
			return errors.New("can't find 'secret_data' attribute in Secret Version data source")
		}

		if resourceSecretData != datasourceSecretData {
			return fmt.Errorf("expected %s, got %s, secret_data doesn't match", resourceSecretData, datasourceSecretData)
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
    auto {}
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

func testAccDatasourceSecretManagerSecretVersion_withBase64SecretData(randomString, data string) string {
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

data "google_secret_manager_secret_version" "basic-base64" {
  secret = google_secret_manager_secret_version.secret-version-basic-base64.secret
  is_secret_data_base64 = true
}
`, randomString, data)
}
