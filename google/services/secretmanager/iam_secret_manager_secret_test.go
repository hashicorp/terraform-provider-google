// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecretManagerSecretIam_iamMemberConditionUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/secretmanager.secretAccessor",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecretIam_iamMemberCondition_basic(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.default",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s %s serviceAccount:%s %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-secret-%s", context["random_suffix"]), context["role"], fmt.Sprintf("tf-test-sa-%s@%s.iam.gserviceaccount.com", context["random_suffix"], envvar.GetTestProjectFromEnv()), fmt.Sprintf("tf-test-condition-%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerSecretIam_iamMemberCondition_update(context),
			},
			{
				ResourceName:      "google_secret_manager_secret_iam_member.default",
				ImportStateId:     fmt.Sprintf("projects/%s/secrets/%s %s serviceAccount:%s %s", envvar.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-secret-%s", context["random_suffix"]), context["role"], fmt.Sprintf("tf-test-sa-%s@%s.iam.gserviceaccount.com", context["random_suffix"], envvar.GetTestProjectFromEnv()), fmt.Sprintf("tf-test-condition-new-%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecretManagerSecretIam_iamMemberCondition_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "default" {
  account_id   = "tf-test-sa-%{random_suffix}"
  display_name = "Secret manager IAM testing account"
}

resource "google_secret_manager_secret" "default" {
  secret_id = "tf-test-secret-%{random_suffix}"
  ttl       = "3600s"

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

resource "google_secret_manager_secret_iam_member" "default" {
  secret_id = google_secret_manager_secret.default.id
  role      = "%{role}"
  member    = "serviceAccount:${google_service_account.default.email}"
  condition {
    title       = "tf-test-condition-%{random_suffix}"
    description = "test condition"
    expression  = "request.time < timestamp(\"2022-03-01T00:00:00Z\")"
  }
}
`, context)
}

func testAccSecretManagerSecretIam_iamMemberCondition_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "default" {
  account_id   = "tf-test-sa-%{random_suffix}"
  display_name = "Secret manager IAM testing account"
}

resource "google_secret_manager_secret" "default" {
  secret_id = "tf-test-secret-%{random_suffix}"
  ttl       = "3600s"

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

resource "google_secret_manager_secret_iam_member" "default" {
  secret_id = google_secret_manager_secret.default.id
  role      = "%{role}"
  member    = "serviceAccount:${google_service_account.default.email}"
  condition {
    title       = "tf-test-condition-new-%{random_suffix}"
    description = "test new condition"
    expression  = "request.time < timestamp(\"2024-03-01T00:00:00Z\")"
  }
}
`, context)
}
