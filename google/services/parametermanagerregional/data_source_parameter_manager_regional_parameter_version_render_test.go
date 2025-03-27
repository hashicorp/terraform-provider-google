// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanagerregional_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersionRender_basicWithResourceReference(t *testing.T) {
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
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_basicWithResourceReferenceWithoutDatasource(context),
			},
			{
				// We've kept sleep because we need to grant the `Secret Manager Secret Accessor` role to the principal
				// of the parameter and it can take up to 7 minutes for the role to take effect. For more information
				// see the access change propagation documentation: https://cloud.google.com/iam/docs/access-change-propagation.
				PreConfig: func() {
					time.Sleep(7 * time.Minute)
				},
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_basicWithResourceReferenceWithDatasource(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-basic", "google_parameter_manager_regional_parameter_version.regional-parameter-version-basic"),
					testAccCheckParameterManagerRegionalRegionalRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-basic", "\"tempsecret\": \"regional-parameter-version-data\"\n"),
				),
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_basicWithResourceReferenceWithoutDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "YAML"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_basicWithResourceReferenceWithDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "YAML"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_regional_parameter_version_render" "regional-parameter-version-basic" {
  parameter = google_parameter_manager_regional_parameter_version.regional-parameter-version-basic.parameter
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-basic.parameter_version_id
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersionRender_withJsonData(t *testing.T) {
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
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_withJsonDataWithoutDatasource(context),
			},
			{
				// We've kept sleep because we need to grant the `Secret Manager Secret Accessor` role to the principal
				// of the parameter and it can take up to 7 minutes for the role to take effect. For more information
				// see the access change propagation documentation: https://cloud.google.com/iam/docs/access-change-propagation.
				PreConfig: func() {
					time.Sleep(7 * time.Minute)
				},
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_withJsonDataWithDatasource(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-with-json-data", "google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data"),
					testAccCheckParameterManagerRegionalRegionalRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-with-json-data", "{\"tempsecret\":\"regional-parameter-version-data\"}"),
				),
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_withJsonDataWithoutDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret_json_data%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-json-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = jsonencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_withJsonDataWithDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret_json_data%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-json-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = jsonencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_regional_parameter_version_render" "regional-parameter-version-with-json-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-json-data.parameter_version_id
  location = "us-central1"
}
`, context)
}

func TestAccDataSourceParameterManagerRegionalRegionalParameterVersionRender_withYamlData(t *testing.T) {
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
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_withYamlDataWithoutDatasource(context),
			},
			{
				// We've kept sleep because we need to grant the `Secret Manager Secret Accessor` role to the principal
				// of the parameter and it can take up to 7 minutes for the role to take effect. For more information
				// see the access change propagation documentation: https://cloud.google.com/iam/docs/access-change-propagation.
				PreConfig: func() {
					time.Sleep(7 * time.Minute)
				},
				Config: testAccParameterManagerRegionalRegionalParameterVersionRender_withYamlDataWithDatasource(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerRegionalRegionalParameterDataDataSourceMatchesResource("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-with-yaml-data", "google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data"),
					testAccCheckParameterManagerRegionalRegionalRenderedParameterDataMatchesDataSourceRenderedData("data.google_parameter_manager_regional_parameter_version_render.regional-parameter-version-with-yaml-data", "\"tempsecret\": \"regional-parameter-version-data\"\n"),
				),
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_withYamlDataWithoutDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "YAML"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret_yaml_data%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameterVersionRender_withYamlDataWithDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "YAML"
}

resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf_temp_secret_yaml_data%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.id
  secret_data = "regional-parameter-version-data"
}

resource "google_secret_manager_regional_secret_iam_member" "member" {
  secret_id = google_secret_manager_regional_secret.secret-basic.secret_id
  location = google_secret_manager_regional_secret.secret-basic.location
  role = "roles/secretmanager.secretAccessor"
  member = "${google_parameter_manager_regional_parameter.regional-parameter-basic.policy_member[0].iam_policy_uid_principal}"
}

resource "google_parameter_manager_regional_parameter_version" "regional-parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = yamlencode({
	"tempsecret": "__REF__(//secretmanager.googleapis.com/${google_secret_manager_regional_secret_version.secret-version-basic.name})"
  })
}

data "google_parameter_manager_regional_parameter_version_render" "regional-parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_regional_parameter.regional-parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_regional_parameter_version.regional-parameter-version-with-yaml-data.parameter_version_id
  location = "us-central1"
}
`, context)
}

func testAccCheckParameterManagerRegionalRegionalRenderedParameterDataMatchesDataSourceRenderedData(dataSource, expectedParameterData string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("can't find Regional Parameter Version Render data source: %s", dataSource)
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		dataSourceParameterData, ok := ds.Primary.Attributes["rendered_parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Regional Parameter Version Render data source")
		}

		if expectedParameterData != dataSourceParameterData {
			return fmt.Errorf("expected %s, got %s, rendered_parameter_data doesn't match", expectedParameterData, dataSourceParameterData)
		}
		return nil
	}
}
