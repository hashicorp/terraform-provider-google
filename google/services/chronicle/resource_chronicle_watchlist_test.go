// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package chronicle_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccChronicleWatchlist_chronicleWatchlistBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"chronicle_id":  envvar.GetTestChronicleInstanceIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckChronicleWatchlistDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccChronicleWatchlist_chronicleWatchlistBasicExample_basic(context),
			},
			{
				ResourceName:            "google_chronicle_watchlist.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "location", "watchlist_id"},
			},
			{
				Config: testAccChronicleWatchlist_chronicleWatchlistBasicExample_update(context),
			},
			{
				ResourceName:            "google_chronicle_watchlist.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "location", "watchlist_id"},
			},
		},
	})
}

func testAccChronicleWatchlist_chronicleWatchlistBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_watchlist" "example" {
  location = "us"
  instance = "%{chronicle_id}"
  watchlist_id = "tf-test-watchlist-name%{random_suffix}"
  description = "tf-test-watchlist-description%{random_suffix}"
  display_name = "tf-test-watchlist-name%{random_suffix}"
  multiplying_factor = 1
  entity_population_mechanism {
    manual {

    }
  }
}
`, context)
}

func testAccChronicleWatchlist_chronicleWatchlistBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_chronicle_watchlist" "example" {
  location = "us"
  instance = "%{chronicle_id}"
  watchlist_id = "tf-test-watchlist-name%{random_suffix}"
  description = "tf-test-watchlist-updated-description%{random_suffix}"
  display_name = "tf-test-updated-watchlist%{random_suffix}"
  multiplying_factor = 2
  entity_population_mechanism {
    manual {

    }
  }
  watchlist_user_preferences {
    pinned = true
  }
}
`, context)
}
