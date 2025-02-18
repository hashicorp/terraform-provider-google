// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/servicemanagement/Service.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples/base_configs/iam_test_file.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package servicemanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccServiceManagementServiceIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"project_name":  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceManagementServiceIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_endpoints_service_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("services/%s roles/viewer", fmt.Sprintf("endpoint%s.endpoints.%s.cloud.goog", context["random_suffix"], context["project_name"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccServiceManagementServiceIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_endpoints_service_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("services/%s roles/viewer", fmt.Sprintf("endpoint%s.endpoints.%s.cloud.goog", context["random_suffix"], context["project_name"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceManagementServiceIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"project_name":  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccServiceManagementServiceIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_endpoints_service_iam_member.foo",
				ImportStateId:     fmt.Sprintf("services/%s roles/viewer user:admin@hashicorptest.com", fmt.Sprintf("endpoint%s.endpoints.%s.cloud.goog", context["random_suffix"], context["project_name"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceManagementServiceIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"project_name":  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceManagementServiceIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_endpoints_service_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_endpoints_service_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("services/%s", fmt.Sprintf("endpoint%s.endpoints.%s.cloud.goog", context["random_suffix"], context["project_name"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceManagementServiceIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_endpoints_service_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("services/%s", fmt.Sprintf("endpoint%s.endpoints.%s.cloud.goog", context["random_suffix"], context["project_name"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceManagementServiceIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog"
  project = "%{project_name}"
  grpc_config = <<EOF
type: google.api.Service
config_version: 3
name: endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF
  protoc_output_base64 = "${filebase64("test-fixtures/test_api_descriptor.pb")}"
}

resource "google_endpoints_service_iam_member" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccServiceManagementServiceIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog"
  project = "%{project_name}"
  grpc_config = <<EOF
type: google.api.Service
config_version: 3
name: endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF
  protoc_output_base64 = "${filebase64("test-fixtures/test_api_descriptor.pb")}"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_endpoints_service_iam_policy" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_endpoints_service_iam_policy" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  depends_on = [
    google_endpoints_service_iam_policy.foo
  ]
}
`, context)
}

func testAccServiceManagementServiceIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog"
  project = "%{project_name}"
  grpc_config = <<EOF
type: google.api.Service
config_version: 3
name: endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF
  protoc_output_base64 = "${filebase64("test-fixtures/test_api_descriptor.pb")}"
}

data "google_iam_policy" "foo" {
}

resource "google_endpoints_service_iam_policy" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccServiceManagementServiceIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog"
  project = "%{project_name}"
  grpc_config = <<EOF
type: google.api.Service
config_version: 3
name: endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF
  protoc_output_base64 = "${filebase64("test-fixtures/test_api_descriptor.pb")}"
}

resource "google_endpoints_service_iam_binding" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccServiceManagementServiceIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog"
  project = "%{project_name}"
  grpc_config = <<EOF
type: google.api.Service
config_version: 3
name: endpoint%{random_suffix}.endpoints.%{project_name}.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF
  protoc_output_base64 = "${filebase64("test-fixtures/test_api_descriptor.pb")}"
}

resource "google_endpoints_service_iam_binding" "foo" {
  service_name = google_endpoints_service.endpoints_service.service_name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
