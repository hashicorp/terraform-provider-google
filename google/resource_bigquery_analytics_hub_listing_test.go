package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryAnalyticsHubListingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingBasicExample(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(context),
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(context map[string]interface{}) string {
	return Nprintf(`
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
  description      = "example data exchange update%{random_suffix}"

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
`, context)
}
