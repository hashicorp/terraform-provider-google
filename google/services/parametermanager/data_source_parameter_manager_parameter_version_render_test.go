// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerParameterVersionRender_basicWithResourceReference(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerParameterVersionRender_basicWithResourceReference(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version_render.parameter-version-basic", "google_parameter_manager_parameter_version.parameter-version-basic"),
					testAccCheckParameterManagerRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_parameter_version_render.parameter-version-basic", "\"tempsecret\": \"parameter-version-data\"\n"),
				),
			},
		},
	})
}

func testAccParameterManagerParameterVersionRender_basicWithResourceReference(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "YAML"
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-temp-secret-basic%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "parameter-version-data"
}

resource "google_secret_manager_secret_iam_member" "member" {
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_parameter.parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_parameter_version_render" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter_version.parameter-version-basic.parameter
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-basic.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerParameterVersionRender_withJsonData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerParameterVersionRender_withJsonData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version_render.parameter-version-with-json-data", "google_parameter_manager_parameter_version.parameter-version-with-json-data"),
					testAccCheckParameterManagerRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_parameter_version_render.parameter-version-with-json-data", "{\"tempsecret\":\"parameter-version-data\"}"),
				),
			},
		},
	})
}

func testAccParameterManagerParameterVersionRender_withJsonData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "JSON"
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-temp-secret-json-data%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "parameter-version-data"
}

resource "google_secret_manager_secret_iam_member" "member" {
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_parameter.parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-with-json-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = jsonencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_parameter_version_render" "parameter-version-with-json-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-with-json-data.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerParameterVersionRender_withYamlData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerParameterVersionRender_withYamlData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version_render.parameter-version-with-yaml-data", "google_parameter_manager_parameter_version.parameter-version-with-yaml-data"),
					testAccCheckParameterManagerRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_parameter_version_render.parameter-version-with-yaml-data", "\"tempsecret\": \"parameter-version-data\"\n"),
				),
			},
		},
	})
}

func testAccParameterManagerParameterVersionRender_withYamlData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "YAML"
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-temp-secret-yaml-data%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "parameter-version-data"
}

resource "google_secret_manager_secret_iam_member" "member" {
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_parameter.parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_parameter_version_render" "parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-with-yaml-data.parameter_version_id
}
`, context)
}

func testAccCheckParameterManagerRenderedParameterDataMatchesDataSourceRenderedData(dataSource, expectedParameterData string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("can't find Parameter Version Render data source: %s", dataSource)
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		dataSourceParameterData, ok := ds.Primary.Attributes["rendered_parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Parameter Version Render data source")
		}

		if expectedParameterData != dataSourceParameterData {
			return fmt.Errorf("expected %s, got %s, rendered_parameter_data doesn't match", expectedParameterData, dataSourceParameterData)
		}
		return nil
	}
}
