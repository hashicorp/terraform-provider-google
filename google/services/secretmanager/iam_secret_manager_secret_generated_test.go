// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package secretmanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccSecretManagerSecretIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccSecretManagerSecretIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_secret_manager_secret_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerSecretIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamBindingGenerated_withCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamBinding_withConditionGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamBindingGenerated_withAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamBinding_withAndWithoutConditionGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo2",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_binding.foo3",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamMemberGenerated_withCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamMember_withConditionGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor user:admin@hashicorptest.com %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamMemberGenerated_withAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
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
				Config: testAccSecretManagerSecretIamMember_withAndWithoutConditionGenerated(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.foo2",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor user:admin@hashicorptest.com %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.foo3",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s roles/secretmanager.secretAccessor user:admin@hashicorptest.com %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerSecretIamPolicyGenerated_withCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/secretmanager.secretAccessor",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	// Test should have 2 bindings: one with a description and one without. Any < chars are converted to a unicode character by the API.
	expectedPolicyData := acctest.Nprintf(`{"bindings":[{"condition":{"description":"%{condition_desc}","expression":"%{condition_expr}","title":"%{condition_title}"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"},{"condition":{"expression":"%{condition_expr}","title":"%{condition_title}-no-description"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"}]}`, context)
	expectedPolicyData = strings.Replace(expectedPolicyData, "<", "\\u003c", -1)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecretIamPolicy_withConditionGenerated(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO(SarahFrench) - uncomment once https://github.com/GoogleCloudPlatform/magic-modules/pull/6466 merged
					// resource.TestCheckResourceAttr("data.google_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttr("google_secret_manager_secret_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttrWith("data.google_iam_policy.foo", "policy_data", tpgresource.CheckGoogleIamPolicy),
				),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecretManagerSecretIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_member" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccSecretManagerSecretIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_secret_manager_secret_iam_policy" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_secret_manager_secret_iam_policy" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  depends_on = [
    google_secret_manager_secret_iam_policy.foo
  ]
}
`, context)
}

func testAccSecretManagerSecretIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

data "google_iam_policy" "foo" {
}

resource "google_secret_manager_secret_iam_policy" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccSecretManagerSecretIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_binding" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccSecretManagerSecretIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_binding" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccSecretManagerSecretIamBinding_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_binding" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
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

func testAccSecretManagerSecretIamBinding_withAndWithoutConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_binding" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}

resource "google_secret_manager_secret_iam_binding" "foo2" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_secret_manager_secret_iam_binding" "foo3" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
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

func testAccSecretManagerSecretIamMember_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_member" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
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

func testAccSecretManagerSecretIamMember_withAndWithoutConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_iam_member" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}

resource "google_secret_manager_secret_iam_member" "foo2" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_secret_manager_secret_iam_member" "foo3" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
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

func testAccSecretManagerSecretIamPolicy_withConditionGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
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
}

resource "google_secret_manager_secret_iam_policy" "foo" {
  project = google_secret_manager_secret.secret-basic.project
  secret_id = google_secret_manager_secret.secret-basic.secret_id
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}
