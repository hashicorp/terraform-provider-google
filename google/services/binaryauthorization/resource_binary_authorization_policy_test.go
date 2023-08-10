// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package binaryauthorization_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/binaryauthorization"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccBinaryAuthorizationPolicy_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func TestAccBinaryAuthorizationPolicy_full(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	note := acctest.RandString(t, 10)
	attestor := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, org, billingId, note, attestor, "ENABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

// Use an attestor created in the default CI project
func TestAccBinaryAuthorizationPolicy_separateProject(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	note := acctest.RandString(t, 10)
	attestor := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicy_separateProject(pid, org, billingId, note, attestor),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func TestAccBinaryAuthorizationPolicy_update(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := "tf-test-" + acctest.RandString(t, 10)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	note := acctest.RandString(t, 10)
	attestor := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, org, billingId, note, attestor, "ENABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, org, billingId, note, attestor, "DISABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func testAccCheckBinaryAuthorizationPolicyDefault(t *testing.T, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		url := fmt.Sprintf("https://binaryauthorization.googleapis.com/v1/projects/%s/policy", pid)
		pol, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}

		// new fields will cause this test to fail- if they're simple outputs, we can just ignore them.
		delete(pol, "updateTime")
		delete(pol, "etag")

		defaultPol := binaryauthorization.DefaultBinaryAuthorizationPolicy(pid)
		if !reflect.DeepEqual(pol, defaultPol) {
			return fmt.Errorf("Policy for project %s was %v, expected default policy %v", pid, pol, defaultPol)
		}
		return nil
	}
}

func testAccBinaryAuthorizationPolicyDefault(pid, org, billing string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}
`, pid, pid, org, billing)
}

func testAccBinaryAuthorizationPolicyBasic(pid, org, billing string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}

resource "google_binary_authorization_policy" "policy" {
  project = google_project.project.project_id

  admission_whitelist_patterns {
    name_pattern = "gcr.io/google_containers/*"
  }

  default_admission_rule {
    evaluation_mode  = "ALWAYS_DENY"
    enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"
  }

  depends_on = [google_project_service.binauthz]
}
`, pid, pid, org, billing)
}

func testAccBinaryAuthorizationPolicyFull(pid, org, billing, note, attestor, gpmode string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}

resource "google_container_analysis_note" "note" {
  project = google_project.project.project_id

  name = "tf-test-%s"
  attestation_authority {
    hint {
      human_readable_name = "My attestor"
    }
  }

  depends_on = [google_project_service.binauthz]
}

resource "google_binary_authorization_attestor" "attestor" {
  project = google_project.project.project_id

  name        = "tf-test-%s"
  description = "my description"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
  }

  depends_on = [google_project_service.binauthz]
}

resource "google_binary_authorization_policy" "policy" {
  project = google_project.project.project_id

  admission_whitelist_patterns {
    name_pattern = "gcr.io/google_containers/*"
  }

  default_admission_rule {
    evaluation_mode  = "ALWAYS_ALLOW"
    enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"
  }

  cluster_admission_rules {
    cluster                 = "us-central1-a.prod-cluster"
    evaluation_mode         = "REQUIRE_ATTESTATION"
    enforcement_mode        = "ENFORCED_BLOCK_AND_AUDIT_LOG"
    require_attestations_by = [google_binary_authorization_attestor.attestor.name]
  }

  global_policy_evaluation_mode = "%s"
}
`, pid, pid, org, billing, note, attestor, gpmode)
}

func testAccBinaryAuthorizationPolicy_separateProject(pid, org, billing, note, attestor string) string {
	return fmt.Sprintf(`
// Use a separate project since each project can only have one policy
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

data "google_client_config" "current" {
}

resource "google_project_service" "binauthz" {
  project = google_project.project.project_id
  service = "binaryauthorization.googleapis.com"
}

resource "google_container_analysis_note" "note" {
  project = google_project.project.project_id

  name = "tf-test-%s"
  attestation_authority {
    hint {
      human_readable_name = "My attestor"
    }
  }

  depends_on = [google_project_service.binauthz]
}

resource "google_binary_authorization_attestor" "attestor" {
  name        = "tf-test-%s"
  description = "my description"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
  }

  depends_on = [google_project_service.binauthz]
}

resource "google_binary_authorization_policy" "policy" {
  project = google_project.project.project_id

  admission_whitelist_patterns {
    name_pattern = "gcr.io/google_containers/*"
  }

  default_admission_rule {
    evaluation_mode  = "ALWAYS_ALLOW"
    enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"
  }

  cluster_admission_rules {
    cluster                 = "us-central1-a.prod-cluster"
    evaluation_mode         = "REQUIRE_ATTESTATION"
    enforcement_mode        = "ENFORCED_BLOCK_AND_AUDIT_LOG"
    require_attestations_by = ["projects/${data.google_client_config.current.project}/attestors/${google_binary_authorization_attestor.attestor.name}"]
  }
}
`, pid, pid, org, billing, note, attestor)
}
