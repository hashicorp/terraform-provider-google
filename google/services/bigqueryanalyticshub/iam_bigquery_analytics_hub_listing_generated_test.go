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

package bigqueryanalyticshub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryAnalyticsHubListingIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListingIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s roles/viewer", envvar.GetTestProjectFromEnv(), "US", fmt.Sprintf("tf_test_my_data_exchange%s", context["random_suffix"]), fmt.Sprintf("tf_test_my_listing%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccBigqueryAnalyticsHubListingIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s roles/viewer", envvar.GetTestProjectFromEnv(), "US", fmt.Sprintf("tf_test_my_data_exchange%s", context["random_suffix"]), fmt.Sprintf("tf_test_my_listing%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryAnalyticsHubListingIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccBigqueryAnalyticsHubListingIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s roles/viewer user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), "US", fmt.Sprintf("tf_test_my_data_exchange%s", context["random_suffix"]), fmt.Sprintf("tf_test_my_listing%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryAnalyticsHubListingIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListingIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigquery_analytics_hub_listing_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s", envvar.GetTestProjectFromEnv(), "US", fmt.Sprintf("tf_test_my_data_exchange%s", context["random_suffix"]), fmt.Sprintf("tf_test_my_listing%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListingIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s", envvar.GetTestProjectFromEnv(), "US", fmt.Sprintf("tf_test_my_data_exchange%s", context["random_suffix"]), fmt.Sprintf("tf_test_my_listing%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListingIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

resource "google_bigquery_analytics_hub_listing_iam_member" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccBigqueryAnalyticsHubListingIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_bigquery_analytics_hub_listing_iam_policy" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_bigquery_analytics_hub_listing_iam_policy" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  depends_on = [
    google_bigquery_analytics_hub_listing_iam_policy.foo
  ]
}
`, context)
}

func testAccBigqueryAnalyticsHubListingIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

data "google_iam_policy" "foo" {
}

resource "google_bigquery_analytics_hub_listing_iam_policy" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccBigqueryAnalyticsHubListingIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

resource "google_bigquery_analytics_hub_listing_iam_binding" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccBigqueryAnalyticsHubListingIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

resource "google_bigquery_analytics_hub_listing_iam_binding" "foo" {
  project = google_bigquery_analytics_hub_listing.listing.project
  location = google_bigquery_analytics_hub_listing.listing.location
  data_exchange_id = google_bigquery_analytics_hub_listing.listing.data_exchange_id
  listing_id = google_bigquery_analytics_hub_listing.listing.listing_id
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
