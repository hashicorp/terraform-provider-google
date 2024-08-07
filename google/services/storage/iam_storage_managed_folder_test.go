// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccStorageManagedFolderIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccStorageManagedFolderIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccStorageManagedFolderIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer user:admin@hashicorptest.com", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	// This may skip test, so do it first
	sa := envvar.GetTestServiceAccountFromEnv(t)
	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}
	context["service_account"] = sa

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_storage_managed_folder_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageManagedFolderIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamBindingGenerated_withCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamBinding_withConditionGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamBindingGenerated_withAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamBinding_withAndWithoutConditionGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo2",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_binding.foo3",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamMemberGenerated_withCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamMember_withConditionGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamMemberGenerated_withAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamMember_withAndWithoutConditionGenerated(context),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer user:admin@hashicorptest.com", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_member.foo2",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_member.foo3",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/ roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageManagedFolderIamPolicyGenerated_withCondition(t *testing.T) {
	t.Parallel()

	// This may skip test, so do it first
	sa := envvar.GetTestServiceAccountFromEnv(t)
	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}
	context["service_account"] = sa

	// Test should have 3 bindings: one with a description and one without, and a third for an admin role. Any < chars are converted to a unicode character by the API.
	expectedPolicyData := acctest.Nprintf(`{"bindings":[{"members":["serviceAccount:%{service_account}"],"role":"%{admin_role}"},{"condition":{"description":"%{condition_desc}","expression":"%{condition_expr}","title":"%{condition_title}"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"},{"condition":{"expression":"%{condition_expr}","title":"%{condition_title}-no-description"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"}]}`, context)
	expectedPolicyData = strings.Replace(expectedPolicyData, "<", "\\u003c", -1)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolderIamPolicy_withConditionGenerated(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO(SarahFrench) - uncomment once https://github.com/GoogleCloudPlatform/magic-modules/pull/6466 merged
					// resource.TestCheckResourceAttr("data.google_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttr("google_storage_managed_folder_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttrWith("data.google_iam_policy.foo", "policy_data", tpgresource.CheckGoogleIamPolicy),
				),
			},
			{
				ResourceName:      "google_storage_managed_folder_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s/managedFolders/managed/folder/name/", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStorageManagedFolderIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_member" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccStorageManagedFolderIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
  binding {
    role = "%{admin_role}"
    members = ["serviceAccount:%{service_account}"]
  }
}

resource "google_storage_managed_folder_iam_policy" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_storage_managed_folder_iam_policy" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  depends_on = [
    google_storage_managed_folder_iam_policy.foo
  ]
}
`, context)
}

func testAccStorageManagedFolderIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

data "google_iam_policy" "foo" {
}

resource "google_storage_managed_folder_iam_policy" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccStorageManagedFolderIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_binding" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccStorageManagedFolderIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_binding" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccStorageManagedFolderIamBinding_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_binding" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}
`, context)
}

func testAccStorageManagedFolderIamBinding_withAndWithoutConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_binding" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}

resource "google_storage_managed_folder_iam_binding" "foo2" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_storage_managed_folder_iam_binding" "foo3" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    # Check that lack of description doesn't cause any issues
    # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
    title       = "%{condition_title_no_desc}"
    expression  = "%{condition_expr_no_desc}"
  }
}
`, context)
}

func testAccStorageManagedFolderIamMember_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_member" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}
`, context)
}

func testAccStorageManagedFolderIamMember_withAndWithoutConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

resource "google_storage_managed_folder_iam_member" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}

resource "google_storage_managed_folder_iam_member" "foo2" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_storage_managed_folder_iam_member" "foo3" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    # Check that lack of description doesn't cause any issues
    # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
    title       = "%{condition_title_no_desc}"
    expression  = "%{condition_expr_no_desc}"
  }
}
`, context)
}

func testAccStorageManagedFolderIamPolicy_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_managed_folder" "folder" {
  bucket = google_storage_bucket.bucket.name
  name   = "managed/folder/name/"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
    condition {
      # Check that lack of description doesn't cause any issues
      # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
      title       = "%{condition_title_no_desc}"
      expression  = "%{condition_expr_no_desc}"
    }
  }
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
    condition {
      title       = "%{condition_title}"
      description = "%{condition_desc}"
      expression  = "%{condition_expr}"
    }
  }
  binding {
    role = "%{admin_role}"
    members = ["serviceAccount:%{service_account}"]
  }
}

resource "google_storage_managed_folder_iam_policy" "foo" {
  bucket         = google_storage_managed_folder.folder.bucket
  managed_folder = google_storage_managed_folder.folder.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}
