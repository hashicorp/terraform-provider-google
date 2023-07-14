// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package orgpolicy_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	orgpolicy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
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
}


`, context)
}

func testAccOrgPolicyPolicy_OrganizationPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "organizations/%{org_id}/policies/gcp.detailedAuditLoggingMode"
  parent = "organizations/%{org_id}"

  spec {
    reset = true
  }
}


`, context)
}

func testAccOrgPolicyPolicy_OrganizationPolicyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name   = "organizations/%{org_id}/policies/gcp.detailedAuditLoggingMode"
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
}


`, context)
}

func testAccCheckOrgPolicyPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_org_policy_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &orgpolicy.Policy{
				Name:   dcl.String(rs.Primary.Attributes["name"]),
				Parent: dcl.String(rs.Primary.Attributes["parent"]),
			}

			client := transport_tpg.NewDCLOrgPolicyClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetPolicy(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_org_policy_policy still exists %v", obj)
			}
		}
		return nil
	}
}
