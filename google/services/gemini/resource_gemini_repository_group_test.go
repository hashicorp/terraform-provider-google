// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gemini_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// To run tests locally please replace the `oauth_token_secret_version` with your secret manager version.
// More details: https://cloud.google.com/developer-connect/docs/connect-github-repo#before_you_begin

func TestAccGeminiRepositoryGroup_update(t *testing.T) {
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic-rg-test", "us-central1", "", map[string]string{"ccfe_debug_note": "terraform_e2e_do_not_delete"})
	context := map[string]interface{}{
		"random_suffix":         acctest.RandString(t, 10),
		"project_id":            os.Getenv("GOOGLE_PROJECT"),
		"code_repository_index": codeRepositoryIndexId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroup_basic(context),
			},
			{
				ResourceName:            "google_gemini_repository_group.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index", "labels", "location", "repository_group_id", "terraform_labels"},
			},
			{
				Config: testAccGeminiRepositoryGroup_update(context),
			},
			{
				ResourceName:            "google_gemini_repository_group.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index", "labels", "location", "repository_group_id", "terraform_labels"},
			},
		},
	})
}

func TestAccGeminiRepositoryGroup_noBootstrap(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    os.Getenv("GOOGLE_PROJECT"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroup_noBootstrap(context),
			},
			{
				ResourceName:            "google_gemini_repository_group.example_e",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_repository_index", "labels", "location", "repository_group_id", "terraform_labels"},
			},
		},
	})
}

func testAccGeminiRepositoryGroup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "tf-test-rg-repository-group-id-%{random_suffix}" 
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "tf-test-repository-conn"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion2-%{random_suffix}"
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
func testAccGeminiRepositoryGroup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "tf-test-rg-repository-group-id-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1", "label2": "value2"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "tf-test-repository-conn"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion3-%{random_suffix}"
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

func testAccGeminiRepositoryGroup_noBootstrap(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_code_repository_index" "cri" {
  labels = {"ccfe_debug_note": "terraform_e2e_should_be_deleted"}
  location = "us-central1"
  code_repository_index_id = "tf-test-rg-index-example-%{random_suffix}"
}

resource "google_gemini_repository_group" "example_a" {
  location = "us-central1"
  code_repository_index = google_gemini_code_repository_index.cri.code_repository_index_id
  repository_group_id = "tf-test-rg-nb-repository-group-id1-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_gemini_repository_group" "example_b" {
  location = "us-central1"
  code_repository_index = google_gemini_code_repository_index.cri.code_repository_index_id
  repository_group_id = "tf-test-rg-nb-repository-group-id2-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_gemini_repository_group" "example_c" {
  location = "us-central1"
  code_repository_index = google_gemini_code_repository_index.cri.code_repository_index_id
  repository_group_id = "tf-test-rg-nb-repository-group-id3-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_gemini_repository_group" "example_d" {
  location = "us-central1"
  code_repository_index = google_gemini_code_repository_index.cri.code_repository_index_id
  repository_group_id = "tf-test-rg-nb-repository-group-id4-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_gemini_repository_group" "example_e" {
  location = "us-central1"
  code_repository_index = google_gemini_code_repository_index.cri.code_repository_index_id
  repository_group_id = "tf-test-rg-nb-repository-group-id5-%{random_suffix}"
  repositories {
    resource = "projects/%{project_id}/locations/us-central1/connections/${google_developer_connect_connection.github_conn.connection_id}/gitRepositoryLinks/${google_developer_connect_git_repository_link.conn.git_repository_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}

resource "google_developer_connect_git_repository_link" "conn" {
  git_repository_link_id = "tf-test-repository-conn"
  parent_connection = google_developer_connect_connection.github_conn.connection_id
  clone_uri = "https://github.com/CC-R-github-robot/tf-test.git"
  location = "us-central1"
  annotations = {}
}

resource "google_developer_connect_connection" "github_conn" {
  location = "us-central1"
  connection_id = "tf-test-cloudaicompanion1-%{random_suffix}"
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
