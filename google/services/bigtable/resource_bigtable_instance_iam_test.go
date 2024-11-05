// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigtableInstanceIamBinding(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s %s",
		envvar.GetTestProjectFromEnv(), context["instance"].(string), context["role"].(string))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamBinding_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_binding.binding", "role", context["role"].(string)),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigtableInstanceIamBinding_update(context),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamMember(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s %s serviceAccount:%s",
		envvar.GetTestProjectFromEnv(),
		context["instance"].(string),
		context["role"].(string),
		envvar.ServiceAccountCanonicalEmail(context["account"].(string)))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamMember(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "role", context["role"].(string)),
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "member", "serviceAccount:"+envvar.ServiceAccountCanonicalEmail(context["account"].(string))),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamPolicy(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s",
		envvar.GetTestProjectFromEnv(), context["instance"].(string))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamPolicy(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigtable_instance_iam_policy.policy", "policy_data"),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamBinding_withCondition(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	conditionExpression := "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
	context := map[string]interface{}{
		"instance":        "tf-bigtable-iam-" + randomString,
		"cluster":         "c-" + randomString,
		"account":         "tf-bigtable-iam-" + randomString,
		"role":            "roles/bigtable.user",
		"condition_title": "expires_after_2019_12_31",
		"condition_expr":  strconv.Quote(conditionExpression),
		"condition_desc":  "Expiring at midnight of 2019-12-31",
	}

	importIdWithCondition := fmt.Sprintf("projects/%s/instances/%s %s %s",
		envvar.GetTestProjectFromEnv(),
		context["instance"].(string),
		context["role"].(string),
		context["condition_title"].(string),
	)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstanceIamBinding_withCondition(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_binding.binding", "condition.#", "1"),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_binding.binding", "condition.0.title", context["condition_title"].(string)),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_binding.binding", "condition.0.expression", conditionExpression),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_binding.binding", "condition.0.description", context["condition_desc"].(string)),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importIdWithCondition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamMember_withCondition(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	conditionExpression := "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
	context := map[string]interface{}{
		"instance":        "tf-bigtable-iam-" + randomString,
		"cluster":         "c-" + randomString,
		"account":         "tf-bigtable-iam-" + randomString,
		"role":            "roles/bigtable.user",
		"condition_title": "expires_after_2019_12_31",
		"condition_expr":  strconv.Quote(conditionExpression),
		"condition_desc":  "Expiring at midnight of 2019-12-31",
	}

	importIdWithCondition := fmt.Sprintf("projects/%s/instances/%s %s serviceAccount:%s %s",
		envvar.GetTestProjectFromEnv(),
		context["instance"].(string),
		context["role"].(string),
		envvar.ServiceAccountCanonicalEmail(context["account"].(string)),
		context["condition_title"].(string),
	)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamMember_withCondition(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_member.member", "condition.#", "1"),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_member.member", "condition.0.title", context["condition_title"].(string)),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_member.member", "condition.0.expression", conditionExpression),
					resource.TestCheckResourceAttr("google_bigtable_instance_iam_member.member", "condition.0.description", context["condition_desc"].(string)),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_member.member",
				ImportStateId:     importIdWithCondition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamPolicy_withCondition(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	conditionExpression := "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
	context := map[string]interface{}{
		"instance":        "tf-bigtable-iam-" + randomString,
		"cluster":         "c-" + randomString,
		"account":         "tf-bigtable-iam-" + randomString,
		"role":            "roles/bigtable.user",
		"condition_title": "expires_after_2019_12_31",
		"condition_expr":  strconv.Quote(conditionExpression),
		"condition_desc":  "Expiring at midnight of 2019-12-31",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s",
		envvar.GetTestProjectFromEnv(), context["instance"].(string))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamPolicy_withCondition(context),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableInstanceIamBinding_basic(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account1" {
  account_id   = "%{account}-1"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account}-2"
  display_name = "Bigtable instance Iam Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, context)
}

func testAccBigtableInstanceIamBinding_update(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account1" {
  account_id   = "%{account}-1"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account}-2"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, context)
}

func testAccBigtableInstanceIamBinding_withCondition(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account1" {
  account_id   = "%{account}-1"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account}-2"
  display_name = "Bigtable instance Iam Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = %{condition_expr}
  }
}
`, context)
}

func testAccBigtableInstanceIamMember(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_bigtable_instance_iam_member" "member" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  member   = "serviceAccount:${google_service_account.test-account.email}"
}
`, context)
}

func testAccBigtableInstanceIamMember_withCondition(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_bigtable_instance_iam_member" "member" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  member   = "serviceAccount:${google_service_account.test-account.email}"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = %{condition_expr}
  }
}
`, context)
}

func testAccBigtableInstanceIamPolicy(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%{role}"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
}
`, context)
}

func testAccBigtableInstanceIamPolicy_withCondition(context map[string]interface{}) string {
	return testBigtableInstanceConfig(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%{role}"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
    condition {
      title       = "%{condition_title}"
      description = "%{condition_desc}"
      expression  = %{condition_expr}
    }
  }
}

resource "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
}
`, context)
}

// Smallest instance possible for testing
func testBigtableInstanceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigtable_instance" "instance" {
	name                  = "%{instance}"
    instance_type = "DEVELOPMENT"

    cluster {
      cluster_id   = "%{cluster}"
      zone         = "us-central1-b"
      storage_type = "HDD"
    }

    deletion_protection = false
}
`, context)
}
