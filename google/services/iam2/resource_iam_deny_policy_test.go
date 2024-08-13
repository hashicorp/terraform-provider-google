// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iam2_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIAM2DenyPolicy_iamDenyPolicyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAM2DenyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate2(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
		},
	})
}

func TestAccIAM2DenyPolicy_iamDenyPolicyFolderParent(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAM2DenyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyFolder(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyFolderUpdate(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
		},
	})
}

func testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_iam_deny_policy" "example" {
  parent   = urlencode("cloudresourcemanager.googleapis.com/projects/${google_project.project.project_id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "First rule"
    deny_rule {
      denied_principals = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test-account.email}"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.update"]
    }
  }
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.update"]
      exception_principals = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test-account.email}"]
    }
  }
}

resource "google_service_account" "test-account" {
  account_id   = "tf-test-deny-account%{random_suffix}"
  display_name = "Test Service Account"
  project      = google_project.project.project_id
}
`, context)
}

func testAccIAM2DenyPolicy_iamDenyPolicyUpdate2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_iam_deny_policy" "example" {
  parent   = urlencode("cloudresourcemanager.googleapis.com/projects/${google_project.project.project_id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some other expr"
        expression = "!resource.matchTag('87654321/env', 'test')"
        location = "/some/file"
        description = "A denial condition"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.update"]
    }
  }
}

resource "google_service_account" "test-account" {
  account_id   = "tf-test-deny-account%{random_suffix}"
  display_name = "Test Service Account"
  project      = google_project.project.project_id
}
`, context)
}

func testAccIAM2DenyPolicy_iamDenyPolicyFolder(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_deny_policy" "example" {
  parent   = urlencode("cloudresourcemanager.googleapis.com/${google_folder.folder.id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
    }
  }
}

resource "google_folder" "folder" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}
`, context)
}

func testAccIAM2DenyPolicy_iamDenyPolicyFolderUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_deny_policy" "example" {
  parent   = urlencode("cloudresourcemanager.googleapis.com/${google_folder.folder.id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
    }
  }
}

resource "google_folder" "folder" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}
`, context)
}
