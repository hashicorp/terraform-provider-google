// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securesourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecureSourceManagerBranchRule_secureSourceManagerBranchRuleWithFieldsExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"prevent_destroy": false,
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecureSourceManagerBranchRule_secureSourceManagerBranchRuleWithFieldsExample_full(context),
			},
			{
				ResourceName:            "google_secure_source_manager_branch_rule.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branch_rule_id", "location", "repository_id"},
			},
			{
				Config: testAccSecureSourceManagerBranchRule_secureSourceManagerBranchRuleWithFieldsExample_update(context),
			},
			{
				ResourceName:            "google_secure_source_manager_branch_rule.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branch_rule_id", "location", "repository_id"},
			},
		},
	})
}

func testAccSecureSourceManagerBranchRule_secureSourceManagerBranchRuleWithFieldsExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
    location = "us-central1"
    instance_id = "tf-test-my-initial-instance%{random_suffix}"
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_repository" "repository" {
    repository_id = "tf-test-my-initial-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name
    location = google_secure_source_manager_instance.instance.location
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_branch_rule" "default" {
    branch_rule_id = "tf-test-my-initial-branchrule%{random_suffix}"
    location = google_secure_source_manager_repository.repository.location
    repository_id = google_secure_source_manager_repository.repository.repository_id
    include_pattern = "test"
    minimum_approvals_count   = 2
    minimum_reviews_count     = 2
    require_comments_resolved = true
    require_linear_history    = true
    require_pull_request      = true
    disabled = false
    allow_stale_reviews = false
}
`, context)
}

func testAccSecureSourceManagerBranchRule_secureSourceManagerBranchRuleWithFieldsExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
	location = "us-central1"
	instance_id = "tf-test-my-initial-instance%{random_suffix}"
	# Prevent accidental deletions.
	lifecycle {
		prevent_destroy = "%{prevent_destroy}"
	}
}

resource "google_secure_source_manager_repository" "repository" {
    repository_id = "tf-test-my-initial-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name
    location = google_secure_source_manager_instance.instance.location
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_branch_rule" "default" {
    branch_rule_id = "tf-test-my-initial-branchrule%{random_suffix}"
	location = google_secure_source_manager_repository.repository.location
    repository_id = google_secure_source_manager_repository.repository.repository_id
    include_pattern = "test"
    minimum_approvals_count   = 1
    minimum_reviews_count     = 1
    require_linear_history    = false
}
`, context)
}
