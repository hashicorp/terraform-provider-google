// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package orgpolicy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccOrgPolicyPolicy_EnforcePolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_EnforcePolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func TestAccOrgPolicyPolicy_FolderPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_FolderPolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
			{
				Config: testAccOrgPolicyPolicy_FolderPolicyUpdate0(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func TestAccOrgPolicyPolicy_OrganizationPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_OrganizationPolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
			{
				Config: testAccOrgPolicyPolicy_OrganizationPolicyUpdate0(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func TestAccOrgPolicyPolicy_ProjectPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_ProjectPolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
			{
				Config: testAccOrgPolicyPolicy_ProjectPolicyUpdate0(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func TestAccOrgPolicyPolicy_DryRunSpecHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_DryRunSpecHandWritten(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}

func testAccOrgPolicyPolicy_EnforcePolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/iam.disableServiceAccountKeyUpload"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      enforce = "FALSE"
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}


`, context)
}

func testAccOrgPolicyPolicy_FolderPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "${google_folder.basic.name}/policies/gcp.resourceLocations"
  parent = google_folder.basic.name

  spec {
    inherit_from_parent = true

    rules {
      deny_all = "TRUE"
    }
  }
}

resource "google_folder" "basic" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder%{random_suffix}"
  deletion_protection = false
}


`, context)
}

func testAccOrgPolicyPolicy_FolderPolicyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "${google_folder.basic.name}/policies/gcp.resourceLocations"
  parent = google_folder.basic.name

  spec {
    inherit_from_parent = false

    rules {
      condition {
        description = "A sample condition for the policy"
        expression  = "resource.matchLabels('labelKeys/123', 'labelValues/345')"
        title       = "sample-condition"
      }

      values {
        allowed_values = ["projects/allowed-project"]
        denied_values  = ["projects/denied-project"]
      }
    }

    rules {
      allow_all = "TRUE"
    }
  }
}

resource "google_folder" "basic" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder%{random_suffix}"
  deletion_protection = false
}


`, context)
}

func testAccOrgPolicyPolicy_OrganizationPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE", "UPDATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}

resource "google_org_policy_policy" "primary" {
  name   = "organizations/%{org_id}/policies/${google_org_policy_custom_constraint.constraint.name}"
  parent = "organizations/%{org_id}"

  spec {
    reset = true
  }
}
`, context)
}

func testAccOrgPolicyPolicy_OrganizationPolicyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE", "UPDATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}

resource "google_org_policy_policy" "primary" {
  name   = "organizations/%{org_id}/policies/${google_org_policy_custom_constraint.constraint.name}"
  parent = "organizations/%{org_id}"

  spec {
    reset = false

    rules {
      enforce = "TRUE"
    }
  }
}
`, context)
}

func testAccOrgPolicyPolicy_ProjectPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/gcp.resourceLocations"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      condition {
        description = "A sample condition for the policy"
        expression  = "resource.matchLabels('labelKeys/123', 'labelValues/345')"
        location    = "sample-location.log"
        title       = "sample-condition"
      }

      values {
        allowed_values = ["projects/allowed-project"]
        denied_values  = ["projects/denied-project"]
      }
    }

    rules {
      allow_all = "TRUE"
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}


`, context)
}

func testAccOrgPolicyPolicy_ProjectPolicyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/gcp.resourceLocations"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      condition {
        description = "A new sample condition for the policy"
        expression  = "false"
        location    = "new-sample-location.log"
        title       = "new-sample-condition"
      }

      values {
        allowed_values = ["projects/new-allowed-project"]
        denied_values  = ["projects/new-denied-project"]
      }
    }

    rules {
      deny_all = "TRUE"
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}


`, context)
}

func testAccOrgPolicyPolicy_DryRunSpecHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.disableGkeAutoUpgrade%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}

resource "google_org_policy_policy" "primary" {
  name   = "organizations/%{org_id}/policies/${google_org_policy_custom_constraint.constraint.name}"
  parent = "organizations/%{org_id}"

  spec {
    rules {
      enforce = "FALSE"
    }
  }
  dry_run_spec {
    inherit_from_parent = false
    reset               = false
    rules {
      enforce = "FALSE"
    }
  }
}

`, context)
}

func testAccCheckOrgPolicyPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_org_policy_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{OrgPolicyBasePath}}{{parent}}/policies/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				Project:              billingProject,
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsOrgpolicyRetryableError},
			})
			if err == nil {
				return fmt.Errorf("OrgPolicyPolicy still exists at %s", url)
			}
		}

		return nil
	}
}
func TestAccOrgPolicyPolicy_EnforceParameterizedMCPolicy(t *testing.T) {
	// Skip this test as no constraints yet launched in production, verified functionality with manual testing.
	t.Skip()
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_EnforceParameterizedMCPolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func testAccOrgPolicyPolicy_EnforceParameterizedMCPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/essentialcontacts.managed.allowedContactDomains"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      enforce = "TRUE"
      parameters = "{\"allowedDomains\": [\"@google.com\"]}"
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}


`, context)
}

func TestAccOrgPolicyPolicy_EnforceParameterizedMCDryRunPolicy(t *testing.T) {
	// Skip this test as no constraints yet launched in production, verified functionality with manual testing.
	t.Skip()
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyPolicy_EnforceParameterizedMCDryRunPolicy(context),
			},
			{
				ResourceName:            "google_org_policy_policy.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "spec.0.rules.0.condition.0.expression"},
			},
		},
	})
}
func testAccOrgPolicyPolicy_EnforceParameterizedMCDryRunPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/essentialcontacts.managed.allowedContactDomains"
  parent = "projects/${google_project.basic.name}"

  dry_run_spec {
    rules {
      enforce = "TRUE"
      parameters = "{\"allowedDomains\": [\"@google.com\"]}"
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
  deletion_policy = "DELETE"
}


`, context)
}
