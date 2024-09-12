// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecretManagerRegionalRegionalSecretIam_iamPolicyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/secretmanager.secretAccessor",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalRegionalSecretIam_iamPolicyBasic(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_iam_policy.default",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/secrets/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-tf-reg-secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerRegionalRegionalSecretIam_iamPolicyUpdate(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_iam_policy.default",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/secrets/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-tf-reg-secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecretManagerRegionalRegionalSecretIam_iamPolicyBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "default" {
  secret_id = "tf-test-tf-reg-secret%{random_suffix}"
  location = "us-central1"
  ttl = "3600s"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }
}

data "google_iam_policy" "default" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_secret_manager_regional_secret_iam_policy" "default" {
  project = google_secret_manager_regional_secret.default.project
  location = google_secret_manager_regional_secret.default.location
  secret_id = google_secret_manager_regional_secret.default.secret_id
  policy_data = data.google_iam_policy.default.policy_data
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretIam_iamPolicyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "default" {
  secret_id = "tf-test-tf-reg-secret%{random_suffix}"
  location = "us-central1"
  ttl = "3600s"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }
}

data "google_iam_policy" "default" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
  }
}

resource "google_secret_manager_regional_secret_iam_policy" "default" {
  project = google_secret_manager_regional_secret.default.project
  location = google_secret_manager_regional_secret.default.location
  secret_id = google_secret_manager_regional_secret.default.secret_id
  policy_data = data.google_iam_policy.default.policy_data
}
`, context)
}

func TestAccSecretManagerRegionalRegionalSecretIam_iamBindingUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/secretmanager.secretAccessor",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalRegionalSecretIam_iamBindingBasic(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_iam_binding.default",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/secrets/%s roles/secretmanager.secretAccessor", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-tf-reg-secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerRegionalRegionalSecretIam_iamBindingUpdate(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_iam_binding.default",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/secrets/%s roles/secretmanager.secretAccessor", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-tf-reg-secret%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecretManagerRegionalRegionalSecretIam_iamBindingBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "default" {
  secret_id = "tf-test-tf-reg-secret%{random_suffix}"
  location = "us-central1"
  ttl = "3600s"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }
}

resource "google_secret_manager_regional_secret_iam_binding" "default" {
  project = google_secret_manager_regional_secret.default.project
  location = google_secret_manager_regional_secret.default.location
  secret_id = google_secret_manager_regional_secret.default.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretIam_iamBindingUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "default" {
  secret_id = "tf-test-tf-reg-secret%{random_suffix}"
  location = "us-central1"
  ttl = "3600s"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }
}

resource "google_secret_manager_regional_secret_iam_binding" "default" {
  project = google_secret_manager_regional_secret.default.project
  location = google_secret_manager_regional_secret.default.location
  secret_id = google_secret_manager_regional_secret.default.secret_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}
