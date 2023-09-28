// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package identityplatform_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIdentityPlatformConfig_update(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"org_id":           envvar.GetTestOrgFromEnv(t),
		"billing_acct":     envvar.GetTestBillingAccountFromEnv(t),
		"quota_start_time": time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"random_suffix":    acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPlatformConfig_basic(context),
			},
			{
				ResourceName:      "google_identity_platform_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityPlatformConfig_update(context),
			},
			{
				ResourceName:      "google_identity_platform_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityPlatformConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "basic" {
  project_id = "tf-test-my-project%{random_suffix}"
  name       = "tf-test-my-project%{random_suffix}"
  org_id     = "%{org_id}"
  billing_account =  "%{billing_acct}"
  labels = {
    firebase = "enabled"
  }
}

resource "google_project_service" "identitytoolkit" {
  project = google_project.basic.project_id
  service = "identitytoolkit.googleapis.com"
}

resource "google_identity_platform_config" "basic" {
  project = google_project.basic.project_id
  autodelete_anonymous_users = true
  sign_in {
    allow_duplicate_emails = true

    anonymous {
        enabled = true
    }
    email {
        enabled = true
        password_required = false
    }
    phone_number {
        enabled = true
        test_phone_numbers = {
            "+11231231234" = "000000"
        }
    }
  }
  sms_region_config {
    allow_by_default {
      disallowed_regions = [
        "CA",
        "US",
      ]
    }
  }
}
`, context)
}

func testAccIdentityPlatformConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "basic" {
  project_id = "tf-test-my-project%{random_suffix}"
  name       = "tf-test-my-project%{random_suffix}"
  org_id     = "%{org_id}"
  billing_account =  "%{billing_acct}"
  labels = {
    firebase = "enabled"
  }
}

resource "google_project_service" "identitytoolkit" {
  project = google_project.basic.project_id
  service = "identitytoolkit.googleapis.com"
}

resource "google_identity_platform_config" "basic" {
  project = google_project.basic.project_id
  sign_in {
    allow_duplicate_emails = false

    anonymous {
        enabled = false
    }
    email {
        enabled = true
        password_required = true
    }
    phone_number {
        enabled = true
        test_phone_numbers = {
	    "+17651212343" = "111111"
        }
    }
  }
  sms_region_config {
    allowlist_only {
      allowed_regions = [
        "AU",
        "NZ",
      ]
    }
  }
}
`, context)
}
