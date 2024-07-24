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

package privilegedaccessmanager_test

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

func TestAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivilegedAccessManagerEntitlementDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample(context),
			},
			{
				ResourceName:            "google_privileged_access_manager_entitlement.tfentitlement",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entitlement_id", "location", "parent"},
			},
		},
	})
}

func testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privileged_access_manager_entitlement" "tfentitlement" {
    entitlement_id = "tf-test-example-entitlement%{random_suffix}"
    location = "global"
    max_request_duration = "43200s"
    parent = "projects/%{project}"
    requester_justification_config {    
        unstructured{}
    }
    eligible_users {
        principals = [
          "group:test@google.com"
        ]
    }
    privileged_access{
        gcp_iam_access{
            role_bindings{
                role = "roles/storage.admin"
                condition_expression = "request.time < timestamp(\"2024-04-23T18:30:00.000Z\")"
            }
            resource = "//cloudresourcemanager.googleapis.com/projects/%{project}"
            resource_type = "cloudresourcemanager.googleapis.com/Project"
        }
    }
    additional_notification_targets {
      admin_email_recipients     = [
        "user@example.com",
      ]
      requester_email_recipients = [
        "user@example.com"
      ]
    }
    approval_workflow {
    manual_approvals {
      require_approver_justification = true
      steps {
        approvals_needed          = 1
        approver_email_recipients = [
          "user@example.com"
        ]
        approvers {
          principals = [
            "group:test@google.com"
          ]
        }
      }
    }
  }
}
`, context)
}

func testAccCheckPrivilegedAccessManagerEntitlementDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_privileged_access_manager_entitlement" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{PrivilegedAccessManagerBasePath}}{{parent}}/locations/{{location}}/entitlements/{{entitlement_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("PrivilegedAccessManagerEntitlement still exists at %s", url)
			}
		}

		return nil
	}
}
