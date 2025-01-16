// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gemini_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccGeminiCodeRepositoryIndex_update(t *testing.T) {
	bootstrappedKMS := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    os.Getenv("GOOGLE_PROJECT"),
		"kms_key":       bootstrappedKMS.CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiCodeRepositoryIndex_basic(context),
			},
			{
				ResourceName:            "google_gemini_code_repository_index.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccGeminiCodeRepositoryIndex_update(context),
			},
			{
				ResourceName:            "google_gemini_code_repository_index.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

// TestAccGeminiCodeRepositoryIndex_delete checks if there is no error in deleting CRI along with children resource
// note: this is an example of a bad usage, where RGs refer to the CRI using a string id, not a reference, as they
// will be force-removed upon CRI deletion, because the CRI provider uses --force option by default
// The plan after the _delete function should not be empty due to the child resource in plan
func TestAccGeminiCodeRepositoryIndex_delete(t *testing.T) {
	bootstrappedKMS := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix": randomSuffix,
		"project_id":    os.Getenv("GOOGLE_PROJECT"),
		"kms_key":       bootstrappedKMS.CryptoKey.Name,
		"cri_id":        fmt.Sprintf("tf-test-cri-index-delete-example-%s", randomSuffix),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiCodeRepositoryIndex_withChildren_basic(context),
			},
			{
				ResourceName:            "google_gemini_code_repository_index.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index_id", "labels", "location", "terraform_labels", "force_destroy"},
			},
			{
				Config:             testAccGeminiCodeRepositoryIndex_withChildren_delete(context),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func testAccGeminiCodeRepositoryIndex_withChildren_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_code_repository_index" "example" {
  labels = {"ccfe_debug_note": "terraform_e2e_should_be_deleted"}
  location = "us-central1"
  code_repository_index_id = "%{cri_id}"
  force_destroy = true
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{cri_id}"
  repository_group_id = "tf-test-rg-repository-group-id-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
  depends_on = [
    google_gemini_code_repository_index.example
  ]
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "tf-test-repository-conn-delete"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion-delete-%{random_suffix}"
  disabled = false

  github_config {
    github_app = "DEVELOPER_CONNECT"
    app_installation_id = 54180648

    authorizer_credential {
      oauth_token_secret_version = "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1"
    }
  }
}
`, context)
}

// Removed depends_on to not break plan test
func testAccGeminiCodeRepositoryIndex_withChildren_delete(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{cri_id}"
  repository_group_id = "tf-test-rg-repository-group-id-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "tf-test-repository-conn-delete"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion-delete-%{random_suffix}"
  disabled = false

  github_config {
    github_app = "DEVELOPER_CONNECT"
    app_installation_id = 54180648

    authorizer_credential {
      oauth_token_secret_version = "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1"
    }
  }
}
`, context)
}

func testAccGeminiCodeRepositoryIndex_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_code_repository_index" "example" {
  labels = {"ccfe_debug_note": "terraform_e2e_should_be_deleted"}
  location = "us-central1"
  code_repository_index_id = "tf-test-cri-index-example-%{random_suffix}"
  kms_key = "%{kms_key}"
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key_binding]
}

data "google_project" "project" {
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_binding" {
  crypto_key_id = "%{kms_key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloudaicompanion.iam.gserviceaccount.com",
  ]
}
`, context)
}

func testAccGeminiCodeRepositoryIndex_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_code_repository_index" "example" {
  labels = {"ccfe_debug_note": "terraform_e2e_should_be_deleted", "new_label": "new_val"}
  location = "us-central1"
  code_repository_index_id = "tf-test-cri-index-example-%{random_suffix}"
  kms_key = "%{kms_key}"
}
`, context)
}
