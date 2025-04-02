// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanagerregional_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersion_basicWithResourceReference(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameterVersion_basicWithResourceReference(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version.regional-parameter-version-basic", "google_parameter_manager_regional_parameter_version.regional-parameter-version-basic"),
				),
			},
		},
	})

}

func testAccParameterManagerRegionalRegionalParameterVersion_basicWithResourceReference(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_regional_parameter%{random_suffix}"
  location = "us-central1"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_regional_parameter_version%{random_suffix}"
  parameter_data = "test-regional-parameter-data-with-resource-reference"
}

data "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter_version.regional-parameter-version-basic.parameter
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-basic.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersion_basicWithParameterName(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameterVersion_basicWithParameterName(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version.regional-parameter-version-basic", "google_parameter_manager_regional_parameter_version.regional-parameter-version-basic"),
				),
			},
		},
	})

}

func testAccParameterManagerRegionalRegionalParameterVersion_basicWithParameterName(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_regional_parameter%{random_suffix}"
  location = "us-central1"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_regional_parameter_version%{random_suffix}"
  parameter_data = "test-regional-parameter-data-with-regional-parameter-name"
}

data "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-basic.parameter_version_id
  location = "us-central1"
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersion_withJsonData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameterVersion_withJsonData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data", "google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data"),
				),
			},
		},
	})

}

func testAccParameterManagerRegionalRegionalParameterVersion_withJsonData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_regional_parameter%{random_suffix}"
  format = "JSON"
  location = "us-central1"

}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-json-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_regional_parameter_version%{random_suffix}"
  parameter_data = jsonencode({
	"key1": "val1",
	"key2": "val2"
  })
}

data "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-json-data" {
  parameter = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data.parameter
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersion_withYamlData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameterVersion_withYamlData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data", "google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data"),
				),
			},
		},
	})

}

func testAccParameterManagerRegionalRegionalParameterVersion_withYamlData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_regional_parameter%{random_suffix}"
  format = "YAML"
  location = "us-central1"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_regional_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"key1": "val1",
	"key2": "val2"
  })
}

data "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data.parameter
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersion_withKmsKey(t *testing.T) {
	t.Parallel()

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pm.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	context := map[string]interface{}{
		"kms_key":       acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-parameter-manager-managed-central-key1").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameterVersion_withKmsKey(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version.regional-parameter-version-with-kms-key", "google_parameter_manager_regional_parameter_version.regional-parameter-version-with-kms-key"),
				),
			},
		},
	})

}

func testAccParameterManagerRegionalRegionalParameterVersion_withKmsKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_regional_parameter%{random_suffix}"
  format = "YAML"
  location = "us-central1"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-kms-key" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_regional_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"key1": "val1",
	"key2": "val2"
  })
}

data "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-kms-key" {
  parameter = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-kms-key.parameter
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-kms-key.parameter_version_id
}
`, context)
}

func testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource(dataSource, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("can't find Regional Parameter Version resource: %s", resource)
		}

		ds, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("can't find Regional Parameter Version data source: %s", dataSource)
		}

		if rs.Primary.ID == "" {
			return errors.New("resource ID not set.")
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		resourceParameterData, ok := rs.Primary.Attributes["parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Regoinal Parameter Version resource")
		}

		dataSourceParameterData, ok := ds.Primary.Attributes["parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Regional Parameter Version data source")
		}

		if resourceParameterData != dataSourceParameterData {
			return fmt.Errorf("expected %s, got %s, parameter_data doesn't match", resourceParameterData, dataSourceParameterData)
		}
		return nil
	}
}
