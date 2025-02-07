// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gemini_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// To run tests locally please replace the `oauth_token_secret_version` with your secret manager version.
// More details: https://cloud.google.com/developer-connect/docs/connect-github-repo#before_you_begin

func TestAccGeminiRepositoryGroupIamBinding(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", map[string]string{"ccfe_debug_note": "terraform_e2e_do_not_delete"})
	developerConnectionId := acctest.BootstrapDeveloperConnection(t, "basic", location, "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1", 54180648)
	gitRepositoryLinkId := acctest.BootstrapGitRepository(t, "basic", location, "https://github.com/CC-R-github-robot/tf-test.git", developerConnectionId)
	repositoryGroupId := "tf-test-iam-repository-group-id-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"role":                  "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index": codeRepositoryIndexId,
		"location":              location,
		"project":               envvar.GetTestProjectFromEnv(),
		"connection_id":         developerConnectionId,
		"git_link_id":           gitRepositoryLinkId,
		"repository_group_id":   repositoryGroupId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroupIamBinding_basic(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccGeminiRepositoryGroupIamBinding_update(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGeminiRepositoryGroupIamMember(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", map[string]string{"ccfe_debug_note": "terraform_e2e_do_not_delete"})
	developerConnectionId := acctest.BootstrapDeveloperConnection(t, "basic", location, "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1", 54180648)
	gitRepositoryLinkId := acctest.BootstrapGitRepository(t, "basic", location, "https://github.com/CC-R-github-robot/tf-test.git", developerConnectionId)
	repositoryGroupId := "tf-test-iam-repository-group-id-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"role":                  "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index": codeRepositoryIndexId,
		"location":              location,
		"project":               envvar.GetTestProjectFromEnv(),
		"connection_id":         developerConnectionId,
		"git_link_id":           gitRepositoryLinkId,
		"repository_group_id":   repositoryGroupId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccGeminiRepositoryGroupIamMember_basic(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s roles/cloudaicompanion.repositoryGroupsUser user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGeminiRepositoryGroupIamPolicy(t *testing.T) {
	location := "us-central1"
	codeRepositoryIndexId := acctest.BootstrapSharedCodeRepositoryIndex(t, "basic", location, "", map[string]string{"ccfe_debug_note": "terraform_e2e_do_not_delete"})
	developerConnectionId := acctest.BootstrapDeveloperConnection(t, "basic", location, "projects/502367051001/secrets/tf-test-cloudaicompanion-github-oauthtoken-c42e5c/versions/1", 54180648)
	gitRepositoryLinkId := acctest.BootstrapGitRepository(t, "basic", location, "https://github.com/CC-R-github-robot/tf-test.git", developerConnectionId)
	repositoryGroupId := "tf-test-iam-repository-group-id-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"role":                  "roles/cloudaicompanion.repositoryGroupsUser",
		"code_repository_index": codeRepositoryIndexId,
		"location":              location,
		"project":               envvar.GetTestProjectFromEnv(),
		"connection_id":         developerConnectionId,
		"git_link_id":           gitRepositoryLinkId,
		"repository_group_id":   repositoryGroupId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiRepositoryGroupIamPolicy_basic(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_gemini_repository_group_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGeminiRepositoryGroupIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_gemini_repository_group_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/codeRepositoryIndexes/%s/repositoryGroups/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), codeRepositoryIndexId, repositoryGroupId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGeminiRepositoryGroupIamMember_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_member" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "%{repository_group_id}"
  repositories {
    resource = "projects/%{project}/locations/us-central1/connections/%{connection_id}/gitRepositoryLinks/%{git_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}
`, context)
}

func testAccGeminiRepositoryGroupIamPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  depends_on = [
    google_gemini_repository_group_iam_policy.foo
  ]
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "%{repository_group_id}"
  repositories {
    resource = "projects/%{project}/locations/us-central1/connections/%{connection_id}/gitRepositoryLinks/%{git_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}
`, context)
}

func testAccGeminiRepositoryGroupIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_iam_policy" "foo" {
}

resource "google_gemini_repository_group_iam_policy" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  policy_data = data.google_iam_policy.foo.policy_data
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "%{repository_group_id}"
  repositories {
    resource = "projects/%{project}/locations/us-central1/connections/%{connection_id}/gitRepositoryLinks/%{git_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}
`, context)
}

func testAccGeminiRepositoryGroupIamBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_binding" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "%{repository_group_id}"
  repositories {
    resource = "projects/%{project}/locations/us-central1/connections/%{connection_id}/gitRepositoryLinks/%{git_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}
`, context)
}

func testAccGeminiRepositoryGroupIamBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gemini_repository_group_iam_binding" "foo" {
  project = "%{project}"
  location = "%{location}"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = google_gemini_repository_group.example.repository_group_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}

resource "google_gemini_repository_group" "example" {
  location = "us-central1"
  code_repository_index = "%{code_repository_index}"
  repository_group_id = "%{repository_group_id}"
  repositories {
    resource = "projects/%{project}/locations/us-central1/connections/%{connection_id}/gitRepositoryLinks/%{git_link_id}"
    branch_pattern = "main"
  }
  labels = {"label1": "value1"}
}
`, context)
}
