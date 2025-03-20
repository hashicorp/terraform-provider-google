// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigqueryanalyticshub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryAnalyticsHubListingSubscription_differentProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubListingSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListingSubscription_differentProject(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_subscription.subscription",
				ImportStateIdFunc: testAccBigqueryAnalyticsHubListingSubscription_stateId,
				ImportState:       true,
				// skipping ImportStateVerify as the resource ID won't match original
				// since the user cannot input the project and destination projects simultaneously
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListingSubscription_stateId(state *terraform.State) (string, error) {
	resourceName := "google_bigquery_analytics_hub_listing_subscription.subscription"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}

	return fmt.Sprintf("projects/%s/locations/US/subscriptions/%s", envvar.GetTestProjectFromEnv(), rawState["subscription_id"]), nil
}

func testAccBigqueryAnalyticsHubListingSubscription_differentProject(context map[string]interface{}) string {
	return acctest.Nprintf(`


# Dataset created in default project
resource "google_bigquery_dataset" "subscription" {
	dataset_id                  = "tf_test_my_listing%{random_suffix}"
	friendly_name               = "tf_test_my_listing%{random_suffix}"
	description                 = ""
	location                    = "US"
}

resource "google_project" "project" {
	project_id      = "tf-test-%{random_suffix}"
	name            = "tf-test-%{random_suffix}"
	org_id          = "%{org_id}"
	deletion_policy = "DELETE"
}


resource "google_project_service" "analyticshub" {
	project  = google_project.project.project_id
	service  = "analyticshub.googleapis.com"
	disable_on_destroy = false # Need it enabled in the project when the test disables services in post-test cleanup
}

resource "google_bigquery_analytics_hub_data_exchange" "subscription" {
  project          = google_project.project.project_id
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = ""
  depends_on = [google_project_service.analyticshub]
}

resource "google_bigquery_analytics_hub_listing" "subscription" {
  project          = google_project.project.project_id
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.subscription.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = ""

  bigquery_dataset {
    dataset = google_bigquery_dataset.subscription.id
  }
}

resource "google_bigquery_analytics_hub_listing_subscription" "subscription" {
  project          = google_project.project.project_id
  location = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.subscription.data_exchange_id
  listing_id       = google_bigquery_analytics_hub_listing.subscription.listing_id
  destination_dataset {
    description = "A test subscription"
    friendly_name = "👋"
    labels = {
      testing = "123"
    }
    location = "US"
    dataset_reference {
      dataset_id = "tf_test_destination_dataset%{random_suffix}"
      project_id = google_bigquery_dataset.subscription.project
    }
  }
}
`, context)
}
