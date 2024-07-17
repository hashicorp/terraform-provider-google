// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privilegedaccessmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementProjectExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_name":  envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivilegedAccessManagerEntitlementDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample_basic(context),
			},
			{
				ResourceName:            "google_privileged_access_manager_entitlement.tfentitlement",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entitlement_id", "parent"},
			},
			{
				Config: testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample_update(context),
			},
			{
				ResourceName:            "google_privileged_access_manager_entitlement.tfentitlement",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entitlement_id", "parent"},
			},
		},
	})
}

func testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privileged_access_manager_entitlement" "tfentitlement" {
    entitlement_id = "tf-test-example-entitlement%{random_suffix}"
    location = "global"
    max_request_duration = "43200s"
    parent = "projects/%{project_name}"
    requester_justification_config { 
      unstructured{}
    }
    eligible_users {
      principals = ["group:test@google.com"]
    }
    privileged_access{
      gcp_iam_access{
        role_bindings{
          role = "roles/storage.admin"
          condition_expression = "request.time < timestamp(\"2024-04-23T18:30:00.000Z\")"
        }
        resource = "//cloudresourcemanager.googleapis.com/projects/%{project_name}"
        resource_type = "cloudresourcemanager.googleapis.com/Project"
      }
    }
    additional_notification_targets {
      admin_email_recipients     = ["user@example.com"]
      requester_email_recipients = ["user@example.com"]
    }
    approval_workflow {
    manual_approvals {
      require_approver_justification = true
      steps {
        approvals_needed          = 1
          approver_email_recipients = ["user@example.com"]
          approvers {
            principals = ["group:test@google.com"]
        }
      }
    }
  }
}
`, context)
}

func testAccPrivilegedAccessManagerEntitlement_privilegedAccessManagerEntitlementBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privileged_access_manager_entitlement" "tfentitlement" {
    entitlement_id = "tf-test-example-entitlement%{random_suffix}"
    location = "global"
    max_request_duration = "4300s"
    parent = "projects/%{project_name}"
    requester_justification_config {    
      not_mandatory{}
    }
    eligible_users {
      principals = ["group:test@google.com"]
    }
    privileged_access{
      gcp_iam_access{
        role_bindings{
          role = "roles/storage.admin"
          condition_expression = "request.time < timestamp(\"2024-04-23T18:30:00.000Z\")"
        }
        resource = "//cloudresourcemanager.googleapis.com/projects/%{project_name}"
        resource_type = "cloudresourcemanager.googleapis.com/Project"
      }
    }
    additional_notification_targets {
      admin_email_recipients     = ["user1@example.com"]
      requester_email_recipients = ["user2@example.com"]
    }
    approval_workflow {
    manual_approvals {
      require_approver_justification = false
      steps {
        approvals_needed          = 1
          approver_email_recipients = ["user3@example.com"]
          approvers {
            principals = ["group:test@google.com"]
        }
      }
    }
  }
}
`, context)
}
