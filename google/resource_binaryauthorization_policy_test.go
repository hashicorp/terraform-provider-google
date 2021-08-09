package google

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBinaryAuthorizationPolicy_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	billingId := getTestBillingAccountFromEnv(t)
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func TestAccBinaryAuthorizationPolicy_full(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	billingId := getTestBillingAccountFromEnv(t)
	note := randString(t, 10)
	attestor := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, pname, org, billingId, note, attestor, "ENABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

// Use an attestor created in the default CI project
func TestAccBinaryAuthorizationPolicy_separateProject(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	billingId := getTestBillingAccountFromEnv(t)
	note := randString(t, 10)
	attestor := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicy_separateProject(pid, pname, org, billingId, note, attestor),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func TestAccBinaryAuthorizationPolicy_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	billingId := getTestBillingAccountFromEnv(t)
	note := randString(t, 10)
	attestor := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, pname, org, billingId, note, attestor, "ENABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyFull(pid, pname, org, billingId, note, attestor, "DISABLE"),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billingId),
			},
			{
				ResourceName:      "google_binary_authorization_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy the policy without destroying the project so we can check
			// that it was restored to the default.
			{
				Config: testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billingId),
				Check:  testAccCheckBinaryAuthorizationPolicyDefault(t, pid),
			},
		},
	})
}

func testAccCheckBinaryAuthorizationPolicyDefault(t *testing.T, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		url := fmt.Sprintf("https://binaryauthorization.googleapis.com/v1/projects/%s/policy", pid)
		pol, err := sendRequest(config, "GET", "", url, config.userAgent, nil)
		if err != nil {
			return err
		}
		delete(pol, "updateTime")

		defaultPol := defaultBinaryAuthorizationPolicy(pid)
		if !reflect.DeepEqual(pol, defaultPol) {
			return fmt.Errorf("Policy for project %s was %v, expected default policy %v", pid, pol, defaultPol)
		}
		return nil
	}
}

func testAccBinaryAuthorizationPolicyDefault(pid, pname, org, billing string) string {
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
`, pid, pname, org, billing)
}

func testAccBinaryAuthorizationPolicyBasic(pid, pname, org, billing string) string {
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
`, pid, pname, org, billing)
}

func testAccBinaryAuthorizationPolicyFull(pid, pname, org, billing, note, attestor, gpmode string) string {
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
`, pid, pname, org, billing, note, attestor, gpmode)
}

func testAccBinaryAuthorizationPolicy_separateProject(pid, pname, org, billing, note, attestor string) string {
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
`, pid, pname, org, billing, note, attestor)
}
